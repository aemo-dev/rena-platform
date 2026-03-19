<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useI18n } from 'vue-i18n'
import { Languages, Plus, Layout, Calendar, ChevronRight, LogOut, Search, Grid, List, Trash2, Folder, Smartphone, Globe, ChevronDown } from 'lucide-vue-next'
import iconImg from '../assets/icon.png'
import { setLanguage } from '../i18n'
import gsap from 'gsap'
import * as projectsApi from '../services/projects'

const router = useRouter()
const { t } = useI18n()
const auth = useAuthStore()
const projects = ref<{ id: string, name: string, created_at: string, icon_url?: string }[]>([])
const loading = ref(true)
const searchQuery = ref('')
const showLangMenu = ref(false)
const viewMode = ref<'grid' | 'list'>('grid')
const showCreateDialog = ref(false)
const showDeleteDialog = ref(false)
const showLogoutDialog = ref(false)
const projectToDelete = ref<{ id: string, name: string } | null>(null)
const deleteConfirmationText = ref('')
const isProjectsExpanded = ref(true)

const userDisplayName = computed(() => {
  return auth.user?.user_metadata?.full_name || auth.user?.user_metadata?.name || auth.user?.email?.split('@')[0] || 'User'
})

const userInitial = computed(() => {
  return userDisplayName.value.charAt(0).toUpperCase()
})

const filteredProjects = computed(() => {
  if (!searchQuery.value) return projects.value
  const query = searchQuery.value.toLowerCase()
  return projects.value.filter(p => p.name.toLowerCase().includes(query))
})

const newProject = ref({
  name: '',
  packageName: '',
  platform: 'android',
  color: '#2e60ff'
})

const errors = ref({
  name: '',
  packageName: ''
})

const colorPresets = ['#2e60ff', '#ff4081', '#4caf50', '#ffc107', '#00bcd4', '#3f51b5', '#795548', '#607d8b']

const changeLanguage = (lang: string) => {
  setLanguage(lang)
  showLangMenu.value = false
}

const fetchProjects = async () => {
  if (!auth.user) return
  
  const response = await projectsApi.getProjects()
  
  if (response.error) {
    console.error('Error fetching projects:', response.error)
  } else if (response.data) {
    projects.value = response.data.projects || []
  }
  
  loading.value = false
  
  // Animate cards after fetch
  setTimeout(() => {
    const cards = document.querySelectorAll(".project-card")
    if (cards.length > 0) {
      gsap.from(cards, {
        y: 20,
        opacity: 0,
        duration: 0.5,
        stagger: 0.1,
        ease: "power2.out"
      })
    }
  }, 100)
}

const validateAppName = (name: string) => {
  if (!name) return 'App Name is required'
  // React Native project names must be alphanumeric and start with a letter
  if (!/^[a-zA-Z][a-zA-Z0-9]*$/.test(name)) {
    return 'Name must start with a letter and contain only letters and numbers (no spaces)'
  }
  if (projects.value.some(p => p.name.toLowerCase() === name.toLowerCase())) {
    return 'Project name already exists'
  }
  return ''
}

const validatePackageName = (pkg: string) => {
  if (!pkg) return 'Package Name is required'
  // Standard Android package name validation: com.example.app
  if (!/^[a-z][a-z0-9_]*(\.[a-z0-9_]+)+[0-9a-z_]$/.test(pkg)) {
    return 'Invalid package name (e.g., com.example.myapp)'
  }
  return ''
}

const deleteProject = (project: { id: string, name: string }) => {
  projectToDelete.value = project
  deleteConfirmationText.value = ''
  showDeleteDialog.value = true
}

const confirmDelete = async () => {
  if (!projectToDelete.value || deleteConfirmationText.value !== projectToDelete.value.name) return

  const response = await projectsApi.deleteProject(projectToDelete.value.id)

  if (response.error) {
    alert('Error deleting project: ' + response.error)
  } else {
    projects.value = projects.value.filter(p => p.id !== projectToDelete.value?.id)
    closeDeleteDialog()
  }
}

const closeDeleteDialog = () => {
  showDeleteDialog.value = false
  projectToDelete.value = null
}

const openCreateDialog = () => {
  showCreateDialog.value = true
  newProject.value = {
    name: '',
    packageName: 'com.rena.myapp',
    platform: 'android',
    color: '#2e60ff'
  }
  errors.value = { name: '', packageName: '' }
}

const closeCreateDialog = () => {
  showCreateDialog.value = false
}

const handleCreateProject = async () => {
  errors.value.name = validateAppName(newProject.value.name)
  errors.value.packageName = validatePackageName(newProject.value.packageName)

  if (errors.value.name || errors.value.packageName) return

  loading.value = true

  try {
    // 2. Create the project
    const response = await projectsApi.createProject({ 
      name: newProject.value.name, 
      package_name: newProject.value.packageName,
      platform: newProject.value.platform,
      color: newProject.value.color,
      icon_url: iconImg, // Default icon URL
    })
    
    if (response.error) {
      if (response.error.includes('column "color"')) {
        alert('Missing database column: Please run the SQL command provided in the chat to add the "color" column to your projects table.')
      } else {
        alert('Error creating project: ' + response.error)
      }
    } else if (response.data && response.data.project) {
      projects.value.unshift(response.data.project)
      closeCreateDialog()
    }
  } catch (err) {
    console.error('Error in project creation flow:', err)
    alert('Failed to create project.')
  } finally {
    loading.value = false
  }
}

const openProject = (id: string) => {
  router.push(`/project/${id}`)
}

const confirmLogout = () => {
  showLogoutDialog.value = true
}

const handleLogout = async () => {
  await auth.logout()
  await auth.fetchUser()
  router.push('/dashboard')
}

const closeLogoutDialog = () => {
  showLogoutDialog.value = false
}

onMounted(() => {
  document.title = 'Rena Creator'
  fetchProjects()
})
</script>

<template>
  <div class="dashboard-wrapper">
    <!-- Sidebar-style Header (Modern & Rare) -->
    <header class="side-nav">
      <div class="nav-top">
        <div class="logo" @click="router.push('/')">
          <img :src="iconImg" alt="Rena Logo" class="logo-img" />
          <span class="logo-text">Rena Builder</span>
        </div>
        
        <nav class="main-menu">
          <button class="menu-item active" @click="isProjectsExpanded = !isProjectsExpanded">
            <Layout :size="20" />
            <span>My Projects</span>
            <ChevronDown :size="16" :class="{ 'rotated': !isProjectsExpanded }" class="expand-icon" />
          </button>
        </nav>

        <div class="sidebar-section" v-if="projects.length > 0 && isProjectsExpanded">
          <div class="sidebar-project-list">
            <div v-for="p in projects" :key="p.id" class="sidebar-project-item" @click="openProject(p.id)">
              <Folder :size="16" />
              <span>{{ p.name }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="nav-bottom">
        <div class="lang-switcher">
          <button @click="showLangMenu = !showLangMenu" class="nav-icon-btn">
            <Languages :size="20" />
          </button>
          <div v-if="showLangMenu" class="lang-dropdown">
            <button @click="changeLanguage('en')">English</button>
            <button @click="changeLanguage('ar')">العربية</button>
          </div>
        </div>
        
        <div class="user-profile" v-if="auth.user">
          <div class="avatar-container tooltip">
            <div class="avatar">{{ userInitial }}</div>
            <div class="tooltip-content">
              <span class="tooltip-name">{{ userDisplayName }}</span>
              <span class="tooltip-email">{{ auth.user.email }}</span>
            </div>
          </div>
          <div class="user-details">
            <span class="email-truncate">{{ userDisplayName }}</span>
          </div>
          <button @click="confirmLogout" class="logout-btn" :title="t('dashboard.logout')">
            <LogOut :size="18" />
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content Area -->
    <main class="content-area">
      <div class="top-bar">
        <div class="search-box">
          <Search :size="18" />
          <input type="text" v-model="searchQuery" :placeholder="t('dashboard.search')" />
        </div>
        
        <div class="actions">
          <div class="view-toggle">
            <button @click="viewMode = 'grid'" :class="{ active: viewMode === 'grid' }"><Grid :size="18" /></button>
            <button @click="viewMode = 'list'" :class="{ active: viewMode === 'list' }"><List :size="18" /></button>
          </div>
          <button @click="openCreateDialog" class="btn-primary">
            <Plus :size="20" />
            <span>{{ t('dashboard.create') }}</span>
          </button>
          
          <div class="user-top-profile tooltip" v-if="auth.user">
            <div class="top-avatar">
              {{ userInitial }}
            </div>
            <div class="tooltip-content">
              <span class="tooltip-name">{{ userDisplayName }}</span>
              <span class="tooltip-email">{{ auth.user.email }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="projects-scroll">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>{{ t('dashboard.loading') }}</p>
        </div>
        
        <div v-else-if="projects.length === 0" class="empty-state">
          <div class="empty-icon">📂</div>
          <h3>{{ t('dashboard.empty') }}</h3>
          <button @click="openCreateDialog" class="btn-primary">{{ t('dashboard.create') }}</button>
        </div>

        <div v-else :class="['projects-container', viewMode]">
          <div 
            v-for="p in filteredProjects" 
            :key="p.id" 
            class="project-card clickable" 
            @click="openProject(p.id)"
          >
            <div class="card-top">
              <div class="project-icon" :style="{ background: ((p as any).color || '#2e60ff') + '15', color: (p as any).color || '#2e60ff' }">
                <img v-if="p.icon_url" :src="p.icon_url" class="project-icon-img" />
                <img v-else :src="iconImg" class="project-icon-img" />
              </div>
              <div class="project-info">
                <h3>{{ p.name }}</h3>
                <div class="date">
                  <Calendar :size="12" />
                  {{ new Date(p.created_at).toLocaleDateString() }}
                </div>
              </div>
            </div>
            <div class="card-actions">
              <button class="delete-btn" @click.stop="deleteProject(p)" :title="'Delete ' + p.name">
                <Trash2 :size="18" />
              </button>
              <button class="open-btn">
                {{ t('dashboard.open') }}
                <ChevronRight :size="16" class="rtl-flip" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Create Project Dialog -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showCreateDialog" class="modal-overlay" @click.self="closeCreateDialog">
          <div class="modal-content">
            <div class="modal-header">
              <h2>{{ t('dashboard.createDialog.title') }}</h2>
              <button class="close-btn" @click="closeCreateDialog">&times;</button>
            </div>
            
            <div class="modal-body">
              <div class="form-group">
                <label>{{ t('dashboard.createDialog.appName') }}</label>
                <input 
                  type="text" 
                  v-model="newProject.name" 
                  placeholder="MyAwesomeApp"
                  :class="{ 'error': errors.name }"
                />
                <span v-if="errors.name" class="error-msg">{{ errors.name }}</span>
              </div>

              <div class="form-group">
                <label>{{ t('dashboard.createDialog.packageName') }}</label>
                <input 
                  type="text" 
                  v-model="newProject.packageName" 
                  placeholder="com.rena.myapp"
                  :class="{ 'error': errors.packageName }"
                />
                <span v-if="errors.packageName" class="error-msg">{{ errors.packageName }}</span>
              </div>

              <div class="form-group">
                <label>{{ t('dashboard.createDialog.platform') }}</label>
                <div class="platform-selector">
                  <div class="platform-option active">
                    <Smartphone :size="20" />
                    <span>Android</span>
                  </div>
                  <div class="platform-option disabled" title="Coming soon">
                    <Globe :size="20" />
                    <span>Web</span>
                  </div>
                  <div class="platform-option disabled" title="Coming soon">
                    <Layout :size="20" />
                    <span>iOS</span>
                  </div>
                </div>
              </div>

              <div class="form-group">
                <label>{{ t('dashboard.createDialog.themeColor') }}</label>
                <div class="color-grid">
                  <div 
                    v-for="c in colorPresets" 
                    :key="c"
                    class="color-option"
                    :style="{ background: c }"
                    :class="{ active: newProject.color === c }"
                    @click="newProject.color = c"
                  ></div>
                </div>
              </div>
            </div>

            <div class="modal-footer">
              <button class="btn-secondary" @click="closeCreateDialog">{{ t('dashboard.createDialog.cancel') }}</button>
              <button class="btn-primary" @click="handleCreateProject">{{ t('dashboard.createDialog.submit') }}</button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Delete Confirmation Dialog -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showDeleteDialog" class="modal-overlay" @click.self="closeDeleteDialog">
          <div class="modal-content delete-modal">
            <div class="modal-header">
              <h2 class="text-danger">Delete Project</h2>
              <button class="close-btn" @click="closeDeleteDialog">&times;</button>
            </div>
            
            <div class="modal-body">
              <div class="warning-icon">⚠️</div>
              <p>This action <strong>cannot</strong> be undone. This will permanently delete the project <strong>{{ projectToDelete?.name }}</strong>.</p>
              
              <div class="form-group mt-4">
                <label>Please type <strong>{{ projectToDelete?.name }}</strong> to confirm:</label>
                <input 
                  type="text" 
                  v-model="deleteConfirmationText" 
                  :placeholder="projectToDelete?.name"
                  class="confirm-input"
                />
              </div>
            </div>

            <div class="modal-footer">
              <button class="btn-secondary" @click="closeDeleteDialog">Cancel</button>
              <button 
                class="btn-danger" 
                :disabled="deleteConfirmationText !== projectToDelete?.name"
                @click="confirmDelete"
              >
                Delete Project
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Logout Confirmation Dialog -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showLogoutDialog" class="modal-overlay" @click.self="closeLogoutDialog">
          <div class="modal-content logout-modal">
            <div class="modal-header">
              <h2>Confirm Logout</h2>
              <button class="close-btn" @click="closeLogoutDialog">&times;</button>
            </div>
            
            <div class="modal-body text-center">
              <div class="warning-icon">👋</div>
              <p>Are you sure you want to log out of your account?</p>
            </div>

            <div class="modal-footer">
              <button class="btn-secondary" @click="closeLogoutDialog">Cancel</button>
              <button class="btn-primary" @click="handleLogout">Yes, Logout</button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.dashboard-wrapper {
  display: flex;
  height: 100vh;
  background: var(--bg);
  overflow: hidden;
}

/* Side Navigation (The "Rare" Look) */
.side-nav {
  width: 260px;
  background: var(--surface);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 24px;
  z-index: 100;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 22px;
  font-weight: 800;
  margin-bottom: 40px;
  cursor: pointer;
}

.logo-img {
  width: 32px;
  height: 32px;
  border-radius: 8px;
}

.logo-text {
  margin-left: 4px;
  letter-spacing: 0.5px;
}

.main-menu {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 12px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
  width: 100%;
}

.expand-icon {
  margin-left: auto;
  transition: transform 0.3s ease;
}

.expand-icon.rotated {
  transform: rotate(-90deg);
}

.menu-item.active {
  background: var(--primary-light);
  color: var(--primary);
}

.sidebar-section {
  margin-top: 8px;
}

.sidebar-project-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.sidebar-project-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 16px;
  border-radius: 8px;
  color: var(--text);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.sidebar-project-item:hover {
  background: var(--bg-alt);
  color: var(--primary);
}

.sidebar-project-item span {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav-bottom {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--bg-alt);
  border-radius: 16px;
}

.avatar {
  width: 32px;
  height: 32px;
  background: var(--primary);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  flex-shrink: 0;
}

.user-details {
  flex: 1;
  min-width: 0;
}

.email-truncate {
  display: block;
  font-size: 13px;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.logout-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 4px;
}

.logout-btn:hover {
  color: var(--secondary);
}

/* Main Content Area */
.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 32px 40px;
  overflow: hidden;
}

.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 40px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--surface);
  border: 1px solid var(--border);
  padding: 10px 20px;
  border-radius: 14px;
  width: 100%;
  max-width: 400px;
  color: var(--text-muted);
}

.search-box input {
  background: none;
  border: none;
  color: var(--text);
  width: 100%;
  font-weight: 500;
}

.search-box input:focus { outline: none; }

.actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-top-profile {
  display: flex;
  align-items: center;
  margin-left: 8px;
}

.top-avatar {
  width: 40px;
  height: 40px;
  background: var(--primary);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.top-avatar:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(0,0,0,0.15);
}

/* Custom Tooltip Styles */
.tooltip {
  position: relative;
}

.tooltip-content {
  position: absolute;
  bottom: 125%;
  left: 50%;
  transform: translateX(-50%) translateY(10px);
  background: #333;
  color: white;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 12px;
  white-space: nowrap;
  opacity: 0;
  visibility: hidden;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
  z-index: 1000;
}

.tooltip-content::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  margin-left: -5px;
  border-width: 5px;
  border-style: solid;
  border-color: #333 transparent transparent transparent;
}

.tooltip:hover .tooltip-content {
  opacity: 1;
  visibility: visible;
  transform: translateX(-50%) translateY(0);
}

.tooltip-name {
  font-weight: 700;
  font-size: 13px;
}

.tooltip-email {
  opacity: 0.8;
  font-size: 11px;
}

.user-top-profile .tooltip-content {
  bottom: auto;
  top: 125%;
}

.user-top-profile .tooltip-content::after {
  top: auto;
  bottom: 100%;
  border-color: transparent transparent #333 transparent;
}

.view-toggle {
  display: flex;
  background: var(--bg-alt);
  padding: 4px;
  border-radius: 10px;
}

.view-toggle button {
  background: none;
  border: none;
  padding: 6px 10px;
  border-radius: 8px;
  cursor: pointer;
  color: var(--text-muted);
}

.view-toggle button.active {
  background: var(--surface);
  color: var(--primary);
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.btn-primary {
  background: var(--primary);
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 14px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  box-shadow: 0 4px 15px rgba(46, 96, 255, 0.2);
}

/* Projects Display */
.projects-scroll {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 20px;
}

.projects-container.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

.projects-container.list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.project-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 20px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  transition: all 0.3s;
}

.project-card.clickable {
  cursor: pointer;
}

.project-card:hover {
  border-color: var(--primary);
  transform: translateY(-4px);
  box-shadow: var(--card-shadow);
}

.card-top {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 24px;
}

.project-icon {
  width: 48px;
  height: 48px;
  background: var(--primary-light);
  color: var(--primary);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.project-icon-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.project-info h3 {
  font-size: 18px;
  margin: 0 0 4px;
  font-weight: 700;
}

.date {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-muted);
}

.card-actions {
  display: flex;
  gap: 12px;
}

.delete-btn {
  padding: 10px;
  background: var(--bg-alt);
  border: 1px solid var(--border);
  border-radius: 12px;
  color: var(--text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.delete-btn:hover {
  background: #fee2e2;
  color: #ef4444;
  border-color: #fca5a5;
}

.card-actions .open-btn {
  flex: 1;
  padding: 10px;
  background: var(--bg-alt);
  border: 1px solid var(--border);
  border-radius: 12px;
  color: var(--text);
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
}

.project-card:hover .open-btn {
  background: var(--primary);
  color: white;
  border-color: var(--primary);
}

/* Modal / Dialog Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.modal-content {
  background: var(--surface);
  width: 100%;
  max-width: 500px;
  border-radius: 24px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.2);
  border: 1px solid var(--border);
  overflow: hidden;
  animation: modalSlideUp 0.3s ease-out;
}

@keyframes modalSlideUp {
  from { transform: translateY(30px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

.modal-header {
  padding: 24px 32px;
  border-bottom: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  font-size: 20px;
  font-weight: 800;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 28px;
  color: var(--text-muted);
  cursor: pointer;
  line-height: 1;
}

.modal-body {
  padding: 32px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 700;
  color: var(--text);
}

.form-group input {
  padding: 12px 16px;
  border-radius: 12px;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text);
  font-weight: 500;
  transition: all 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 4px var(--primary-light);
}

.form-group input.error {
  border-color: #ef4444;
}

.error-msg {
  font-size: 12px;
  color: #ef4444;
  font-weight: 500;
}

.platform-selector {
  display: flex;
  gap: 12px;
}

.platform-option {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid var(--border);
  background: var(--bg-alt);
  cursor: pointer;
  transition: all 0.2s;
  color: var(--text-muted);
}

.platform-option.active {
  border-color: var(--primary);
  background: var(--primary-light);
  color: var(--primary);
}

.platform-option.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.color-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 10px;
}

.color-option {
  aspect-ratio: 1;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid transparent;
  transition: transform 0.2s;
}

.color-option:hover {
  transform: scale(1.1);
}

.color-option.active {
  transform: scale(1.1);
  box-shadow: 0 0 0 3px var(--surface), 0 0 0 5px var(--brand);
}

.modal-footer {
  padding: 24px 32px;
  background: var(--bg-alt);
  display: flex;
  justify-content: flex-end;
  gap: 16px;
}

.btn-secondary {
  padding: 12px 24px;
  border-radius: 12px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text);
  font-weight: 700;
  cursor: pointer;
}

/* Delete Modal Specific Styles */
.delete-modal {
  max-width: 450px !important;
}

.text-danger {
  color: #ff4d4d;
}

.warning-icon {
  font-size: 48px;
  text-align: center;
  margin-bottom: 16px;
}

.btn-danger {
  background: #ff4d4d;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-danger:hover:not(:disabled) {
  background: #ff3333;
  transform: translateY(-1px);
}

.btn-danger:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.mt-4 {
  margin-top: 1.5rem;
}

.confirm-input {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--bg);
  color: var(--text);
  margin-top: 8px;
  box-sizing: border-box;
}

.confirm-input:focus {
  border-color: #ff4d4d;
  outline: none;
}

/* Logout Modal Specific Styles */
.logout-modal {
  max-width: 400px !important;
}

.text-center {
  text-align: center;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Modal Animations */
.modal-content {
  animation: modalIn 0.3s ease-out;
}

@keyframes modalIn {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* List View Overrides */
.list .project-card {
  flex-direction: row;
  align-items: center;
  padding: 16px 24px;
}
.list .card-top { margin-bottom: 0; flex: 1; }
.list .card-actions .open-btn { width: auto; padding: 10px 20px; }

/* Empty/Loading States */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 80px 0;
}

.empty-icon { font-size: 64px; margin-bottom: 24px; }

.empty-state h3 {
  margin-bottom: 24px;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--border);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s infinite linear;
  margin-bottom: 16px;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* RTL Support */
[dir="rtl"] .side-nav { border-right: none; border-left: 1px solid var(--border); }
[dir="rtl"] .menu-item { text-align: right; }
[dir="rtl"] .lang-dropdown { right: auto; left: 0; }
.rtl-flip { transition: transform 0.3s; }
[dir="rtl"] .rtl-flip { transform: rotate(180deg); }

/* Language Dropdown in Sidebar */
.lang-switcher { position: relative; }
.nav-icon-btn {
  background: none;
  border: 1px solid var(--border);
  color: var(--text-muted);
  padding: 8px;
  border-radius: 10px;
  cursor: pointer;
  display: flex;
}
.lang-dropdown {
  position: absolute;
  bottom: 100%;
  left: 0;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 8px;
  box-shadow: var(--card-shadow);
  min-width: 140px;
  margin-bottom: 8px;
}
.lang-dropdown button {
  width: 100%;
  padding: 8px 12px;
  border: none;
  background: none;
  text-align: left;
  color: var(--text);
  cursor: pointer;
  border-radius: 8px;
}
.lang-dropdown button:hover { background: var(--bg-alt); color: var(--primary); }
</style>

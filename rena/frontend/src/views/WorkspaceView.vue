<script setup lang="ts">
import { ref, onMounted, onUnmounted, shallowRef, markRaw } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as projectsApi from '../services/projects'
import { 
  ChevronLeft, Play, Save, Settings, 
  Undo2, Redo2, Moon, Sun, Download, Upload
} from 'lucide-vue-next'
import * as Blockly from 'blockly'
import DarkTheme from '@blockly/theme-dark'
import iconImg from '../assets/icon.png'
import { mainToolbox, ReactNativeGenerator, RenaTheme, initializeTheme } from '../blocks'
import PhonePreview from '../components/PhonePreview.vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const projectId = route.params.id as string
const project = ref<any>(null)
const loading = ref(true)
const blocklyDiv = ref<HTMLElement | null>(null)
const workspace = shallowRef<Blockly.WorkspaceSvg | null>(null)
const generatedCode = ref('')
const isDarkMode = ref(false)
const audioContext = ref<AudioContext | null>(null)
const isExporting = ref(false)
const isImporting = ref(false)
const isSaving = ref(false)
const saveStatus = ref('')
const saveError = ref('')
const hasUnsavedChanges = ref(false)
let autoSaveTimer: number | null = null

// Sidebar state (removed - no longer needed)

const fetchProject = async () => {
  const result = await projectsApi.getProject(projectId)
  if (result.error || !result.data?.project) {
    console.error('Error fetching project:', result.error)
    router.push('/dashboard')
  } else {
    project.value = result.data.project
  }
  loading.value = false
}

const loadWorkspaceFromProject = () => {
  if (!workspace.value || !project.value || !project.value.workspace_xml) return

  try {
    const xml = (Blockly.Xml as any).textToDom(project.value.workspace_xml)
    Blockly.Xml.domToWorkspace(xml, workspace.value)
    hasUnsavedChanges.value = false
    saveStatus.value = 'Workspace loaded'
  } catch (err) {
    console.warn('Could not load saved workspace XML:', err)
    saveStatus.value = 'Workspace ready'
  }
}

const initBlockly = () => {
  if (!blocklyDiv.value) return

  const ws = Blockly.inject(blocklyDiv.value, {
    toolbox: mainToolbox,
    theme: RenaTheme,
    grid: {
      spacing: 20,
      length: 3,
      colour: '#ccc',
      snap: true
    },
    zoom: {
      controls: true,
      wheel: true,
      startScale: 1.0,
      maxScale: 3,
      minScale: 0.3,
      scaleSpeed: 1.2
    },
    trashcan: true,
    renderer: 'geras' // Use default geras renderer with custom theme
  })
  
  workspace.value = markRaw(ws)
  initializeTheme(ws)
  setupAudioFeedback()
  setupBlockValidation()

  loadWorkspaceFromProject()
  setupAutoSave()

  // Handle resize
  window.addEventListener('resize', () => {
    if (workspace.value) Blockly.svgResize(workspace.value as Blockly.WorkspaceSvg)
  })
}

onMounted(async () => {
  await fetchProject()
  initBlockly()
})

onUnmounted(() => {
  if (workspace.value) {
    workspace.value.dispose()
  }
  if (autoSaveTimer) {
    clearTimeout(autoSaveTimer)
  }
  // Clean up audio context
  if (audioContext.value && audioContext.value.state !== 'closed') {
    audioContext.value.close()
  }
})

const scheduleAutoSave = () => {
  if (autoSaveTimer) {
    clearTimeout(autoSaveTimer)
  }
  autoSaveTimer = window.setTimeout(async () => {
    if (hasUnsavedChanges.value) {
      await saveProject(true)
    }
  }, 1500)
}

const setupAutoSave = () => {
  if (!workspace.value) return
  workspace.value.addChangeListener(() => {
    hasUnsavedChanges.value = true
    saveStatus.value = 'Unsaved changes'
    scheduleAutoSave()
  })
}

const saveProject = async (isAutoSave = false) => {
  if (!workspace.value || !project.value) return

  if (!isAutoSave) {
    saveStatus.value = 'Saving...'
  }

  isSaving.value = true
  saveError.value = ''

  try {
    const xml = Blockly.Xml.workspaceToDom(workspace.value)
    const xmlText = Blockly.Xml.domToText(xml)
    const code = ReactNativeGenerator.workspaceToCode(workspace.value)
    const response = await projectsApi.saveWorkspace(projectId, xmlText, code)

    if (response.error) {
      saveError.value = response.error
      saveStatus.value = 'Save failed'
      console.error('Failed to save project:', response.error)
    } else {
      hasUnsavedChanges.value = false
      saveStatus.value = isAutoSave ? 'Auto-saved' : 'Saved successfully'
      if (!isAutoSave) {
        alert('Project saved successfully!')
      }
    }
  } catch (err) {
    saveError.value = err instanceof Error ? err.message : 'Unknown save error'
    saveStatus.value = 'Save failed'
    console.error('Error saving project:', err)
    if (!isAutoSave) {
      alert('Error saving project: ' + saveError.value)
    }
  } finally {
    isSaving.value = false
  }
}

const runProject = () => {
  if (!workspace.value) return
  
  try {
    const code = ReactNativeGenerator.workspaceToCode(workspace.value)
    generatedCode.value = code
    console.log('Generated React Native Code:\n', code)
    alert('Code Generated! Check Console.')
  } catch (err) {
    console.error('Code generation failed:', err)
  }
}

const exportProject = async () => {
  if (!workspace.value || !project.value) return
  
  isExporting.value = true
  
  try {
    // First, save the current workspace state
    const xml = Blockly.Xml.workspaceToDom(workspace.value)
    const xmlText = Blockly.Xml.domToText(xml)
    
    // Generate the code
    const code = ReactNativeGenerator.workspaceToCode(workspace.value)
    
    // Save to database first
    const saveResult = await projectsApi.saveWorkspace(projectId, xmlText, code)
    if (saveResult.error) {
      console.error('Error saving workspace before export:', saveResult.error)
      throw new Error('Failed to save workspace')
    }
    
    // Now download from backend
    const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'
    const userId = useAuthStore().user?.id
    if (!userId) {
      throw new Error('User not authenticated')
    }

    const response = await fetch(`${backendUrl}/api/projects/${projectId}/export`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/zip',
        'X-User-ID': userId,
      },
    })
    
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Export failed')
    }
    
    // Create blob and download
    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${project.value.name}.rnp`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    console.log('Project exported successfully')
    alert('Project exported successfully!')
  } catch (error) {
    console.error('Export failed:', error)
    alert('Failed to export project: ' + (error as Error).message)
  } finally {
    isExporting.value = false
  }
}

const importProject = async () => {
  // Create file input
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.rnp'
  
  input.onchange = async (e: Event) => {
    const target = e.target as HTMLInputElement
    const file = target.files?.[0]
    
    if (!file) return
    
    isImporting.value = true
    
    try {
      const formData = new FormData()
      formData.append('file', file)
      const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'
      const token = localStorage.getItem('auth_token')
      const headers: Record<string, string> = {}
      if (token) headers.Authorization = `Bearer ${token}`

      const response = await fetch(`${backendUrl}/api/projects/import`, {
        method: 'POST',
        headers,
        body: formData,
      })
      
      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Import failed')
      }
      
      const result = await response.json()
      console.log('Import successful:', result)
      alert(`Project imported successfully as "${result.project_name}"!`)
      
      // Navigate to projects view
      router.push('/dashboard')
    } catch (error) {
      console.error('Import failed:', error)
      alert('Failed to import project: ' + (error as Error).message)
    } finally {
      isImporting.value = false
    }
  }
  
  input.click()
}

const toggleDarkMode = () => {
  if (!workspace.value) return
  
  isDarkMode.value = !isDarkMode.value
  
  // Apply the appropriate theme
  if (isDarkMode.value) {
    workspace.value.setTheme(DarkTheme)
  } else {
    workspace.value.setTheme(RenaTheme)
  }
}

// Audio feedback functions
const initAudioContext = () => {
  if (!audioContext.value) {
    audioContext.value = new (window.AudioContext || (window as any).webkitAudioContext)()
  }
}

const playTone = (frequency: number, duration: number, type: OscillatorType = 'sine') => {
  if (!audioContext.value) return
  
  const oscillator = audioContext.value.createOscillator()
  const gainNode = audioContext.value.createGain()
  
  oscillator.connect(gainNode)
  gainNode.connect(audioContext.value.destination)
  
  oscillator.frequency.value = frequency
  oscillator.type = type
  
  gainNode.gain.setValueAtTime(0.1, audioContext.value.currentTime)
  gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.value.currentTime + duration)
  
  oscillator.start(audioContext.value.currentTime)
  oscillator.stop(audioContext.value.currentTime + duration)
}

const playClickSound = () => {
  playTone(800, 0.05, 'triangle')
}

const playDeleteSound = () => {
  playTone(400, 0.1, 'sawtooth')
}

const playConnectSound = () => {
  playTone(1200, 0.08, 'sine')
}

const playDisconnectSound = () => {
  playTone(600, 0.06, 'triangle')
}

const setupAudioFeedback = () => {
  initAudioContext()
  
  if (!workspace.value) return
  
  // Initialize audio on first user interaction
  const initAudioOnInteraction = () => {
    initAudioContext()
    if (audioContext.value?.state === 'suspended') {
      audioContext.value.resume()
    }
    document.removeEventListener('click', initAudioOnInteraction)
  }
  
  document.addEventListener('click', initAudioOnInteraction)
  
  // Hook into Blockly events for audio feedback
  workspace.value.addChangeListener((event: any) => {
    switch(event.type) {
      case Blockly.Events.CLICK:
        playClickSound()
        break
      case Blockly.Events.BLOCK_DELETE:
        playDeleteSound()
        break
      case Blockly.Events.BLOCK_MOVE:
        if (event.newParentId) {
          playConnectSound()
        } else if (event.oldParentId) {
          playDisconnectSound()
        }
        break
    }
  })
}

// Block validation and error/warning indicators
const setupBlockValidation = () => {
  if (!workspace.value) return
  
  // Listen for block change events
  workspace.value.addChangeListener((event: any) => {
    if (event.type === Blockly.Events.BLOCK_CHANGE) {
      validateBlock(event.blockId)
    } else if (event.type === Blockly.Events.BLOCK_CREATE) {
      validateBlock(event.blockId)
    }
  })
}

const validateBlock = (blockId: string) => {
  if (!workspace.value) return
  
  const block = workspace.value.getBlockById(blockId)
  if (!block) return
  
  // Check for warnings and errors
  const hasError = checkBlockForErrors(block)
  const hasWarning = checkBlockForWarnings(block)
  
  // Update block indicators
  if (hasError) {
    block.setWarningText('Error: Invalid configuration')
  } else if (hasWarning) {
    block.setWarningText('Warning: Check configuration')
  } else {
    block.setWarningText(null)
  }
}

const checkBlockForErrors = (block: Blockly.Block): boolean => {
  // Add your custom error checking logic here
  // Example: Check if required fields are filled
  // For now, return false by default
  console.log('Checking block for errors:', block.type)
  return false
}

const checkBlockForWarnings = (block: Blockly.Block): boolean => {
  // Add your custom warning checking logic here
  // Example: Check for deprecated blocks or suboptimal configurations
  // For now, return false by default
  console.log('Checking block for warnings:', block.type)
  return false
}
</script>

<template>
  <div class="workspace-wrapper" v-if="!loading">
    <!-- Top Header -->
    <header class="workspace-header">
      <div class="header-left">
        <button class="back-btn" @click="router.push('/dashboard')">
          <ChevronLeft :size="20" />
        </button>
        <div class="project-meta">
          <img :src="iconImg" class="mini-logo" />
          <div class="meta-text">
            <span class="project-name">{{ project?.name }}</span>
            <span class="package-name">{{ project?.package_name }}</span>
          </div>
        </div>
      </div>

      <div class="header-center">
        <!-- Removed view switcher -->
      </div>

      <div class="header-right">
        <div class="workspace-actions">
          <button class="tool-btn" title="Undo" @click="workspace?.undo(false)"><Undo2 :size="18" /></button>
          <button class="tool-btn" title="Redo" @click="workspace?.undo(true)"><Redo2 :size="18" /></button>
          <div class="divider"></div>
          <button class="action-btn save" @click="saveProject(false)" :disabled="isSaving">
            <Save :size="18" />
            <span>{{ isSaving ? 'Saving...' : 'Save' }}</span>
          </button>
          <button class="action-btn run" @click="runProject">
            <Play :size="18" />
            <span>Run</span>
          </button>
          <div class="divider"></div>
          <button 
            class="action-btn export" 
            @click="exportProject" 
            :disabled="isExporting"
            title="Export Project as .rnp file"
          >
            <Download :size="18" />
            <span>{{ isExporting ? 'Exporting...' : 'Export' }}</span>
          </button>
          <button 
            class="action-btn import" 
            @click="importProject" 
            :disabled="isImporting"
            title="Import Project from .rnp file"
          >
            <Upload :size="18" />
            <span>{{ isImporting ? 'Importing...' : 'Import' }}</span>
          </button>
          <div class="divider"></div>
          <button class="tool-btn" :title="isDarkMode ? 'Light Mode' : 'Dark Mode'" @click="toggleDarkMode">
            <Moon v-if="!isDarkMode" :size="18" />
            <Sun v-else :size="18" />
          </button>
          <button class="tool-btn settings"><Settings :size="20" /></button>
        </div>
        <div class="save-status-row">
          <span class="save-status">{{ saveStatus }}</span>
          <span class="save-error" v-if="saveError">{{ saveError }}</span>
        </div>
      </div>
    </header>

    <div class="main-workspace">
      <!-- Blockly Editor -->
      <div class="editor-area">
        <div ref="blocklyDiv" class="blockly-container"></div>
      </div>

      <!-- Phone Preview Component -->
      <PhonePreview :project="project" />
    </div>
  </div>

  <div v-else class="loading-workspace">
    <div class="spinner"></div>
    <p>Loading Workspace...</p>
  </div>
</template>

<style scoped>
.workspace-wrapper {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--bg);
  overflow: hidden;
}

.workspace-header {
  height: 64px;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  z-index: 100;
}

.header-left, .header-right {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
}

.header-right { justify-content: flex-end; }

.back-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  transition: all 0.2s;
}

.back-btn:hover {
  background: var(--bg-alt);
  color: var(--text);
}

.project-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.mini-logo {
  width: 32px;
  height: 32px;
  border-radius: 6px;
}

.meta-text {
  display: flex;
  flex-direction: column;
}

.project-name {
  font-weight: 700;
  font-size: 14px;
  color: var(--text);
}

.package-name {
  font-size: 11px;
  color: var(--text-muted);
}

.view-switcher {
  display: flex;
  background: var(--bg-alt);
  padding: 4px;
  border-radius: 10px;
}

.view-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 16px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.view-btn.active {
  background: var(--surface);
  color: var(--primary);
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.workspace-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
.save-status-row {
  margin-left: 16px;
  color: var(--text-muted);
  font-size: 11px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.save-status {
  color: var(--text-muted);
}
.save-error {
  color: #ff6b6b;
}

.tool-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
}

.tool-btn:hover { color: var(--text); background: var(--bg-alt); }

.divider {
  width: 1px;
  height: 24px;
  background: var(--border);
  margin: 0 4px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  border: none;
  font-weight: 700;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn.save {
  background: var(--bg-alt);
  color: var(--text);
}

.action-btn.run {
  background: var(--primary);
  color: white;
}

.action-btn.export {
  background: var(--success);
  color: white;
}

.action-btn.import {
  background: var(--info);
  color: white;
}

.action-btn:hover { transform: translateY(-1px); }

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* Main Area */
.main-workspace {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* Editor Area - Takes all available space */
.editor-area {
  flex: 1;
  position: relative;
  min-width: 0; /* Allows flex items to shrink below content size */
}

.blockly-container {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.app-header {
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  color: white;
}

.app-title { font-weight: 700; font-size: 16px; }

.app-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
}

.non-visible-tray {
  height: 40px;
  background: #f8f9fa;
  border-top: 1px solid #eee;
  display: flex;
  align-items: center;
  padding: 0 12px;
  gap: 12px;
}

.tray-icon {
  width: 24px;
  height: 24px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
}

.loading-workspace {
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
}
</style>

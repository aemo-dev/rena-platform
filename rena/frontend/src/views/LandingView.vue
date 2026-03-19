<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import gsap from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { 
  Globe, ArrowRight, Layers, 
  ShieldCheck, Languages, Smartphone, Code, 
  Cloud, Download, Share2, Palette, Zap, Sparkles, Box
} from 'lucide-vue-next'
import iconImg from '../assets/icon.png'
import { setLanguage } from '../i18n'

gsap.registerPlugin(ScrollTrigger)

const router = useRouter()
const { t } = useI18n()
const auth = useAuthStore()
const showLangMenu = ref(false)
const isScrolled = ref(false)

const userDisplayName = computed(() => {
  return auth.user?.user_metadata?.full_name || auth.user?.user_metadata?.name || auth.user?.email?.split('@')[0] || 'User'
})

const userInitial = computed(() => {
  return userDisplayName.value.charAt(0).toUpperCase()
})

const changeLanguage = (lang: string) => {
  setLanguage(lang)
  showLangMenu.value = false
}

const floatingIcons = [Zap, Sparkles, Box, Layers, Code, Globe]
const heroShapes = ref<any[]>([])

onMounted(async () => {
  document.title = 'Rena'
  // Generate hero background shapes
  for (let i = 0; i < 15; i++) {
    heroShapes.value.push({
      id: i,
      icon: floatingIcons[Math.floor(Math.random() * floatingIcons.length)],
      size: Math.random() * 30 + 15,
      top: `${Math.random() * 80 + 10}%`,
      left: `${Math.random() * 90 + 5}%`,
      opacity: Math.random() * 0.1 + 0.05,
      duration: Math.random() * 5 + 5,
      delay: Math.random() * 2
    })
  }

  // Handle scroll for navbar
  const handleScroll = () => {
    isScrolled.value = window.scrollY > 50
  }
  window.addEventListener('scroll', handleScroll)
  handleScroll() // Initial check

  // Ensure auth state is loaded
  if (auth.loading) {
    await auth.fetchUser()
  }

  // Animate hero shapes
  setTimeout(() => {
    const elShapes = document.querySelectorAll('.hero-shape')
    elShapes.forEach((el) => {
      gsap.to(el, {
        y: 'random(-60, 60)',
        x: 'random(-40, 40)',
        rotation: 'random(-90, 90)',
        duration: 'random(4, 8)',
        repeat: -1,
        yoyo: true,
        ease: 'sine.inOut'
      })
    })
  }, 100)

  // Hero Animation
  gsap.from(".hero-content > *", {
    y: 30,
    opacity: 0,
    duration: 1,
    stagger: 0.2,
    ease: "power4.out"
  })

  // Reveal Animations on Scroll
  const revealSections = [".reveal-up", ".reveal-left", ".reveal-right"]
  revealSections.forEach(selector => {
    gsap.utils.toArray(selector).forEach((el: any) => {
      gsap.from(el, {
        scrollTrigger: {
          trigger: el,
          start: "top 85%",
          toggleActions: "play none none none"
        },
        y: selector === ".reveal-up" ? 50 : 0,
        x: selector === ".reveal-left" ? -50 : (selector === ".reveal-right" ? 50 : 0),
        opacity: 0,
        duration: 1,
        ease: "power3.out"
      })
    })
  })
})

const startNow = () => {
  if (auth.user) {
    router.push('/dashboard')
  } else {
    auth.fetchUser().then(() => {
      router.push('/dashboard')
    })
  }
}
</script>

<template>
  <div class="landing">
    <!-- Grid Background -->
    <div class="grid-background"></div>

    <!-- Navbar -->
    <nav class="navbar" :class="{ 'scrolled': isScrolled }">
      <div class="nav-container">
        <div class="logo" @click="router.push('/')">
          <img :src="iconImg" alt="Rena Logo" class="logo-img" />
          <span class="logo-text">Rena Builder</span>
        </div>
        
        <div class="nav-links">
          <div class="lang-switcher">
            <button @click="showLangMenu = !showLangMenu" class="nav-btn-icon">
              <Languages :size="20" />
            </button>
            <div v-if="showLangMenu" class="lang-dropdown">
              <button @click="changeLanguage('en')">English</button>
              <button @click="changeLanguage('ar')">العربية</button>
            </div>
          </div>
          
          <template v-if="!auth.user">
            <button @click="startNow" class="login-btn">{{ t('nav.login') }}</button>
            <button @click="startNow" class="cta-nav-btn">{{ t('nav.start') }}</button>
          </template>
          <template v-else>
            <div class="user-nav-profile tooltip" v-if="auth.user">
              <div class="user-avatar" @click="router.push('/dashboard')">
                {{ userInitial }}
              </div>
              <div class="tooltip-content">
                <span class="tooltip-name">{{ userDisplayName }}</span>
                <span class="tooltip-email">{{ auth.user.email }}</span>
              </div>
            </div>
            <button @click="router.push('/dashboard')" class="cta-nav-btn">{{ t('dashboard.title') }}</button>
          </template>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <header class="hero">
      <div class="hero-bg-shapes">
        <div 
          v-for="shape in heroShapes" 
          :key="shape.id"
          class="hero-shape"
          :style="{
            top: shape.top,
            left: shape.left,
            width: shape.size + 'px',
            height: shape.size + 'px',
            opacity: shape.opacity,
            color: 'var(--primary)'
          }"
        >
          <component :is="shape.icon" :size="shape.size" />
        </div>
      </div>
      <div class="hero-content">
        <h1 class="hero-title">
          {{ t('hero.titlePrefix') }} <span class="text-primary">{{ t('hero.titleAccent') }}</span> {{ t('hero.titleSuffix') }}
        </h1>
        <div class="hero-features">
          <span><Layers :size="18" /> {{ t('hero.feature1') }}</span>
          <span><Code :size="18" /> {{ t('hero.feature2') }}</span>
          <span><Globe :size="18" /> {{ t('hero.feature3') }}</span>
        </div>
        <button @click="startNow" class="hero-cta-btn">
          {{ t('hero.cta') }} <ArrowRight :size="20" class="rtl-flip" />
        </button>
        <div class="scroll-indicator">
          <span>{{ t('hero.scroll') }}</span>
          <div class="mouse"></div>
        </div>
      </div>
    </header>

    <!-- Visual Designer Section -->
    <section class="info-section bg-alt">
      <div class="container grid-2">
        <div class="text-content reveal-left">
          <span class="section-badge">{{ t('sections.designer.badge') }}</span>
          <h2>{{ t('sections.designer.title') }}</h2>
          <p>{{ t('sections.designer.desc') }}</p>
          <ul class="feature-list">
            <li><div class="dot"></div> {{ t('sections.designer.feat1') }}</li>
            <li><div class="dot"></div> {{ t('sections.designer.feat2') }}</li>
            <li><div class="dot"></div> {{ t('sections.designer.feat3') }}</li>
          </ul>
        </div>
        <div class="visual-content reveal-right">
          <div class="designer-preview">
            <div class="palette-mock"></div>
            <div class="canvas-mock">
              <div class="mock-button"></div>
              <div class="mock-card"></div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Visual Blocks Section -->
    <section class="info-section">
      <div class="container grid-2 reverse">
        <div class="text-content reveal-right">
          <span class="section-badge">{{ t('sections.blocks.badge') }}</span>
          <h2>{{ t('sections.blocks.title') }}</h2>
          <p>{{ t('sections.blocks.desc') }}</p>
          <ul class="feature-list">
            <li><div class="dot"></div> {{ t('sections.blocks.feat1') }}</li>
            <li><div class="dot"></div> {{ t('sections.blocks.feat2') }}</li>
            <li><div class="dot"></div> {{ t('sections.blocks.feat3') }}</li>
          </ul>
        </div>
        <div class="visual-content reveal-left">
          <div class="blocks-preview">
            <div class="block b-1">When Button.Click</div>
            <div class="block b-2">Do: Play Sound</div>
            <div class="block b-3">Store Data: "Score"</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Project Management & Build Section -->
    <section class="info-section bg-alt">
      <div class="container grid-2">
        <div class="text-content reveal-left">
          <span class="section-badge">{{ t('sections.management.badge') }}</span>
          <h2>{{ t('sections.management.title') }}</h2>
          <p>{{ t('sections.management.desc') }}</p>
          <div class="btn-group">
            <div class="mini-feat"><Cloud :size="20" /> {{ t('sections.management.cloud') }}</div>
            <div class="mini-feat"><Smartphone :size="20" /> {{ t('sections.management.access') }}</div>
          </div>
        </div>
        <div class="text-content reveal-right">
          <span class="section-badge">{{ t('sections.export.badge') }}</span>
          <h2>{{ t('sections.export.title') }}</h2>
          <p>{{ t('sections.export.desc') }}</p>
          <div class="btn-group">
            <div class="mini-feat"><Download :size="20" /> {{ t('sections.export.oneclick') }}</div>
            <div class="mini-feat"><Share2 :size="20" /> {{ t('sections.export.qr') }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Why Choose Section -->
    <section class="why-choose">
      <div class="container">
        <div class="section-header reveal-up">
          <h2>{{ t('features.title') }}</h2>
          <p>{{ t('features.sub') }}</p>
        </div>
        <div class="features-grid reveal-up">
          <div v-for="(feat, i) in whyChoose" :key="i" class="why-card">
            <div class="why-icon" :style="{ background: feat.color + '15', color: feat.color }">
              <component :is="feat.icon" :size="24" />
            </div>
            <h4>{{ t(`features.f${i+1}.title`) }}</h4>
            <p>{{ t(`features.f${i+1}.desc`) }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Final CTA -->
    <section class="final-cta">
      <div class="container reveal-up">
        <h2>{{ t('cta.title') }}</h2>
        <p>{{ t('cta.sub') }}</p>
        <button @click="startNow" class="hero-cta-btn large">{{ t('cta.btn') }}</button>
      </div>
    </section>

    <!-- Footer -->
    <footer class="footer">
      <div class="container footer-grid">
        <div class="footer-brand">
          <div class="logo">
            <img :src="iconImg" alt="Rena Logo" class="logo-img" />
            <span class="logo-text">Rena Builder</span>
          </div>
          <p>{{ t('footer.desc') }}</p>
        </div>
        <div class="footer-links-group">
          <div class="link-col">
            <h5>{{ t('footer.col1') }}</h5>
            <a href="#">{{ t('footer.col1_item1') }}</a>
            <a href="#">{{ t('footer.col1_item2') }}</a>
            <a href="#">{{ t('footer.col1_item3') }}</a>
          </div>
          <div class="link-col">
            <h5>{{ t('footer.col2') }}</h5>
            <a href="#">{{ t('footer.col2_item1') }}</a>
            <a href="#">{{ t('footer.col2_item2') }}</a>
            <a href="#">{{ t('footer.col2_item3') }}</a>
          </div>
          <div class="link-col">
            <h5>{{ t('footer.col3') }}</h5>
            <a href="#">{{ t('footer.col3_item1') }}</a>
            <a href="#">{{ t('footer.col3_item2') }}</a>
            <a href="#">{{ t('footer.col3_item3') }}</a>
          </div>
        </div>
      </div>
      <div class="footer-bottom">
        <div class="container">
        <p>&copy; {{ new Date().getFullYear() }} {{ t('footer.copy') }}</p>
      </div>
      </div>
    </footer>
  </div>
</template>

<script lang="ts">
const whyChoose = [
  { title: 'No Coding!', desc: 'Create apps without writing a single line of code.', icon: Code, color: '#2e60ff' },
  { title: 'Cloud Based', desc: 'Safe storage on powerful cloud infrastructure.', icon: Cloud, color: '#4caf50' },
  { title: 'Simple Blocks', desc: 'Drag & drop blocks to build your app logic easily.', icon: Layers, color: '#ff4081' },
  { title: 'Live Testing', desc: 'Test changes instantly on your device.', icon: Smartphone, color: '#ffc107' },
  { title: 'Material Design', desc: 'Beautiful, modern interfaces built-in.', icon: Palette, color: '#3f51b5' },
  { title: 'Secure', desc: 'Enterprise-grade security for your data.', icon: ShieldCheck, color: '#00bcd4' }
]
export default { whyChoose }
</script>

<style scoped>
.landing {
  position: relative;
  background: var(--bg);
  color: var(--text);
  overflow-x: hidden;
}

.grid-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: linear-gradient(rgba(0,0,0,0.03) 1px, transparent 1px),
                    linear-gradient(to right, rgba(0,0,0,0.03) 1px, transparent 1px);
  background-size: 40px 40px;
  pointer-events: none;
  z-index: 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.bg-alt {
  background-color: var(--bg-alt);
}

.grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 80px;
  align-items: center;
}

@media (max-width: 968px) {
  .grid-2 {
    grid-template-columns: 1fr;
    gap: 40px;
    text-align: center;
  }
  .grid-2.reverse .text-content { order: 1; }
  .grid-2.reverse .visual-content { order: 2; }
}

/* Navbar */
.navbar {
  position: fixed;
  top: 0;
  width: 100%;
  background: transparent;
  backdrop-filter: blur(0px);
  border-bottom: 1px solid transparent;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1000;
  height: 88px;
  display: flex;
  align-items: center;
}

.navbar.scrolled {
  background: var(--navbar-bg);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border);
  height: 72px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.03);
}

.nav-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 40px; /* Increased side spacing */
}

.logo-text {
  margin-left: 10px;
  letter-spacing: 0.8px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 22px;
  font-weight: 800;
  cursor: pointer;
}

.logo-img {
  width: 32px;
  height: 32px;
  border-radius: 8px;
}

.nav-links {
  display: flex;
  gap: 24px;
  align-items: center;
}

.nav-btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-muted);
  padding: 8px;
  border-radius: 50%;
  display: flex;
}

.nav-btn-icon:hover {
  background: var(--surface-alt);
  color: var(--primary);
}

.login-btn {
  background: none;
  border: none;
  font-weight: 600;
  color: var(--text);
  cursor: pointer;
  padding: 8px 16px;
}

.user-nav-profile {
  position: relative;
  display: flex;
  align-items: center;
}

.user-avatar {
  width: 40px;
  height: 40px;
  background: var(--grad-brand);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.3s;
  border: 2px solid var(--surface);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.user-avatar:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(0,0,0,0.15);
}

/* Custom Tooltip Styles */
.tooltip {
  position: relative;
}

.tooltip-content {
  position: absolute;
  top: 125%;
  left: 50%;
  transform: translateX(-50%) translateY(-10px);
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
  bottom: 100%;
  left: 50%;
  margin-left: -5px;
  border-width: 5px;
  border-style: solid;
  border-color: transparent transparent #333 transparent;
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

.cta-nav-btn {
  background: var(--primary);
  color: white;
  border: none;
  padding: 12px 28px;
  border-radius: 50px; /* Fully circular */
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 15px rgba(46, 96, 255, 0.2);
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.cta-nav-btn:hover {
  background: var(--primary-dark);
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 8px 25px rgba(46, 96, 255, 0.3);
}

/* Hero */
.hero {
  min-height: 95vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  padding: 120px 24px 60px;
  text-align: center;
  background: var(--grad-surface);
}

.hero-bg-shapes {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
  pointer-events: none;
}

.hero-shape {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  filter: blur(0.5px);
}

.hero-content {
  max-width: 900px;
}

.hero-title {
  font-size: clamp(40px, 7vw, 68px);
  font-weight: 900;
  line-height: 1.15;
  margin-bottom: 24px;
  color: var(--text);
  letter-spacing: -1.5px;
}

.text-primary {
  color: var(--primary);
}

.hero-features {
  display: flex;
  justify-content: center;
  gap: 32px;
  margin-bottom: 48px;
  color: var(--text-muted);
  font-weight: 500;
  flex-wrap: wrap;
}

.hero-features span {
  display: flex;
  align-items: center;
  gap: 8px;
}

.hero-cta-btn {
  background: var(--primary);
  color: white;
  border: none;
  padding: 18px 42px;
  font-size: 18px;
  font-weight: 800;
  border-radius: 50px; /* Fully circular */
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0 auto;
  box-shadow: 0 8px 25px rgba(46, 96, 255, 0.25);
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.hero-cta-btn:hover {
  background: var(--primary-dark);
  transform: translateY(-5px) scale(1.05);
  box-shadow: 0 12px 30px rgba(46, 96, 255, 0.3);
}

.hero-cta-btn.large {
  padding: 20px 60px;
  font-size: 20px;
}

.scroll-indicator {
  margin-top: 80px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: var(--text-muted);
  font-size: 14px;
  font-weight: 600;
}

.mouse {
  width: 24px;
  height: 40px;
  border: 2px solid var(--border);
  border-radius: 12px;
  position: relative;
}

.mouse::before {
  content: '';
  position: absolute;
  top: 8px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 8px;
  background: var(--primary);
  border-radius: 2px;
  animation: scroll 2s infinite;
}

@keyframes scroll {
  0% { transform: translate(-50%, 0); opacity: 1; }
  100% { transform: translate(-50%, 15px); opacity: 0; }
}

/* Info Sections */
.info-section {
  padding: 120px 0;
}

.section-badge {
  color: var(--primary);
  font-weight: 800;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 16px;
  display: block;
}

h2 {
  font-size: clamp(32px, 5vw, 42px);
  font-weight: 800;
  line-height: 1.2;
  margin-bottom: 24px;
}

p {
  font-size: 18px;
  line-height: 1.6;
  color: var(--text-muted);
  margin-bottom: 32px;
}

.feature-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.feature-list li {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  font-weight: 600;
  color: var(--text);
}

.dot {
  width: 8px;
  height: 8px;
  background: var(--primary);
  border-radius: 50%;
}

.btn-group {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.mini-feat {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  font-weight: 700;
  font-size: 15px;
}

/* Visual Mocks */
.designer-preview {
  background: white;
  border-radius: 20px;
  box-shadow: 0 30px 60px rgba(0,0,0,0.1);
  display: flex;
  height: 400px;
  overflow: hidden;
  border: 1px solid var(--border);
}

.palette-mock {
  width: 30%;
  background: #f1f3f9;
  border-right: 1px solid #e2e8f0;
}

.canvas-mock {
  flex: 1;
  padding: 40px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.mock-button { width: 120px; height: 40px; background: #2e60ff; border-radius: 8px; }
.mock-card { width: 100%; height: 180px; background: #f8fafc; border-radius: 12px; border: 2px dashed #cbd5e1; }

.blocks-preview {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 40px;
}

.block {
  padding: 16px 24px;
  border-radius: 12px;
  font-weight: 700;
  color: white;
  width: fit-content;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.b-1 { background: #ffc107; margin-left: 0; }
.b-2 { background: #2e60ff; margin-left: 30px; }
.b-3 { background: #4caf50; margin-left: 60px; }

/* Why Choose */
.why-choose {
  padding: 120px 0;
  background: var(--bg);
}

.section-header {
  text-align: center;
  margin-bottom: 80px;
}

.why-card {
  background: var(--surface);
  padding: 40px;
  border-radius: 24px;
  border: 1px solid var(--border);
  transition: all 0.3s;
  text-align: left;
}

.why-card:hover {
  border-color: var(--primary);
  box-shadow: var(--card-shadow);
  transform: translateY(-5px);
}

.why-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
}

.why-card h4 {
  font-size: 20px;
  margin-bottom: 12px;
}

.why-card p {
  font-size: 16px;
  margin: 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 32px;
}

/* Final CTA */
.final-cta {
  padding: 120px 0;
  background: var(--primary);
  color: white;
  text-align: center;
}

.final-cta h2 { color: white; margin-bottom: 16px; }
.final-cta p { color: rgba(255,255,255,0.8); margin-bottom: 40px; }
.final-cta .hero-cta-btn { background: white; color: var(--primary); box-shadow: 0 10px 30px rgba(0,0,0,0.1); }

/* Footer */
.footer {
  padding: 80px 0 0;
  background: var(--bg-alt);
}

.footer-grid {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 80px;
  margin-bottom: 80px;
}

@media (max-width: 968px) {
  .footer-grid { grid-template-columns: 1fr; gap: 40px; }
}

.footer-brand .logo { margin-bottom: 20px; }
.footer-brand p { max-width: 320px; }

.footer-links-group {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 40px;
}

.link-col h5 {
  font-size: 18px;
  font-weight: 800;
  margin-bottom: 24px;
}

.link-col a {
  display: block;
  color: var(--text-muted);
  text-decoration: none;
  margin-bottom: 12px;
  font-weight: 500;
  transition: color 0.2s;
}

.link-col a:hover { color: var(--primary); }

.footer-bottom {
  padding: 30px 0;
  border-top: 1px solid var(--border);
  text-align: center;
  font-size: 14px;
  color: var(--text-muted);
}

/* Utils */
.lang-switcher { position: relative; }
.lang-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 8px;
  box-shadow: var(--card-shadow);
  min-width: 140px;
  z-index: 1100;
  margin-top: 8px;
}

.lang-dropdown button {
  width: 100%;
  padding: 10px 16px;
  border: none;
  background: none;
  text-align: left;
  cursor: pointer;
  border-radius: 8px;
  font-weight: 600;
  color: var(--text);
}

.lang-dropdown button:hover { background: var(--bg-alt); color: var(--primary); }

.rtl-flip { transition: transform 0.3s; }
[dir="rtl"] .rtl-flip { transform: rotate(180deg); }
/* RTL handling for dropdown */
[dir="rtl"] .lang-dropdown {
  right: auto;
  left: 0;
}
[dir="rtl"] .lang-dropdown button {
  text-align: right;
}
</style>

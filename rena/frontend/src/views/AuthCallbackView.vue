<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { CheckCircle, XCircle, ArrowRight, Loader2 } from 'lucide-vue-next'
import gsap from 'gsap'
import iconImg from '../assets/icon.png'

const router = useRouter()
const auth = useAuthStore()
const status = ref<'loading' | 'success' | 'error'>('loading')
const message = ref('')
const callbackCard = ref<HTMLElement | null>(null)

onMounted(async () => {
  document.title = 'Rena - Processing Authentication'
  
  // Animate card in
  gsap.from(callbackCard.value, {
    y: 30,
    opacity: 0,
    duration: 0.8,
    ease: 'power3.out'
  })

  // Parse URL Hash
  const hash = window.location.hash
  const params = new URLSearchParams(hash.substring(1))
  
  const accessToken = params.get('access_token')
  const error = params.get('error_description') || params.get('error')

  // Immediately clear the hash from the URL for security
  if (hash) {
    window.history.replaceState(null, '', window.location.pathname + window.location.search)
  }

  if (accessToken) {
    status.value = 'success'
    message.value = 'Welcome back! Your authentication was successful.'
    
    // Success animation
    setTimeout(() => {
      gsap.to(".success-icon", {
        scale: 1.2,
        duration: 0.5,
        yoyo: true,
        repeat: 1,
        ease: "back.out(1.7)"
      })
    }, 100)

    // Automatically redirect after a few seconds
    setTimeout(() => {
      router.push('/dashboard')
    }, 3000)
  } else if (error) {
    status.value = 'error'
    message.value = error.replace(/\+/g, ' ')
  } else {
    // Check if we are already logged in (case where user manually goes to this page)
    if (auth.user) {
      status.value = 'success'
      message.value = 'You are already logged in. Redirecting to dashboard...'
      setTimeout(() => router.push('/dashboard'), 2000)
    } else {
      status.value = 'error'
      message.value = 'No authentication data found. Please try logging in again.'
    }
  }
})

const goHome = () => router.push('/')
const goLogin = () => router.push('/login')
const goDashboard = () => router.push('/dashboard')
</script>

<template>
  <div class="callback-container">
    <div ref="callbackCard" class="card">
      <div class="logo-section">
        <img :src="iconImg" alt="Rena Logo" class="logo-img" />
      </div>

      <div v-if="status === 'loading'" class="status-content">
        <Loader2 class="animate-spin text-primary" :size="48" />
        <h2>Processing...</h2>
        <p>Please wait while we verify your credentials.</p>
      </div>

      <div v-else-if="status === 'success'" class="status-content">
        <div class="icon-wrapper success">
          <CheckCircle class="success-icon" :size="64" />
        </div>
        <h2>Congratulations!</h2>
        <p>{{ message }}</p>
        <button @click="goDashboard" class="primary-btn">
          Go to Dashboard <ArrowRight :size="18" />
        </button>
        <p class="auto-redirect">Redirecting automatically in a moment...</p>
      </div>

      <div v-else-if="status === 'error'" class="status-content">
        <div class="icon-wrapper error">
          <XCircle :size="64" />
        </div>
        <h2>Authentication Failed</h2>
        <p>{{ message }}</p>
        <div class="btn-group">
          <button @click="goLogin" class="primary-btn">Try Again</button>
          <button @click="goHome" class="secondary-btn">Back to Home</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.callback-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--bg-alt);
  padding: 24px;
}

.card {
  background: var(--surface);
  width: 100%;
  max-width: 480px;
  padding: 48px;
  border-radius: 32px;
  box-shadow: var(--card-shadow);
  text-align: center;
  border: 1px solid var(--border);
}

.logo-section {
  margin-bottom: 32px;
}

.logo-img {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: var(--grad-brand);
}

.status-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.icon-wrapper {
  margin-bottom: 8px;
}

.icon-wrapper.success { color: #10b981; }
.icon-wrapper.error { color: #ef4444; }

h2 {
  font-size: 28px;
  font-weight: 800;
  color: var(--text);
}

p {
  color: var(--text-muted);
  line-height: 1.6;
  font-size: 16px;
}

.primary-btn {
  background: var(--primary);
  color: white;
  border: none;
  padding: 14px 32px;
  border-radius: 50px;
  font-weight: 700;
  font-size: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  margin-top: 12px;
}

.primary-btn:hover {
  background: var(--primary-dark);
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 8px 25px rgba(46, 96, 255, 0.2);
}

.secondary-btn {
  background: none;
  border: 1px solid var(--border);
  color: var(--text-muted);
  padding: 14px 32px;
  border-radius: 50px;
  font-weight: 600;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 12px;
}

.secondary-btn:hover {
  background: var(--bg-alt);
  color: var(--text);
}

.btn-group {
  display: flex;
  gap: 12px;
  width: 100%;
  justify-content: center;
}

.auto-redirect {
  font-size: 13px;
  margin-top: 16px;
  opacity: 0.7;
}

.text-primary { color: var(--primary); }
</style>

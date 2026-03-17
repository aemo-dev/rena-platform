<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { Github, Layout, Smartphone, Globe, Code, Layers, Zap, Shield, Eye, EyeOff, Mail, CheckCircle } from 'lucide-vue-next'
import gsap from 'gsap'
import iconImg from '../assets/icon.png'

const { t } = useI18n()
const auth = useAuthStore()
const route = useRoute()
const loginCard = ref<HTMLElement | null>(null)

const iconsList = [Layout, Smartphone, Globe, Code, Layers, Zap, Shield, Github]
const colors = ['#2e60ff', '#ff4081', '#4caf50', '#ffc107', '#00bcd4', '#3f51b5']

interface Shape {
  id: number
  icon: any
  size: number
  top: string
  left: string
  color: string
}

const shapes = ref<Shape[]>([])
const rememberMe = ref(false)
const isSignUp = ref(false)
const isForgotPassword = ref(false)
const showSuccessModal = ref(false)
const firstName = ref('')
const lastName = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const agreeToTerms = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const passwordStrength = computed(() => {
  if (!password.value) return 0
  let strength = 0
  if (password.value.length >= 8) strength += 25
  if (/[A-Z]/.test(password.value)) strength += 25
  if (/[0-9]/.test(password.value)) strength += 25
  if (/[^A-Za-z0-9]/.test(password.value)) strength += 25
  return strength
})

const strengthColor = computed(() => {
  const s = passwordStrength.value
  if (s <= 25) return '#ef4444' // Red
  if (s <= 50) return '#f59e0b' // Amber
  if (s <= 75) return '#3b82f6' // Blue
  return '#10b981' // Green
})

onMounted(() => {
  document.title = 'Rena Auth'

  // Check for error in query params and clear URL
  if (route.query.error) {
    const errorCode = route.query.error as string
    const errorDesc = route.query.error_description as string
    
    let message = errorDesc ? errorDesc.replace(/\+/g, ' ') : 'Authentication failed. Please try again.'
    
    if (errorCode === 'bad_oauth_state') {
      message = 'Your session expired. Please try logging in again.'
    } else if (errorCode === 'invalid_request') {
      message = 'Invalid request. Please try logging in again.'
    }
    
    errorMessage.value = message
    // Clear the error from URL
    window.history.replaceState({}, '', window.location.pathname)
  }
  
  // Generate random shapes
  for (let i = 0; i < 20; i++) {
    shapes.value.push({
      id: i,
      icon: iconsList[Math.floor(Math.random() * iconsList.length)],
      size: Math.random() * 40 + 20,
      top: `${Math.random() * 100}%`,
      left: `${Math.random() * 100}%`,
      color: colors[Math.floor(Math.random() * colors.length)]
    })
  }

  // Animate Card
  gsap.from(loginCard.value, { 
    y: 30, 
    opacity: 0, 
    duration: 0.8, 
    ease: 'power3.out' 
  })

  // Animate shapes
  setTimeout(() => {
    const elShapes = document.querySelectorAll('.random-shape')
    elShapes.forEach((el) => {
      gsap.to(el, {
        x: 'random(-150, 150)',
        y: 'random(-150, 150)',
        rotation: 'random(-360, 360)',
        duration: 'random(3, 6)',
        repeat: -1,
        yoyo: true,
        ease: 'sine.inOut'
      })
    })
  }, 100)
})

const login = (provider: 'google' | 'github') => {
  auth.loginWithProvider(provider)
}

const handleEmailAuth = async () => {
  if (isSignUp.value) {
    if (!firstName.value || !lastName.value || !email.value || !password.value || !confirmPassword.value) {
      errorMessage.value = 'Please fill in all fields'
      return
    }
    if (password.value !== confirmPassword.value) {
      errorMessage.value = 'Passwords do not match'
      return
    }
    if (!agreeToTerms.value) {
      errorMessage.value = 'You must agree to the Terms of Service'
      return
    }
  } else {
    if (!email.value || !password.value) {
      errorMessage.value = 'Please enter both email and password'
      return
    }
  }
  
  loading.value = true
  errorMessage.value = ''
  
  try {
    if (isSignUp.value) {
      await auth.signUpWithEmail(email.value, password.value, {
        first_name: firstName.value,
        last_name: lastName.value,
        full_name: `${firstName.value} ${lastName.value}`
      })
      showSuccessModal.value = true
    } else {
      await auth.signInWithEmail(email.value, password.value)
    }
  } catch (err: any) {
    errorMessage.value = err.message || 'Authentication failed'
  } finally {
    loading.value = false
  }
}

const handleForgotPassword = async () => {
  if (!email.value) {
    errorMessage.value = 'Please enter your email address'
    return
  }
  
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    await auth.resetPassword(email.value)
    successMessage.value = 'If an account exists for this email, you will receive a reset link shortly.'
  } catch (err: any) {
    errorMessage.value = err.message || 'Failed to send reset link'
  } finally {
    loading.value = false
  }
}

const toggleForgotPassword = () => {
  isForgotPassword.value = !isForgotPassword.value
  isSignUp.value = false
  errorMessage.value = ''
  successMessage.value = ''
}
</script>

<template>
  <div class="login-container">
    <div class="login-bg-decor">
      <div 
        v-for="shape in shapes" 
        :key="shape.id"
        class="random-shape"
        :style="{
          top: shape.top,
          left: shape.left,
          width: shape.size + 'px',
          height: shape.size + 'px',
          color: shape.color,
          opacity: 0.15
        }"
      >
        <component :is="shape.icon" :size="shape.size" />
      </div>
    </div>

    <div ref="loginCard" class="card">
      <div class="card-header">
        <img :src="iconImg" alt="Rena Logo" class="logo-img" />
        <h1 v-if="isForgotPassword">Reset Password</h1>
        <h1 v-else>{{ isSignUp ? 'Create Account' : t('login.title') }}</h1>
        
        <p v-if="isForgotPassword">Enter your email and we'll send you a link to reset your password.</p>
        <p v-else>{{ isSignUp ? 'Join Rena Builder and start creating today' : t('login.sub') }}</p>
      </div>

      <!-- Forgot Password Form -->
      <form v-if="isForgotPassword" @submit.prevent="handleForgotPassword" class="email-form">
        <div class="form-group">
          <input 
            type="email" 
            v-model="email" 
            placeholder="Email address" 
            required
            class="form-input"
          />
        </div>
        
        <div v-if="errorMessage" class="error-text">{{ errorMessage }}</div>
        <div v-if="successMessage" class="success-text">{{ successMessage }}</div>

        <button type="submit" class="submit-btn" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          <span v-else>Send Reset Link</span>
        </button>
        
        <button type="button" @click="toggleForgotPassword" class="back-to-login">
          Back to Login
        </button>
      </form>

      <!-- Login/SignUp Form -->
      <form v-else @submit.prevent="handleEmailAuth" class="email-form">
        <div v-if="isSignUp" class="form-row">
          <div class="form-group">
            <input 
              type="text" 
              v-model="firstName" 
              placeholder="First Name" 
              required
              class="form-input"
            />
          </div>
          <div class="form-group">
            <input 
              type="text" 
              v-model="lastName" 
              placeholder="Last Name" 
              required
              class="form-input"
            />
          </div>
        </div>

        <div class="form-group">
          <input 
            type="email" 
            v-model="email" 
            placeholder="Email address" 
            required
            class="form-input"
          />
        </div>
        <div class="form-group">
          <input 
            :type="showPassword ? 'text' : 'password'" 
            v-model="password" 
            placeholder="Password" 
            required
            class="form-input with-icon"
          />
          <button type="button" @click="showPassword = !showPassword" class="password-toggle">
            <component :is="showPassword ? EyeOff : Eye" :size="18" />
          </button>
          
          <div v-if="isSignUp" class="strength-circle-wrapper">
            <svg class="strength-svg" viewBox="0 0 24 24">
              <circle class="strength-bg" cx="12" cy="12" r="10"></circle>
              <circle 
                class="strength-fill" 
                cx="12" 
                cy="12" 
                r="10" 
                :style="{
                  strokeDasharray: '62.8',
                  strokeDashoffset: (62.8 - (62.8 * passwordStrength) / 100),
                  stroke: strengthColor
                }"
              ></circle>
            </svg>
          </div>

          <div v-if="!isSignUp" class="forgot-link-wrapper">
            <button type="button" @click="toggleForgotPassword" class="forgot-btn">
              Forgot password?
            </button>
          </div>
        </div>
        <div v-if="isSignUp" class="form-group">
          <input 
            :type="showConfirmPassword ? 'text' : 'password'" 
            v-model="confirmPassword" 
            placeholder="Confirm Password" 
            required
            class="form-input with-icon"
          />
          <button type="button" @click="showConfirmPassword = !showConfirmPassword" class="password-toggle">
            <component :is="showConfirmPassword ? EyeOff : Eye" :size="18" />
          </button>
        </div>

        <div v-if="isSignUp" class="terms-checkbox">
          <label class="checkbox-label">
            <input type="checkbox" v-model="agreeToTerms" required />
            <span class="checkmark"></span>
            I agree to the <span class="terms-link">Terms of Service</span>
          </label>
        </div>

        <div v-else class="remember-me-section">
          <label class="checkbox-label">
            <input type="checkbox" v-model="rememberMe" />
            <span class="checkmark"></span>
            Remember me
          </label>
        </div>
        
        <div v-if="errorMessage" class="error-text">{{ errorMessage }}</div>

        <button type="submit" class="submit-btn" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          <span v-else>{{ isSignUp ? 'Sign Up' : 'Sign In' }}</span>
        </button>
      </form>

      <div v-if="!isForgotPassword" class="divider">
        <span>OR</span>
      </div>

      <div v-if="!isForgotPassword" class="auth-buttons">
        <button @click="login('google')" class="auth-btn google">
          <img src="https://www.gstatic.com/firebasejs/ui/2.0.0/images/auth/google.svg" alt="Google" />
          {{ t('login.google') }}
        </button>
        
        <button @click="login('github')" class="auth-btn github">
          <Github :size="20" />
          {{ t('login.github') }}
        </button>
      </div>

      <div v-if="!isForgotPassword" class="toggle-auth">
        {{ isSignUp ? 'Already have an account?' : "Don't have an account?" }}
        <button @click="isSignUp = !isSignUp" class="toggle-btn">
          {{ isSignUp ? 'Sign In' : 'Sign Up' }}
        </button>
      </div>

      <div class="card-footer">
        <p>{{ t('login.agree') }} <span class="terms-link">{{ t('login.terms') }}</span>.</p>
        <router-link to="/" class="back-link">← {{ t('login.back') }}</router-link>
      </div>
    </div>

    <!-- Success Modal -->
    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showSuccessModal" class="success-modal-overlay" @click.self="showSuccessModal = false">
          <div class="success-modal">
            <div class="success-icon-wrapper">
              <div class="success-ripple"></div>
              <CheckCircle :size="64" class="success-check-icon" />
            </div>
            
            <h2>Verify Your Email</h2>
            <p>We've sent a magic link to <strong>{{ email }}</strong>. Please check your inbox and click the link to activate your account.</p>
            
            <div class="success-modal-actions">
              <button @click="showSuccessModal = false" class="modal-primary-btn">
                <Mail :size="18" />
                Got it, thanks!
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: var(--bg-alt);
  position: relative;
  overflow: hidden;
}

.login-bg-decor {
  position: absolute;
  width: 100%;
  height: 100%;
  z-index: 1;
  pointer-events: none;
}

:deep(.random-shape) {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card {
  background: var(--surface);
  width: 100%;
  max-width: 440px; /* Slightly wider for form */
  padding: 40px;
  border-radius: 28px;
  box-shadow: var(--card-shadow);
  text-align: center;
  z-index: 2;
  border: 1px solid var(--border);
  backdrop-filter: blur(10px);
}

.card-header {
  margin-bottom: 32px;
}

.logo-img {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  object-fit: cover;
  margin-bottom: 24px;
}

h1 {
  font-size: 28px;
  font-weight: 700;
  color: var(--text);
  margin-bottom: 12px;
}

p {
  color: var(--text-muted);
  font-size: 15px;
  line-height: 1.5;
}

.email-form {
  display: grid;
  gap: 16px;
  margin-bottom: 24px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.terms-checkbox,
.remember-me-section {
  text-align: left;
  margin-top: -8px;
}

.form-group {
  position: relative;
}

.form-input {
  width: 100%;
  padding: 14px 20px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text);
  font-size: 15px;
  font-weight: 500;
  transition: all 0.2s;
  box-sizing: border-box;
}

.form-input.with-icon {
  padding-right: 48px;
}

.password-toggle {
  position: absolute;
  right: 14px;
  top: 14px;
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s;
}

.password-toggle:hover {
  color: var(--primary);
}

.strength-circle-wrapper {
  position: absolute;
  right: -36px;
  top: 12px;
  width: 24px;
  height: 24px;
}

.strength-svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.strength-bg {
  fill: none;
  stroke: var(--border);
  stroke-width: 3;
}

.strength-fill {
  fill: none;
  stroke-width: 3;
  stroke-linecap: round;
  transition: stroke-dashoffset 0.4s ease, stroke 0.4s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 4px rgba(46, 96, 255, 0.1);
}

.submit-btn {
  width: 100%;
  padding: 14px;
  border-radius: 14px;
  border: none;
  background: var(--primary);
  color: white;
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.submit-btn:hover {
  background: var(--primary-dark);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(46, 96, 255, 0.2);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.error-text {
  color: #ef4444;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
}

.success-text {
  color: #10b981;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
}

.forgot-link-wrapper {
  text-align: right;
  margin-top: 8px;
}

.forgot-btn {
  background: none;
  border: none;
  color: var(--primary);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  padding: 0;
}

.forgot-btn:hover {
  text-decoration: underline;
}

.back-to-login {
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  margin-top: 16px;
}

.back-to-login:hover {
  color: var(--primary);
}

.divider {
  display: flex;
  align-items: center;
  text-align: center;
  margin: 24px 0;
  color: var(--text-muted);
  font-size: 13px;
  font-weight: 600;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid var(--border);
}

.divider span {
  padding: 0 16px;
}

.auth-buttons {
  display: grid;
  gap: 12px;
  margin-bottom: 24px;
}

.auth-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 12px 24px;
  border-radius: 14px;
  cursor: pointer;
  font-size: 15px;
  font-weight: 600;
  transition: all 0.2s;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text);
}

.auth-btn.github {
  background: #1f2937;
  color: white;
  border-color: #1f2937;
}

.auth-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  border-color: var(--primary);
}

.auth-btn img {
  width: 20px;
}

.toggle-auth {
  font-size: 14px;
  color: var(--text-muted);
  margin-bottom: 32px;
}

.toggle-btn {
  background: none;
  border: none;
  color: var(--primary);
  font-weight: 700;
  cursor: pointer;
  padding: 0 4px;
}

.toggle-btn:hover {
  text-decoration: underline;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255,255,255,0.3);
  border-radius: 50%;
  border-top-color: #fff;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.card-footer {
  border-top: 1px solid var(--border);
  padding-top: 24px;
}

.terms-link {
  color: var(--primary);
  text-decoration: none;
  font-weight: 600;
}

.terms-link:hover {
  text-decoration: underline;
}

.back-link {
  display: block;
  margin-top: 24px;
  color: var(--text-muted);
  text-decoration: none;
  font-weight: 600;
}

.back-link:hover {
  color: var(--primary);
}

[dir="rtl"] .back-link {
  direction: rtl;
}

/* Success Modal Styles */
.success-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(8px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
}

.success-modal {
  background: var(--surface);
  width: 90%;
  max-width: 440px;
  padding: 48px 40px;
  border-radius: 32px;
  text-align: center;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.2);
  border: 1px solid var(--border);
  position: relative;
  overflow: hidden;
}

.success-icon-wrapper {
  position: relative;
  width: 80px;
  height: 80px;
  margin: 0 auto 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.success-check-icon {
  color: #10b981;
  position: relative;
  z-index: 2;
}

.success-ripple {
  position: absolute;
  width: 100%;
  height: 100%;
  background: rgba(16, 185, 129, 0.1);
  border-radius: 50%;
  animation: ripple 2s infinite;
}

@keyframes ripple {
  0% { transform: scale(0.8); opacity: 1; }
  100% { transform: scale(2); opacity: 0; }
}

.success-modal h2 {
  font-size: 26px;
  font-weight: 800;
  color: var(--text);
  margin-bottom: 16px;
}

.success-modal p {
  color: var(--text-muted);
  line-height: 1.6;
  font-size: 16px;
  margin-bottom: 32px;
}

.modal-primary-btn {
  width: 100%;
  background: var(--primary);
  color: white;
  border: none;
  padding: 16px;
  border-radius: 18px;
  font-weight: 700;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.modal-primary-btn:hover {
  background: var(--primary-dark);
  transform: translateY(-3px);
  box-shadow: 0 10px 20px rgba(46, 96, 255, 0.2);
}

/* Transitions */
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.4s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

.modal-fade-enter-active .success-modal {
  animation: modal-in 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}

@keyframes modal-in {
  0% { transform: scale(0.8) translateY(20px); opacity: 0; }
  100% { transform: scale(1) translateY(0); opacity: 1; }
}
</style>

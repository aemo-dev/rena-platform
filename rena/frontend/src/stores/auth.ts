import { defineStore } from 'pinia'
import { setAuthTokens, clearAuthTokens } from '../services/api'

type UserMeta = {
  full_name?: string
  name?: string
}

type AuthUser = {
  id: string
  email: string
  user_metadata?: UserMeta
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as AuthUser | null,
    loading: false
  }),
  actions: {
    async autoLogin() {
      if (this.loading) return

      const existingUserId = localStorage.getItem('user_id')
      const existingEmail = localStorage.getItem('user_email')
      if (existingUserId && existingEmail) {
        this.loadFromStorage()
        return
      }

      this.loading = true
      try {
        const backendUrl = (import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080').replace(/\/$/, '')
        const res = await fetch(`${backendUrl}/api/auth/device-login`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
        })
        const data = await res.json()
        if (!res.ok || !data.token || !data.user_id || !data.email) {
          throw new Error(data.error || 'Device login failed')
        }

        const metadata = data.device_name ? { name: String(data.device_name) } : undefined
        setAuthTokens(String(data.user_id), data.token)
        localStorage.setItem('user_email', String(data.email))
        if (metadata) {
          localStorage.setItem('user_metadata', JSON.stringify(metadata))
        } else {
          localStorage.removeItem('user_metadata')
        }

        this.user = {
          id: String(data.user_id),
          email: String(data.email),
          user_metadata: metadata,
        }
      } finally {
        this.loading = false
      }
    },
    async fetchUser() {
      const userId = localStorage.getItem('user_id')
      const email = localStorage.getItem('user_email')
      if (!userId || !email) {
        this.user = null
        await this.autoLogin()
        return
      }
      const metadata = localStorage.getItem('user_metadata')
      this.user = {
        id: userId,
        email,
        user_metadata: metadata ? JSON.parse(metadata) : undefined,
      }
    },
    async signInWithEmail(email: string, password: string) {
      const res = await fetch(`${(import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080').replace(/\/$/, '')}/api/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      })
      const data = await res.json()
      if (!res.ok || !data.token) {
        throw new Error(data.error || 'Login failed')
      }
      setAuthTokens(String(data.user_id), data.token)
      localStorage.setItem('user_email', email)
      this.user = { id: String(data.user_id), email }
    },
    async signUpWithEmail(email: string, password: string, metadata?: any) {
      const res = await fetch(`${(import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080').replace(/\/$/, '')}/api/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      })
      const data = await res.json()
      if (!res.ok || !data.token) {
        throw new Error(data.error || 'Registration failed')
      }
      setAuthTokens(String(data.user_id), data.token)
      localStorage.setItem('user_email', email)
      if (metadata) {
        localStorage.setItem('user_metadata', JSON.stringify(metadata))
      }
      this.user = { id: String(data.user_id), email, user_metadata: metadata }
    },
    async loginWithProvider(_: string) {
      throw new Error('OAuth login is not supported in local auth')
    },
    async resetPassword(email: string) {
      console.warn('resetPassword not implemented for local auth')
      if (!email) throw new Error('Email required')
    },
    logout() {
      clearAuthTokens()
      localStorage.removeItem('user_email')
      localStorage.removeItem('user_metadata')
      this.user = null
    },
    loadFromStorage() {
      const userId = localStorage.getItem('user_id')
      const email = localStorage.getItem('user_email')
      const metadataStr = localStorage.getItem('user_metadata')
      if (userId && email) {
        this.user = {
          id: userId,
          email,
          user_metadata: metadataStr ? JSON.parse(metadataStr) : undefined,
        }
      } else {
        this.user = null
      }
    }
  }
})

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useAuthStore } from './stores/auth'
import { useRouter } from 'vue-router'
import { setLanguage } from './i18n'

const auth = useAuthStore()
const router = useRouter()
let lastEnsuredUserId: string | null = null
let isEnsuringKeystore = false

const ensureKeystore = async (user: any) => {
  if (!user || !user.email || isEnsuringKeystore || lastEnsuredUserId === user.id) return
  isEnsuringKeystore = true
  const backendUrl = (import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080').replace(/\/$/, '')
  try {
    console.log('[Keystore] Ensuring keystore exists for user...')
    const response = await fetch(`${backendUrl}/api/keystore/generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_id: user.id, email: user.email })
    })
    if (!response.ok) {
      const errorData = await response.json()
      console.error('[Keystore] Generation failed:', errorData.error)
    } else {
      lastEnsuredUserId = user.id
    }
  } catch (err) {
    console.error('[Keystore] Failed to reach backend:', err)
  } finally {
    isEnsuringKeystore = false
  }
}

onMounted(() => {
  setLanguage('en')
  auth.loadFromStorage()
  if (auth.user) {
    ensureKeystore(auth.user)
  }
})

watch(
  () => auth.user,
  (newUser) => {
    if (newUser) {
      ensureKeystore(newUser)
      if (router.currentRoute.value.path === '/login') {
        router.push('/dashboard')
      }
    } else {
      if (router.currentRoute.value.meta.requiresAuth) {
        router.push('/login')
      }
    }
  }
)
</script>

<template>
  <router-view />
</template>

<style>
/* Global styles are in style.css */
</style>

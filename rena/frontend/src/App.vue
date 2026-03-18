<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { supabase } from './supabase'
import { useAuthStore } from './stores/auth'
import { useRouter } from 'vue-router'
import { setLanguage } from './i18n'

const auth = useAuthStore()
const router = useRouter()
let authListener: any = null
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
      body: JSON.stringify({
        user_id: user.id,
        email: user.email
      })
    })

    if (!response.ok) {
      const errorData = await response.json()
      console.error('[Keystore] Generation failed:', errorData.error)
    } else {
      console.log('[Keystore] Keystore verified/generated successfully')
      lastEnsuredUserId = user.id
    }
  } catch (err) {
    console.error('[Keystore] Failed to reach backend:', err)
  } finally {
    isEnsuringKeystore = false
  }
}

onMounted(() => {
  // Initialize language/direction
  setLanguage('en')
  
  // Listen for auth state changes
  const { data: { subscription } } = supabase.auth.onAuthStateChange((event, session) => {
    if (event === 'SIGNED_IN' || event === 'INITIAL_SESSION') {
      auth.user = session?.user || null
      
      if (session?.user) {
        ensureKeystore(session.user)
      }

      // Only auto-redirect from login page
      const currentPath = router.currentRoute.value.path
      if (session && currentPath === '/login') {
        router.push('/dashboard')
      }
    } else if (event === 'SIGNED_OUT') {
      auth.user = null
      lastEnsuredUserId = null
      // Redirect to login if on a protected page
      if (router.currentRoute.value.meta.requiresAuth) {
        router.push('/login')
      }
    }
  })
  
  authListener = subscription
})

onUnmounted(() => {
  if (authListener) {
    authListener.unsubscribe()
  }
})
</script>

<template>
  <router-view />
</template>

<style>
/* Global styles are in style.css */
</style>

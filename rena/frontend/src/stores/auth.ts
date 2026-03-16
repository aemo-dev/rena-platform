import { defineStore } from 'pinia'
import { supabase } from '../supabase'
import type { User, Provider } from '@supabase/supabase-js'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    loading: true
  }),
  actions: {
    async fetchUser() {
      const { data: { user } } = await supabase.auth.getUser()
      this.user = user
      this.loading = false
    },
    async loginWithProvider(provider: Provider) {
      const { error } = await supabase.auth.signInWithOAuth({
        provider: provider,
        options: {
            redirectTo: window.location.origin + '/auth/callback'
        }
      })
      if (error) throw error
    },
    async signInWithEmail(email: string, pass: string) {
      const { data, error } = await supabase.auth.signInWithPassword({
        email,
        password: pass
      })
      if (error) throw error
      this.user = data.user
    },
    async signUpWithEmail(email: string, pass: string, metadata?: any) {
      const { data, error } = await supabase.auth.signUp({
        email,
        password: pass,
        options: {
          data: metadata,
          emailRedirectTo: window.location.origin + '/auth/callback'
        }
      })
      if (error) throw error
      this.user = data.user
    },
    async resetPassword(email: string) {
      const { error } = await supabase.auth.resetPasswordForEmail(email, {
        redirectTo: window.location.origin + '/auth/callback?reset=true'
      })
      if (error) throw error
    },
    async logout() {
      await supabase.auth.signOut()
      this.user = null
    }
  }
})

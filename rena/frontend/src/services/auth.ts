/**
 * Authentication API Service
 * Handles all authentication-related operations with Supabase
 * @module auth
 */

import { supabase } from '../supabase'
import type { User, Provider } from '@supabase/supabase-js'
import apiFetch, { type ApiResponse } from './api'

/**
 * Login request payload
 */
export interface LoginRequest {
  email: string
  password: string
}

/**
 * Registration request payload
 */
export interface RegisterRequest {
  email: string
  password: string
  metadata?: {
    full_name?: string
    first_name?: string
    last_name?: string
    [key: string]: unknown
  }
}

/**
 * Password reset request payload
 */
export interface PasswordResetRequest {
  email: string
}

/**
 * OAuth provider sign-in options
 */
export interface OAuthSignInOptions {
  provider: Provider
  redirectTo?: string
}

/**
 * Email sign-in options
 */
export interface EmailSignInOptions {
  email: string
  password: string
}

/**
 * Sign-up options
 */
export interface SignUpOptions {
  email: string
  password: string
  metadata?: {
    full_name?: string
    first_name?: string
    last_name?: string
    [key: string]: unknown
  }
  emailRedirectTo?: string
}

/**
 * Authentication result with user and session information
 */
export interface AuthResult {
  user: User | null
  sessionId?: string | null
  accessToken?: string | null
  refreshToken?: string | null
}

/**
 * Get current authenticated user from Supabase
 * @returns Promise resolving to user data or null
 * @throws Error if network request fails
 */
export async function getCurrentUser(): Promise<User | null> {
  try {
    const { data: { user }, error } = await supabase.auth.getUser()
    
    if (error) {
      console.warn('Failed to get current user:', error.message)
      return null
    }
    
    return user
  } catch (error) {
    console.error('Unexpected error getting current user:', error)
    return null
  }
}

/**
 * Fetch current user via backend API (alternative method)
 * @returns Promise resolving to API response with user data
 */
export async function fetchCurrentUserFromApi(): Promise<ApiResponse<{ user: User }>> {
  return await apiFetch<{ user: User }>('/api/auth/user')
}

/**
 * Sign in with email and password
 * @param email - User's email address
 * @param password - User's password
 * @returns Promise resolving to authentication result
 * @throws Error if authentication fails
 */
export async function signInWithEmail(email: string, password: string): Promise<AuthResult> {
  try {
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password,
    })

    if (error) {
      throw new AuthError(error.message, 'SIGN_IN_FAILED')
    }

    return {
      user: data.user,
      sessionId: data.session?.provider_token,
      accessToken: data.session?.access_token,
      refreshToken: data.session?.refresh_token,
    }
  } catch (error) {
    if (error instanceof AuthError) {
      throw error
    }
    throw new AuthError(
      error instanceof Error ? error.message : 'Failed to sign in',
      'SIGN_IN_FAILED'
    )
  }
}

/**
 * Sign up with email and password
 * @param email - User's email address
 * @param password - User's password
 * @param metadata - Optional user metadata (name, etc.)
 * @param emailRedirectTo - URL to redirect after email confirmation
 * @returns Promise resolving to authentication result
 * @throws Error if registration fails
 */
export async function signUpWithEmail(
  email: string,
  password: string,
  metadata?: SignUpOptions['metadata'],
  emailRedirectTo?: string
): Promise<AuthResult> {
  try {
    const { data, error } = await supabase.auth.signUp({
      email,
      password,
      options: {
        data: metadata,
        emailRedirectTo: emailRedirectTo ?? `${window.location.origin}/auth/callback`,
      },
    })

    if (error) {
      throw new AuthError(error.message, 'SIGN_UP_FAILED')
    }

    return {
      user: data.user,
      sessionId: data.session?.provider_token,
      accessToken: data.session?.access_token,
      refreshToken: data.session?.refresh_token,
    }
  } catch (error) {
    if (error instanceof AuthError) {
      throw error
    }
    throw new AuthError(
      error instanceof Error ? error.message : 'Failed to register',
      'SIGN_UP_FAILED'
    )
  }
}

/**
 * Sign in with OAuth provider (Google, GitHub, etc.)
 * @param provider - OAuth provider name
 * @param redirectTo - Optional redirect URL after authentication
 * @returns Promise resolving to OAuth redirect result
 * @throws Error if OAuth sign-in fails
 */
export async function signInWithOAuth(provider: Provider, redirectTo?: string): Promise<{ url?: string }> {
  try {
    const siteUrl = import.meta.env.VITE_SITE_URL || window.location.origin
    const { data, error } = await supabase.auth.signInWithOAuth({
      provider,
      options: {
        redirectTo: redirectTo ?? `${siteUrl}/auth/callback`,
      },
    })

    if (error) {
      throw new AuthError(error.message, 'OAUTH_SIGN_IN_FAILED')
    }

    return {
      url: data.url,
    }
  } catch (error) {
    if (error instanceof AuthError) {
      throw error
    }
    throw new AuthError(
      error instanceof Error ? error.message : 'Failed to sign in with OAuth',
      'OAUTH_SIGN_IN_FAILED'
    )
  }
}

/**
 * Send password reset email
 * @param email - User's email address
 * @param redirectTo - URL to redirect for password reset
 * @returns Promise resolving when email is sent
 * @throws Error if password reset fails
 */
export async function sendPasswordResetEmail(
  email: string,
  redirectTo?: string
): Promise<void> {
  try {
    const { error } = await supabase.auth.resetPasswordForEmail(email, {
      redirectTo: redirectTo ?? `${window.location.origin}/auth/callback?reset=true`,
    })

    if (error) {
      throw new AuthError(error.message, 'PASSWORD_RESET_FAILED')
    }
  } catch (error) {
    if (error instanceof AuthError) {
      throw error
    }
    throw new AuthError(
      error instanceof Error ? error.message : 'Failed to send password reset email',
      'PASSWORD_RESET_FAILED'
    )
  }
}

/**
 * Sign out current user
 * @returns Promise resolving when user is signed out
 * @throws Error if logout fails
 */
export async function signOut(): Promise<void> {
  try {
    const { error } = await supabase.auth.signOut()

    if (error) {
      throw new AuthError(error.message, 'LOGOUT_FAILED')
    }
  } catch (error) {
    if (error instanceof AuthError) {
      throw error
    }
    throw new AuthError(
      error instanceof Error ? error.message : 'Failed to logout',
      'LOGOUT_FAILED'
    )
  }
}

/**
 * Refresh current session
 * @returns Promise resolving to refreshed authentication result
 * @throws Error if session refresh fails
 */
export async function refreshSession(): Promise<AuthResult> {
  try {
    const { data, error } = await supabase.auth.refreshSession()

    if (error) {
      throw new AuthError(error.message, 'SESSION_REFRESH_FAILED')
    }

    return {
      user: data.user,
      sessionId: data.session?.provider_token,
      accessToken: data.session?.access_token,
      refreshToken: data.session?.refresh_token,
    }
  } catch (error) {
    if (error instanceof AuthError) {
      throw error
    }
    throw new AuthError(
      error instanceof Error ? error.message : 'Failed to refresh session',
      'SESSION_REFRESH_FAILED'
    )
  }
}

/**
 * Custom authentication error with error code
 */
export class AuthError extends Error {
  public readonly code: string
  public readonly originalError?: Error

  constructor(message: string, code: string, originalError?: Error) {
    super(message)
    this.name = 'AuthError'
    this.code = code
    this.originalError = originalError

    // Maintain proper prototype chain
    Object.setPrototypeOf(this, AuthError.prototype)
  }

  /**
   * Check if error is a specific error type
   */
  isErrorCode(expectedCode: string): boolean {
    return this.code === expectedCode
  }
}

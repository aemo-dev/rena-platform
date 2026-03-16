/**
 * API Client for communicating with the Go backend
 */

const API_BASE_URL = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'

export interface ApiResponse<T> {
  data?: T
  error?: string
  message?: string
}

export interface User {
  id: string
  email: string
  user_metadata?: {
    full_name?: string
    name?: string
    [key: string]: any
  }
  created_at?: string
}

export interface Project {
  id: string
  user_id: string
  name: string
  package_name: string
  platform: string
  color: string
  icon_url?: string
  version_code: number
  version_name: string
  workspace_xml?: string
  generated_code?: string
  status?: string
  created_at: string
  updated_at: string
}

/**
 * Get authentication token from localStorage
 */
function getAuthHeaders(): HeadersInit {
  const token = localStorage.getItem('auth_token')
  const userId = localStorage.getItem('user_id')
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  }
  
  if (userId) {
    headers['X-User-ID'] = userId
  }
  
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  
  return headers
}

/**
 * Generic fetch wrapper with error handling
 */
async function apiFetch<T>(endpoint: string, options: RequestInit = {}): Promise<ApiResponse<T>> {
  const url = `${API_BASE_URL}${endpoint}`
  
  const defaultOptions: RequestInit = {
    headers: getAuthHeaders(),
  }
  
  // Merge provided options with defaults
  const config = {
    ...defaultOptions,
    ...options,
    headers: {
      ...defaultOptions.headers,
      ...(options.headers || {}),
    },
  }
  
  try {
    const response = await fetch(url, config)
    const data = await response.json()
    
    if (!response.ok) {
      return {
        error: data.error || 'An error occurred',
      }
    }
    
    return {
      data: data as T,
      message: data.message,
    }
  } catch (error) {
    console.error('API request failed:', error)
    return {
      error: error instanceof Error ? error.message : 'Network error',
    }
  }
}

/**
 * Set authentication tokens after login
 */
export function setAuthTokens(userId: string, token?: string) {
  localStorage.setItem('user_id', userId)
  if (token) {
    localStorage.setItem('auth_token', token)
  }
}

/**
 * Clear authentication tokens on logout
 */
export function clearAuthTokens() {
  localStorage.removeItem('user_id')
  localStorage.removeItem('auth_token')
}

/**
 * Get current user ID from localStorage
 */
export function getCurrentUserId(): string | null {
  return localStorage.getItem('user_id')
}

export { apiFetch }
export default apiFetch

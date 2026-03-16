/**
 * Projects API Service
 * Handles all project-related API calls
 */

import apiFetch, { type ApiResponse, type Project } from './api'

/**
 * Get all projects for the current user
 */
export async function getProjects(): Promise<ApiResponse<{ projects: Project[] }>> {
  return await apiFetch<{ projects: Project[] }>('/api/projects')
}

/**
 * Get a single project by ID
 */
export async function getProject(projectId: string): Promise<ApiResponse<{ project: Project }>> {
  return await apiFetch<{ project: Project }>(`/api/projects/${projectId}`)
}

/**
 * Create a new project
 */
export async function createProject(projectData: {
  name: string
  package_name: string
  platform?: string
  color?: string
  icon_url?: string
  version_code?: number
  version_name?: string
  workspace_xml?: string
  generated_code?: string
}): Promise<ApiResponse<{ project: Project }>> {
  return await apiFetch<{ project: Project }>('/api/projects', {
    method: 'POST',
    body: JSON.stringify(projectData),
  })
}

/**
 * Update an existing project
 */
export async function updateProject(
  projectId: string,
  updates: {
    name?: string
    package_name?: string
    platform?: string
    color?: string
    icon_url?: string
    version_code?: number
    version_name?: string
    workspace_xml?: string
    generated_code?: string
    status?: string
  }
): Promise<ApiResponse<{ project: Project }>> {
  return await apiFetch<{ project: Project }>(`/api/projects/${projectId}`, {
    method: 'PUT',
    body: JSON.stringify(updates),
  })
}

/**
 * Delete a project
 */
export async function deleteProject(projectId: string): Promise<ApiResponse<void>> {
  return await apiFetch<void>(`/api/projects/${projectId}`, {
    method: 'DELETE',
  })
}

/**
 * Save workspace XML and generated code to a project
 */
export async function saveWorkspace(
  projectId: string,
  workspaceXml: string,
  generatedCode: string
): Promise<ApiResponse<{ project: Project }>> {
  return await updateProject(projectId, {
    workspace_xml: workspaceXml,
    generated_code: generatedCode,
  })
}

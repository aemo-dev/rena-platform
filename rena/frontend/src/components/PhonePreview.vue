<script setup lang="ts">
import { ref } from 'vue'
import { 
  Smartphone, 
  ZoomIn, 
  ZoomOut, 
  RotateCcw} from 'lucide-vue-next'

interface Props {
  project?: {
    name: string
    color: string
  }
}

const props = defineProps<Props>()

// Phone state (removed - using simple rectangle only)
const zoomLevel = ref(1)
const minZoom = 0.5
const maxZoom = 1.5
const zoomStep = 0.1

// Phone dimensions
const screenWidth = ref(320)
const screenHeight = ref(580)

// Preset sizes
const presetSizes = [
  { name: 'iPhone SE', width: 375, height: 667 },
  { name: 'iPhone 14', width: 390, height: 844 },
  { name: 'iPhone 14 Pro Max', width: 430, height: 899 },
  { name: 'Pixel 7', width: 412, height: 915 },
  { name: 'iPad Mini', width: 768, height: 1024 },
  { name: 'iPad Pro', width: 1024, height: 1366 }
]

const setPresetSize = (preset: typeof presetSizes[0]) => {
  screenWidth.value = preset.width
  screenHeight.value = preset.height
}

const resetToDefault = () => {
  screenWidth.value = 320
  screenHeight.value = 580
}

// Computed device info (removed - no longer needed)

// Zoom controls
const zoomIn = () => {
  if (zoomLevel.value < maxZoom) {
    zoomLevel.value = Math.min(maxZoom, zoomLevel.value + zoomStep)
  }
}

const zoomOut = () => {
  if (zoomLevel.value > minZoom) {
    zoomLevel.value = Math.max(minZoom, zoomLevel.value - zoomStep)
  }
}

const resetZoom = () => {
  zoomLevel.value = 1
}
</script>

<template>
  <aside class="phone-preview">
    <!-- Header with Zoom Controls -->
    <div class="preview-header">
      <div class="header-left">
        <Smartphone :size="18" />
        <span class="header-title">Live Preview</span>
      </div>
      
      <!-- Zoom Controls -->
      <div class="zoom-controls">
        <button 
          @click="zoomOut" 
          :disabled="zoomLevel <= minZoom" 
          title="Zoom Out"
          class="zoom-btn"
        >
          <ZoomOut :size="16" />
        </button>
        <span class="zoom-level">{{ Math.round(zoomLevel * 100) }}%</span>
        <button 
          @click="zoomIn" 
          :disabled="zoomLevel >= maxZoom" 
          title="Zoom In"
          class="zoom-btn"
        >
          <ZoomIn :size="16" />
        </button>
        <button 
          @click="resetZoom" 
          class="zoom-btn reset-zoom" 
          title="Reset Zoom"
        >
          <RotateCcw :size="14" />
        </button>
      </div>
    </div>
    
    <!-- Phone Container - Simple Rectangle -->
    <div class="phone-wrapper">
      <div 
        class="phone-screen simple-rectangle" 
        :style="{
          width: `${screenWidth}px`,
          height: `${screenHeight}px`,
          transform: `scale(${zoomLevel})`
        }"
      >
        <div class="app-header" :style="{ background: project?.color }">
          <span class="app-title">{{ project?.name }}</span>
        </div>
        <div class="app-content">
          <div class="preview-placeholder">
            <p>Building your app...</p>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Size Controls - Below Phone -->
    <div class="size-controls">
      <div class="control-group">
        <label>Width (px)</label>
        <input 
          type="number" 
          v-model.number="screenWidth" 
          min="200" 
          max="1200"
          class="size-input"
        />
      </div>
      <div class="control-group">
        <label>Height (px)</label>
        <input 
          type="number" 
          v-model.number="screenHeight" 
          min="300" 
          max="1600"
          class="size-input"
        />
      </div>
      <button @click="resetToDefault" class="reset-size-btn" title="Reset to Default">
        <RotateCcw :size="14" />
      </button>
    </div>
    
    <!-- Preset Sizes - Single Line -->
    <div class="preset-sizes">
      <button 
        v-for="preset in presetSizes" 
        :key="preset.name"
        @click="setPresetSize(preset)"
        class="preset-btn"
        :class="{ active: screenWidth === preset.width && screenHeight === preset.height }"
        :title="`${preset.width}×${preset.height}`"
      >
        {{ preset.name }}
      </button>
    </div>
    
    <!-- Device Info Footer (removed - simple mode) -->
  </aside>
</template>

<style scoped>
.phone-preview {
  width: 360px;
  background: var(--surface);
  border-left: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

/* Header */
.preview-header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 700;
  font-size: 13px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Zoom Controls */
.zoom-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.zoom-btn {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: none;
  background: var(--bg-alt);
  color: var(--text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.zoom-btn:hover:not(:disabled) {
  background: var(--primary);
  color: white;
}

.zoom-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.zoom-level {
  font-size: 10px;
  font-weight: 700;
  min-width: 35px;
  text-align: center;
  color: var(--text-muted);
}

.reset-zoom {
  width: 20px !important;
  height: 20px !important;
}

/* Phone Wrapper */
.phone-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-alt);
  padding: 12px;
  overflow: hidden;
  min-height: 200px;
}

/* Simple Rectangle Screen */
.phone-screen.simple-rectangle {
  background: white;
  border-radius: 12px;
  box-shadow: 0 6px 24px rgba(0,0,0,0.1), 0 2px 6px rgba(0,0,0,0.06);
  border: 1px solid #e0e0e0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: transform 0.3s ease;
  transform-origin: center center;
}

/* Size Controls - Below Phone */
.size-controls {
  display: flex;
  gap: 8px;
  padding: 8px 12px;
  border-top: 1px solid var(--border);
  background: var(--bg-alt);
  align-items: center;
  flex-wrap: wrap;
}

.control-group {
  display: flex;
  align-items: center;
  gap: 6px;
}

.control-group label {
  font-size: 9px;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
}

.size-input {
  width: 65px;
  padding: 4px 6px;
  border: 1px solid var(--border);
  border-radius: 4px;
  background: var(--surface);
  color: var(--text);
  font-size: 11px;
  font-weight: 600;
}

.size-input:focus {
  outline: none;
  border-color: var(--primary);
}

.reset-size-btn {
  width: 28px;
  height: 28px;
  border-radius: 4px;
  border: none;
  background: var(--surface);
  color: var(--text-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.reset-size-btn:hover {
  background: var(--primary);
  color: white;
}

/* Preset Sizes - Single Line */
.preset-sizes {
  display: flex;
  gap: 6px;
  padding: 6px 12px 10px;
  border-top: 1px solid var(--border);
  background: var(--bg-alt);
  flex-wrap: wrap;
  justify-content: center;
}

.preset-btn {
  padding: 4px 10px;
  border-radius: 4px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text-muted);
  cursor: pointer;
  font-size: 10px;
  font-weight: 600;
  transition: all 0.2s;
  white-space: nowrap;
}

.preset-btn:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.preset-btn.active {
  border-color: var(--primary);
  background: var(--primary);
  color: white;
}

.reset-zoom {
  width: 24px !important;
  height: 24px !important;
}

/* Phone Wrapper */
.phone-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-alt);
  padding: 20px;
  overflow: hidden;
  min-height: 400px;
}

/* Simple Rectangle Screen */
.phone-screen.simple-rectangle {
  width: 320px;
  height: 580px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.12), 0 2px 8px rgba(0,0,0,0.08);
  border: 1px solid #e0e0e0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: transform 0.3s ease;
  transform-origin: center center;
}

.app-header {
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  color: white;
}

.app-title {
  font-weight: 700;
  font-size: 16px;
}

.app-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
}

.preview-placeholder {
  text-align: center;
}

/* Device Info Footer (removed - simple mode) */
</style>

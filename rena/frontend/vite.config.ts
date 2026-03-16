import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  envDir: '../../', // Look for .env in the root
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@languages': path.resolve(__dirname, '../languages'),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})

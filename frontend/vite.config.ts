import { defineConfig } from 'vite'
import { VitePWA } from 'vite-plugin-pwa'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
      react(),
      VitePWA({ registerType: 'autoUpdate' })
  ],
  server: {
    proxy: {
      '/api': {
        target:'http://192.168.2.162:42160',
        changeOrigin: false,
        secure: false,
        ws: true,
      },
    },
  }
})

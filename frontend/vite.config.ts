import { defineConfig } from 'vite'
import { VitePWA } from 'vite-plugin-pwa'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
      react(),
      VitePWA({
        registerType: 'autoUpdate',
        manifest: {
          name: 'AltWebInstaller',
          short_name: 'AltWebInstaller',
          theme_color: '#0f7e82'
        }
      })
  ],
  build:{
    outDir:"../dist"
  },
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

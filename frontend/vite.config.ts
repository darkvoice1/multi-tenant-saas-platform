import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { visualizer } from 'rollup-plugin-visualizer'

export default defineConfig({
  plugins: (() => {
    const plugins = [vue()]
    if (process.env.ANALYZE) {
      plugins.push(
        visualizer({
          filename: 'dist/stats.html',
          open: true,
          gzipSize: true,
          brotliSize: true
        })
      )
    }
    return plugins
  })(),
  server: {
    port: 5173
  }
})

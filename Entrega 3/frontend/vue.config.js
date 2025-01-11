const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: [
    'vuetify'
  ],
  devServer: {
    port: process.env.VUE_APP_PORT || 3000, // Usa la variable de entorno o un valor predeterminado
    proxy: {
      '/api': {
        target: process.env.VUE_APP_API_URL || 'http://localhost:8080', // Usa la variable de entorno o un valor predeterminado
        changeOrigin: true
      }
    }
  }
})

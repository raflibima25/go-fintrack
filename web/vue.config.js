const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  outputDir: 'dist',
  devServer: {
    proxy: {
      '/api': {
        target: process.env.VUE_APP_API_BASE_URL,
        changeOrigin: true
      }
    }
  },
  css: {
    loaderOptions: {
      postcss: {
        postcssOptions: {
          plugins: [
              require('tailwindcss'),
              require('autoprefixer')
          ]
        }
      }
    }
  },
  configureWebpack: {
    resolve: {
      alias: {
        '@': require('path').resolve(__dirname, 'src')
      }
    }
  }
})
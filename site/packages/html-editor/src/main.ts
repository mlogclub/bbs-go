import { createApp } from 'vue'
import App from './App.vue'

// 检测系统主题并设置
const setThemeClass = () => {
  if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

// 初始设置
setThemeClass()

// 监听系统主题变化
if (window.matchMedia) {
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', setThemeClass)
}

const app = createApp(App)
app.mount('#app')

import { defineStore } from 'pinia'
import { ref, onMounted } from 'vue'

type Theme = 'light' | 'dark'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref<Theme>('light')

  const applyTheme = (newTheme: Theme) => {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    if (newTheme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  const toggleTheme = () => {
    const newTheme = theme.value === 'light' ? 'dark' : 'light'
    applyTheme(newTheme)
  }

  onMounted(() => {
    const savedTheme = localStorage.getItem('theme') as Theme | null
    if (savedTheme) {
      applyTheme(savedTheme)
    } else {
      // Check system preference
      const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
      applyTheme(prefersDark ? 'dark' : 'light')
    }
  })

  return {
    theme,
    toggleTheme,
    applyTheme
  }
})
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { useThemeStore } from '../../src/stores/theme'
import ThemeToggler from '../../src/components/ThemeToggler.vue'

describe('ThemeToggler.vue', () => {
  it('calls the theme store to toggle the theme on click', async () => {
    const pinia = createPinia()
    setActivePinia(pinia)
    const themeStore = useThemeStore()
    const toggleThemeSpy = vi.spyOn(themeStore, 'toggleTheme')

    const wrapper = mount(ThemeToggler, {
      global: {
        plugins: [pinia],
      },
    })

    await wrapper.find('input[type="checkbox"]').trigger('click')

    expect(toggleThemeSpy).toHaveBeenCalledTimes(1)
  })
})
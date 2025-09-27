import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { axe } from 'vitest-axe'
import App from '../../src/App.vue'
import router from '../../src/router'
import { createPinia } from 'pinia'

// Mocking the ResizeObserver which is not available in JSDOM
const mockResizeObserver = class {
  observe() {}
  unobserve() {}
  disconnect() {}
}
global.ResizeObserver = mockResizeObserver

describe('App.vue', () => {
  it('should have no accessibility violations', async () => {
    const pinia = createPinia()
    const wrapper = mount(App, {
      global: {
        plugins: [router, pinia],
      },
    })

    // Wait for router to be ready to avoid issues with async navigation
    await router.isReady()

    const results = await axe(wrapper.element)
    expect(results).toHaveNoViolations()
  })
})
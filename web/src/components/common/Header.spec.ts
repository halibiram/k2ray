import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import Header from './Header.vue'
import { useAuthStore } from '../../stores/auth'
import { useThemeStore } from '../../stores/theme'

// Define a single, shared mock for the router push function
const mockRouterPush = vi.fn()

// Mock the entire vue-router module
vi.mock('vue-router', () => ({
  // The mock useRouter function now returns the same object every time
  useRouter: () => ({
    push: mockRouterPush,
  }),
}))

// Create a basic i18n instance for testing
const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en: {
      header: { title: 'K2Ray Panel' },
      common: { logout: 'Logout' },
    },
    tr: {
      header: { title: 'K2Ray Paneli' },
      common: { logout: 'Çıkış Yap' },
    },
  },
})

describe('Header.vue', () => {
  beforeEach(() => {
    // Set up a fresh Pinia instance for each test
    setActivePinia(createPinia())
    // Clear all mock history before each test
    vi.clearAllMocks()
  })

  it('renders the header title correctly', () => {
    const wrapper = mount(Header, {
      global: { plugins: [i18n] },
    })
    expect(wrapper.find('h2').text()).toBe('K2Ray Panel')
  })

  it('displays the correct theme toggle button text', async () => {
    const themeStore = useThemeStore()
    const wrapper = mount(Header, {
      global: { plugins: [i18n] },
    })

    expect(wrapper.find('button.bg-gray-200').text()).toBe('Dark')

    themeStore.theme = 'dark'
    await wrapper.vm.$nextTick()

    expect(wrapper.find('button.bg-gray-700').text()).toBe('Light')
  })

  it('calls logout and redirects when the logout button is clicked', async () => {
    const authStore = useAuthStore()
    // Mock the store's logout action to resolve immediately
    const logoutMock = vi.fn(() => Promise.resolve())
    authStore.logout = logoutMock

    const wrapper = mount(Header, {
      global: { plugins: [i18n] },
    })

    const logoutButton = wrapper.findAll('button').find(b => b.text() === 'Logout')
    expect(logoutButton).toBeDefined()

    // Trigger the click
    await logoutButton.trigger('click')

    // Wait for Vue to process the update queue after the promise resolves
    await wrapper.vm.$nextTick()

    // Assert that the store's logout method was called
    expect(logoutMock).toHaveBeenCalledTimes(1)

    // Assert that our shared mock router push function was called
    expect(mockRouterPush).toHaveBeenCalledWith('/login')
  })

  it('changes language when the select is changed', async () => {
    const wrapper = mount(Header, {
      global: { plugins: [i18n] },
    })

    expect(wrapper.find('h2').text()).toBe('K2Ray Panel')

    const select = wrapper.find('select')
    await select.setValue('tr')

    expect(i18n.global.locale.value).toBe('tr')
    expect(wrapper.find('h2').text()).toBe('K2Ray Paneli')
  })
})
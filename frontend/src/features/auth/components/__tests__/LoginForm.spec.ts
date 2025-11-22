import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import LoginForm from '../LoginForm.vue'
import { useAuthStore } from '../../stores/authStore'
import { createRouter, createWebHistory } from 'vue-router'

// Mock router
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: { template: '<div>Home</div>' } },
    { path: '/dashboard', component: { template: '<div>Dashboard</div>' } },
    { path: '/register', component: { template: '<div>Register</div>' } },
  ],
})

describe('LoginForm', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should render login form', () => {
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [router],
      },
    })

    expect(wrapper.find('input#email').exists()).toBe(true)
    expect(wrapper.find('input#password').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('should submit form with email and password', async () => {
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [router],
      },
    })

    const authStore = useAuthStore()
    authStore.login = vi.fn().mockResolvedValue({
      user: {
        id: '123',
        email: 'test@example.com',
        full_name: 'Test User',
        role: 'student',
      },
      token: 'mock-token',
    })

    // Fill form
    await wrapper.find('input#email').setValue('test@example.com')
    await wrapper.find('input#password').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    expect(authStore.login).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'password123',
    })
  })

  it('should show error message on login failure', async () => {
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [router],
      },
    })

    const authStore = useAuthStore()
    authStore.login = vi.fn().mockRejectedValue(new Error('Invalid credentials'))

    await wrapper.find('input#email').setValue('test@example.com')
    await wrapper.find('input#password').setValue('wrongpassword')

    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('Invalid credentials')
  })

  it('should disable submit button while loading', async () => {
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [router],
      },
    })

    const authStore = useAuthStore()
    let resolveLogin: (value: any) => void
    authStore.login = vi.fn().mockImplementation(() => {
      return new Promise(resolve => {
        resolveLogin = resolve
      })
    })

    await wrapper.find('input#email').setValue('test@example.com')
    await wrapper.find('input#password').setValue('password123')

    const submitButton = wrapper.find('button[type="submit"]')
    
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    expect(submitButton.attributes('disabled')).toBeDefined()
    expect(wrapper.text()).toContain('Signing in...')

    resolveLogin!({
      user: { id: '123', email: 'test@example.com', full_name: 'Test User', role: 'student' },
      token: 'mock-token',
    })
  })

  it('should have link to registration page', () => {
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [router],
      },
    })

    const registerLink = wrapper.find('a[href="/register"]')
    expect(registerLink.exists()).toBe(true)
    expect(registerLink.text()).toContain('Register')
  })

  it('should redirect to dashboard on successful login', async () => {
    const wrapper = mount(LoginForm, {
      global: {
        plugins: [router],
      },
    })

    const authStore = useAuthStore()
    authStore.login = vi.fn().mockResolvedValue({
      user: {
        id: '123',
        email: 'test@example.com',
        full_name: 'Test User',
        role: 'student',
      },
      token: 'mock-token',
    })

    const pushSpy = vi.spyOn(router, 'push')

    await wrapper.find('input#email').setValue('test@example.com')
    await wrapper.find('input#password').setValue('password123')
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    expect(pushSpy).toHaveBeenCalledWith('/dashboard')
  })
})

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import RegisterForm from '../RegisterForm.vue'
import { useAuthStore } from '../../stores/authStore'
import { createRouter, createWebHistory } from 'vue-router'

// Mock router
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: { template: '<div>Home</div>' } },
    { path: '/dashboard', component: { template: '<div>Dashboard</div>' } },
    { path: '/login', component: { template: '<div>Login</div>' } },
  ],
})

describe('RegisterForm', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should render registration form without role selection', () => {
    const wrapper = mount(RegisterForm, {
      global: {
        plugins: [router],
      },
    })

    expect(wrapper.find('input#full_name').exists()).toBe(true)
    expect(wrapper.find('input#email').exists()).toBe(true)
    expect(wrapper.find('input#password').exists()).toBe(true)

    // Role selection should NOT exist
    expect(wrapper.find('select#role').exists()).toBe(false)
    expect(wrapper.find('label[for="role"]').exists()).toBe(false)
  })

  it('should submit form with only required fields', async () => {
    const wrapper = mount(RegisterForm, {
      global: {
        plugins: [router],
      },
    })

    const authStore = useAuthStore()
    authStore.register = vi.fn().mockResolvedValue({
      user: {
        id: '123',
        email: 'test@example.com',
        full_name: 'Test User',
        role: 'student', // Assigned by backend
      },
      token: 'mock-token',
    })

    // Fill form
    await wrapper.find('input#full_name').setValue('Test User')
    await wrapper.find('input#email').setValue('test@example.com')
    await wrapper.find('input#password').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')

    await wrapper.vm.$nextTick()

    expect(authStore.register).toHaveBeenCalledWith({
      full_name: 'Test User',
      email: 'test@example.com',
      password: 'password123',
    })
  })

  it('should show error message on registration failure', async () => {
    const wrapper = mount(RegisterForm, {
      global: {
        plugins: [router],
      },
    })

    const authStore = useAuthStore()
    authStore.register = vi.fn().mockRejectedValue(new Error('Email already exists'))

    await wrapper.find('input#full_name').setValue('Test User')
    await wrapper.find('input#email').setValue('test@example.com')
    await wrapper.find('input#password').setValue('password123')

    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 100))

    expect(wrapper.text()).toContain('Email already exists')
  })

  it('should disable submit button while loading', async () => {
    const wrapper = mount(RegisterForm, {
      global: {
        plugins: [router],
      },
    })

    const authStore = useAuthStore()
    let resolveRegister: (value: any) => void
    authStore.register = vi.fn().mockImplementation(() => {
      return new Promise(resolve => {
        resolveRegister = resolve
      })
    })

    await wrapper.find('input#full_name').setValue('Test User')
    await wrapper.find('input#email').setValue('test@example.com')
    await wrapper.find('input#password').setValue('password123')

    const submitButton = wrapper.find('button[type="submit"]')

    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    expect(submitButton.attributes('disabled')).toBeDefined()
    expect(wrapper.text()).toContain('Creating account...')

    resolveRegister!({
      user: { id: '123', email: 'test@example.com', full_name: 'Test User', role: 'student' },
      token: 'mock-token',
    })
  })

  it('should have link to login page', () => {
    const wrapper = mount(RegisterForm, {
      global: {
        plugins: [router],
      },
    })

    const loginLink = wrapper.find('a[href="/login"]')
    expect(loginLink.exists()).toBe(true)
    expect(loginLink.text()).toContain('Sign in')
  })
})

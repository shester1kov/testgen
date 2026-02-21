import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import WelcomeView from '../WelcomeView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        { path: '/', component: WelcomeView },
        { path: '/login', component: { template: '<div>Login</div>' } },
        { path: '/register', component: { template: '<div>Register</div>' } }
    ]
})

describe('WelcomeView.vue', () => {
    it('renders standard layout elements correctly', async () => {
        router.push('/')
        await router.isReady()

        const wrapper = mount(WelcomeView, {
            global: {
                plugins: [router]
            }
        })

        // Check project name
        expect(wrapper.text()).toContain('TestGen')
        // Check main headline text in Russian
        expect(wrapper.text()).toContain('Автоматическая генерация')
        expect(wrapper.text()).toContain('учебных тестов')
        // Check main call to action
        expect(wrapper.text()).toContain('Начать использование')
    })

    it('contains links to login and register pages', async () => {
        router.push('/')
        await router.isReady()

        const wrapper = mount(WelcomeView, {
            global: {
                plugins: [router],
                stubs: {
                    RouterLink: {
                        props: ['to'],
                        template: '<a :href="to" class="router-link-stub"><slot /></a>'
                    }
                }
            }
        })

        const links = wrapper.findAll('.router-link-stub')
        const toValues = links.map(link => link.attributes('href'))

        // Verify correct router links exist
        expect(toValues).toContain('/login')
        expect(toValues).toContain('/register')

        // Verify link text
        const loginLink = links.find(link => link.attributes('href') === '/login')
        const registerLink = links.find(link => link.text().includes('Регистрация'))

        expect(loginLink?.text()).toContain('Вход')
        expect(registerLink?.attributes('href')).toBe('/register')
    })
})

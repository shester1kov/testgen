import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/features/auth/stores/authStore'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/features/auth/components/LoginForm.vue'),
    meta: { requiresAuth: false, layout: 'auth' },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/features/auth/components/RegisterForm.vue'),
    meta: { requiresAuth: false, layout: 'auth' },
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/DashboardView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/documents',
    name: 'Documents',
    component: () => import('@/views/DocumentsView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/documents/:id',
    name: 'DocumentDetails',
    component: () => import('@/views/DocumentDetailsView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/tests',
    name: 'Tests',
    component: () => import('@/views/TestsView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/tests/create',
    name: 'CreateTest',
    component: () => import('@/views/CreateTestView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/tests/:id',
    name: 'TestDetails',
    component: () => import('@/views/TestDetailsView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/tests/:id/edit',
    name: 'EditTest',
    component: () => import('@/views/EditTestView.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFoundView.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if ((to.name === 'Login' || to.name === 'Register') && authStore.isAuthenticated) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router

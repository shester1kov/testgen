import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, VueWrapper } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import UserList from '../UserList.vue'
import { useUsersStore } from '../../stores/usersStore'
import { useAuthStore } from '@/features/auth/stores/authStore'
import type { User } from '@/features/auth/types/auth.types'

vi.mock('@/services/userService')
vi.mock('@/services/authService')
vi.mock('@/utils/logger', () => ({
  default: {
    logStoreAction: vi.fn(),
    logStoreError: vi.fn(),
    info: vi.fn(),
    debug: vi.fn(),
  },
}))

describe('UserList', () => {
  let wrapper: VueWrapper
  let usersStore: ReturnType<typeof useUsersStore>
  let authStore: ReturnType<typeof useAuthStore>

  const mockUsers: User[] = [
    {
      id: '1',
      email: 'admin@test.com',
      full_name: 'Admin User',
      role: 'admin',
      created_at: '2024-01-01T00:00:00Z',
    },
    {
      id: '2',
      email: 'teacher@test.com',
      full_name: 'Teacher User',
      role: 'teacher',
      created_at: '2024-01-02T00:00:00Z',
    },
    {
      id: '3',
      email: 'student@test.com',
      full_name: 'Student User',
      role: 'student',
      created_at: '2024-01-03T00:00:00Z',
    },
  ]

  beforeEach(() => {
    setActivePinia(createPinia())
    usersStore = useUsersStore()
    authStore = useAuthStore()
    vi.clearAllMocks()
  })

  // POSITIVE TEST: Renders users list
  it('renders users list correctly', async () => {
    usersStore.users = mockUsers
    usersStore.total = 3
    authStore.user = mockUsers[0] // Admin user

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('admin@test.com')
    expect(wrapper.text()).toContain('teacher@test.com')
    expect(wrapper.text()).toContain('student@test.com')
  })

  // POSITIVE TEST: Display role badges
  it('displays role badges for each user', async () => {
    usersStore.users = mockUsers
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Admin')
    expect(wrapper.text()).toContain('Teacher')
    expect(wrapper.text()).toContain('Student')
  })

  // POSITIVE TEST: Admin can see change role buttons
  it('shows change role button for admin users', async () => {
    usersStore.users = mockUsers
    authStore.user = mockUsers[0] // Admin user

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const changeRoleButtons = wrapper.findAll('button').filter(btn => btn.text().includes('Change Role'))
    expect(changeRoleButtons.length).toBe(3) // One for each user
  })

  // NEGATIVE TEST: Non-admin cannot change roles
  it('hides change role buttons for non-admin users', async () => {
    usersStore.users = mockUsers
    authStore.user = mockUsers[1] // Teacher user

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const changeRoleButtons = wrapper.findAll('button').filter(btn => btn.text().includes('Change Role'))
    expect(changeRoleButtons.length).toBe(0)
  })

  // POSITIVE TEST: Fetches users on mount
  it('fetches users on component mount', async () => {
    const fetchUsersSpy = vi.spyOn(usersStore, 'fetchUsers').mockResolvedValue()
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    expect(fetchUsersSpy).toHaveBeenCalledWith(1)
  })

  // POSITIVE TEST: Format role names correctly
  it('formats role names correctly', async () => {
    usersStore.users = mockUsers
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const vm = wrapper.vm as any
    expect(vm.formatRole('admin')).toBe('Admin')
    expect(vm.formatRole('teacher')).toBe('Teacher')
    expect(vm.formatRole('student')).toBe('Student')
  })

  // POSITIVE TEST: Get correct badge classes for each role
  it('returns correct badge classes for roles', async () => {
    usersStore.users = mockUsers
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const vm = wrapper.vm as any
    expect(vm.getRoleBadgeClass('admin')).toContain('pink')
    expect(vm.getRoleBadgeClass('teacher')).toContain('orange')
    expect(vm.getRoleBadgeClass('student')).toContain('blue')
  })

  // NEGATIVE TEST: Display loading state
  it('displays loading state correctly', async () => {
    usersStore.loading = true
    usersStore.users = []
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Loading users')
  })

  // NEGATIVE TEST: Display error message
  it('displays error message when present', async () => {
    vi.spyOn(usersStore, 'fetchUsers').mockRejectedValue(new Error('Failed to load users'))
    usersStore.error = 'Failed to load users'
    usersStore.users = [] // Error state shows when there are no users loaded yet
    usersStore.loading = false
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Failed to load users')
    expect(wrapper.text()).toContain('Retry')
  })

  // POSITIVE TEST: Display empty state
  it('displays empty state when no users', async () => {
    vi.spyOn(usersStore, 'fetchUsers').mockResolvedValue()
    usersStore.users = []
    usersStore.total = 0
    usersStore.loading = false
    usersStore.error = null
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('No users found')
  })

  // POSITIVE TEST: Pagination controls visible when needed
  it('shows pagination when total pages > 1', async () => {
    usersStore.users = mockUsers
    usersStore.total = 50 // More than one page (20 per page)
    usersStore.currentPage = 1
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Previous')
    expect(wrapper.text()).toContain('Next')
  })

  // POSITIVE TEST: Handle page change
  it('changes page when pagination button clicked', async () => {
    const fetchUsersSpy = vi.spyOn(usersStore, 'fetchUsers').mockResolvedValue()
    usersStore.users = mockUsers
    usersStore.total = 50
    usersStore.currentPage = 1
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const nextButton = wrapper.findAll('button').find(btn => btn.text().includes('Next'))
    await nextButton?.trigger('click')
    await wrapper.vm.$nextTick()

    expect(fetchUsersSpy).toHaveBeenCalledWith(2)
  })

  // NEGATIVE TEST: Previous button disabled on first page
  it('disables previous button on first page', async () => {
    usersStore.users = mockUsers
    usersStore.total = 50
    usersStore.currentPage = 1
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const prevButton = wrapper.findAll('button').find(btn => btn.text().includes('Previous'))
    expect(prevButton?.attributes('disabled')).toBeDefined()
  })

  // POSITIVE TEST: Open role change modal
  it('opens role change modal when button clicked', async () => {
    usersStore.users = mockUsers
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const changeRoleButtons = wrapper.findAll('button').filter(btn => btn.text().includes('Change Role'))
    await changeRoleButtons[0].trigger('click')
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Change User Role')
  })

  // POSITIVE TEST: Close role change modal
  it('closes role change modal', async () => {
    usersStore.users = mockUsers
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const vm = wrapper.vm as any
    vm.openRoleModal(mockUsers[1])
    await wrapper.vm.$nextTick()

    expect(vm.showRoleModal).toBe(true)

    vm.closeRoleModal()
    await wrapper.vm.$nextTick()

    expect(vm.showRoleModal).toBe(false)
    expect(vm.selectedUser).toBeNull()
  })

  // POSITIVE TEST: Successfully update user role
  it('updates user role successfully', async () => {
    const updateRoleSpy = vi.spyOn(usersStore, 'updateUserRole').mockResolvedValue(mockUsers[1])
    usersStore.users = mockUsers
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const vm = wrapper.vm as any
    vm.openRoleModal(mockUsers[2]) // Student user
    vm.selectedRole = 'teacher'
    await wrapper.vm.$nextTick()

    await vm.handleRoleChange()
    await wrapper.vm.$nextTick()

    expect(updateRoleSpy).toHaveBeenCalledWith('3', 'teacher')
    expect(vm.showRoleModal).toBe(false)
  })

  // NEGATIVE TEST: Handle role update error
  it('handles role update errors gracefully', async () => {
    const updateRoleSpy = vi
      .spyOn(usersStore, 'updateUserRole')
      .mockRejectedValue(new Error('Update failed'))
    usersStore.users = mockUsers
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const vm = wrapper.vm as any
    vm.openRoleModal(mockUsers[2])
    vm.selectedRole = 'teacher'
    await wrapper.vm.$nextTick()

    await vm.handleRoleChange()
    await wrapper.vm.$nextTick()

    expect(updateRoleSpy).toHaveBeenCalled()
    // Modal should close even on error (based on finally block)
    expect(vm.isProcessing).toBe(false)
  })

  // POSITIVE TEST: Display user count
  it('displays correct user count', async () => {
    usersStore.users = mockUsers
    usersStore.total = 3
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('3 users total')
  })

  // POSITIVE TEST: Retry on error
  it('retries loading users when retry button clicked', async () => {
    const fetchUsersSpy = vi.spyOn(usersStore, 'fetchUsers').mockResolvedValue()
    const clearErrorSpy = vi.spyOn(usersStore, 'clearError')
    usersStore.error = 'Failed to load users'
    usersStore.users = mockUsers
    authStore.user = mockUsers[0]

    wrapper = mount(UserList)
    await wrapper.vm.$nextTick()

    const retryButton = wrapper.findAll('button').find(btn => btn.text() === 'Retry')
    await retryButton?.trigger('click')
    await wrapper.vm.$nextTick()

    expect(clearErrorSpy).toHaveBeenCalled()
    expect(fetchUsersSpy).toHaveBeenCalled()
  })
})

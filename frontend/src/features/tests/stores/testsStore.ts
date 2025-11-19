import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Test, TestGenerationRequest, TestExportRequest, MoodleSyncRequest } from '../types/test.types'
import { testService } from '@/services/testService'

export const useTestsStore = defineStore('tests', () => {
  // State
  const tests = ref<Test[]>([])
  const currentTest = ref<Test | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)
  const currentPage = ref(1)
  const totalPages = ref(0)

  // Actions
  async function createTest(data: Partial<Test>) {
    loading.value = true
    error.value = null

    try {
      const test = await testService.createTest(data)
      tests.value.unshift(test)
      return test
    } catch (err: any) {
      error.value = err.message || 'Failed to create test'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchTests(page = 1, limit = 10) {
    loading.value = true
    error.value = null

    try {
      const response = await testService.getTests(page, limit)
      tests.value = response.data
      total.value = response.total
      currentPage.value = response.page
      totalPages.value = response.totalPages
      return response
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch tests'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchTest(id: string) {
    loading.value = true
    error.value = null

    try {
      const test = await testService.getTest(id)
      currentTest.value = test
      return test
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch test'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateTest(id: string, data: Partial<Test>) {
    loading.value = true
    error.value = null

    try {
      const test = await testService.updateTest(id, data)
      // Update test in list
      const index = tests.value.findIndex(t => t.id === id)
      if (index !== -1) {
        tests.value[index] = test
      }
      if (currentTest.value?.id === id) {
        currentTest.value = test
      }
      return test
    } catch (err: any) {
      error.value = err.message || 'Failed to update test'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteTest(id: string) {
    loading.value = true
    error.value = null

    try {
      await testService.deleteTest(id)
      tests.value = tests.value.filter(test => test.id !== id)
    } catch (err: any) {
      error.value = err.message || 'Failed to delete test'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function generateTest(data: TestGenerationRequest) {
    loading.value = true
    error.value = null

    try {
      const test = await testService.generateTest(data)
      tests.value.unshift(test)
      currentTest.value = test
      return test
    } catch (err: any) {
      error.value = err.message || 'Failed to generate test'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function exportTest(data: TestExportRequest) {
    loading.value = true
    error.value = null

    try {
      const blob = await testService.exportTest(data)
      return blob
    } catch (err: any) {
      error.value = err.message || 'Failed to export test'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function syncToMoodle(data: MoodleSyncRequest) {
    loading.value = true
    error.value = null

    try {
      const test = await testService.syncToMoodle(data)
      // Update test in list
      const index = tests.value.findIndex(t => t.id === data.test_id)
      if (index !== -1) {
        tests.value[index] = test
      }
      if (currentTest.value?.id === data.test_id) {
        currentTest.value = test
      }
      return test
    } catch (err: any) {
      error.value = err.message || 'Failed to sync to Moodle'
      throw err
    } finally {
      loading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  return {
    // State
    tests,
    currentTest,
    loading,
    error,
    total,
    currentPage,
    totalPages,
    // Actions
    createTest,
    fetchTests,
    fetchTest,
    updateTest,
    deleteTest,
    generateTest,
    exportTest,
    syncToMoodle,
    clearError,
  }
})

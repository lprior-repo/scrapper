/**
 * Simple All Components Tests
 * 
 * Basic import and export tests for all components
 * Quick validation that components are properly structured
 */

import { test, expect } from "bun:test"

test("OrganizationInput exports correctly", () => {
  const { OrganizationInput } = require('../OrganizationInput')
  expect(typeof OrganizationInput).toBe('function')
})

test("TopicsToggle exports correctly", () => {
  const { TopicsToggle } = require('../TopicsToggle')  
  expect(typeof TopicsToggle).toBe('function')
})

test("LoadButton exports correctly", () => {
  const { LoadButton } = require('../LoadButton')
  expect(typeof LoadButton).toBe('function')
})

test("AppHeader exports correctly", () => {
  const { AppHeader } = require('../AppHeader')
  expect(typeof AppHeader).toBe('function')
})

test("ErrorBoundary exports correctly", () => {
  const { ErrorBoundary, GraphErrorBoundary } = require('../ErrorBoundary')
  expect(typeof ErrorBoundary).toBe('function')
  expect(typeof GraphErrorBoundary).toBe('function')
})

test("LoadingSpinner components export correctly", () => {
  const { 
    LoadingSpinner, 
    SkeletonLoader, 
    GraphLoadingSpinner,
    InlineLoader,
    LoadingOverlay 
  } = require('../LoadingSpinner')
  
  expect(typeof LoadingSpinner).toBe('function')
  expect(typeof SkeletonLoader).toBe('function')
  expect(typeof GraphLoadingSpinner).toBe('function')
  expect(typeof InlineLoader).toBe('function')
  expect(typeof LoadingOverlay).toBe('function')
})

test("all components are properly structured", () => {
  // This test verifies that our component refactoring was successful
  const orgInput = require('../OrganizationInput')
  const topics = require('../TopicsToggle')
  const loadBtn = require('../LoadButton')
  const header = require('../AppHeader')
  const errors = require('../ErrorBoundary')
  const loading = require('../LoadingSpinner')
  
  // All modules should be defined
  expect(orgInput).toBeDefined()
  expect(topics).toBeDefined()
  expect(loadBtn).toBeDefined()
  expect(header).toBeDefined()
  expect(errors).toBeDefined()
  expect(loading).toBeDefined()
  
  // All should export their main components
  expect(orgInput.OrganizationInput).toBeDefined()
  expect(topics.TopicsToggle).toBeDefined()
  expect(loadBtn.LoadButton).toBeDefined()
  expect(header.AppHeader).toBeDefined()
  expect(errors.ErrorBoundary).toBeDefined()
  expect(loading.LoadingSpinner).toBeDefined()
})
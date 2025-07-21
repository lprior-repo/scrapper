/**
 * Simplified Test Setup for Bun Test Runner
 * 
 * This file is preloaded by bunfig.toml for all tests
 * Sets up minimal testing environment for speed
 */

// Mock DOM globals for Node environment
if (typeof window === 'undefined') {
  global.window = {
    matchMedia: (query: string) => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: () => {},
      removeListener: () => {},
      addEventListener: () => {},
      removeEventListener: () => {},
      dispatchEvent: () => {},
    }),
    location: { reload: () => {} },
  } as any

  global.document = {
    body: { contains: () => true },
    createElement: () => ({}),
    querySelector: () => null,
    querySelectorAll: () => [],
  } as any

  global.navigator = {
    clipboard: {
      writeText: async () => Promise.resolve(),
      readText: async () => Promise.resolve(''),
    },
  } as any
}

// Mock ResizeObserver
global.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}

// Mock IntersectionObserver  
global.IntersectionObserver = class IntersectionObserver {
  constructor(public readonly callback: any) {}
  observe() {}
  unobserve() {}
  disconnect() {}
}

// Setup jest mock functions for compatibility
global.jest = {
  fn: (impl?: (...args: readonly any[]) => any) => {
    const calls: readonly any[] = []
    const mockFn = (...args: readonly any[]) => {
      calls.push(args)
      return impl ? impl(...args) : undefined
    }
    
    mockFn.mockImplementation = (newImpl: any) => { impl = newImpl; return mockFn }
    mockFn.mockReturnValue = (value: any) => { impl = () => value; return mockFn }
    mockFn.mockResolvedValue = (value: any) => { impl = () => Promise.resolve(value); return mockFn }
    mockFn.mockRejectedValue = (value: any) => { impl = () => Promise.reject(value); return mockFn }
    mockFn.mockClear = () => { calls.length = 0; return mockFn }
    
    Object.defineProperty(mockFn, 'calls', { get: () => calls })
    
    return mockFn as any
  },
  clearAllMocks: () => {},
}

declare global {
  var jest: any
}
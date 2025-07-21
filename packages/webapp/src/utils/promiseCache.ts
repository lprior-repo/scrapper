/**
 * Promise Cache Utility for React 19 Suspense
 *
 * This utility prevents re-triggering promises on every render
 * and includes cache eviction for failed promises.
 * Works with React 19's `use` hook.
 */

interface CacheEntry<T> {
  readonly promise: Promise<T>
  readonly status: 'pending' | 'resolved' | 'rejected'
  readonly result?: T
  readonly error?: Error
  readonly timestamp: number
}

interface CacheOptions {
  readonly ttl?: number // Time to live in milliseconds
  readonly maxSize?: number // Maximum cache size
  readonly retryOnError?: boolean // Retry failed promises
}

class PromiseCache {
  private readonly cache = new Map<string, CacheEntry<unknown>>()
  private readonly options: Required<CacheOptions>

  constructor(options: CacheOptions = {}) {
    this.options = {
      ttl: options.ttl ?? 5 * 60 * 1000, // 5 minutes default
      maxSize: options.maxSize ?? 100,
      retryOnError: options.retryOnError ?? true,
    }
  }

  /**
   * Get or create a cached promise
   */
  get<T>(key: string, promiseFactory: () => Promise<T>): Promise<T> {
    const existing = this.cache.get(key) as CacheEntry<T> | undefined
    const now = Date.now()

    // Check if we have a valid cached entry
    if (existing) {
      const isExpired = now - existing.timestamp > this.options.ttl
      const shouldRetry =
        existing.status === 'rejected' && this.options.retryOnError

      if (!isExpired && !shouldRetry) {
        return existing.promise
      }

      // Remove expired or failed entries
      if (isExpired || shouldRetry) {
        this.cache.delete(key)
      }
    }

    // Create new promise and cache entry
    const promise = promiseFactory()
    const entry: CacheEntry<T> = {
      promise,
      status: 'pending',
      timestamp: now,
    }

    // Update entry status when promise resolves/rejects
    promise
      .then((result) => {
        const cachedEntry = this.cache.get(key) as CacheEntry<T> | undefined
        if (cachedEntry && cachedEntry.promise === promise) {
          this.cache.set(key, {
            ...cachedEntry,
            status: 'resolved',
            result,
          })
        }
      })
      .catch((error) => {
        const cachedEntry = this.cache.get(key) as CacheEntry<T> | undefined
        if (cachedEntry && cachedEntry.promise === promise) {
          this.cache.set(key, {
            ...cachedEntry,
            status: 'rejected',
            error: error instanceof Error ? error : new Error(String(error)),
          })
        }

        // Auto-evict failed promises if retryOnError is enabled
        if (this.options.retryOnError) {
          setTimeout(() => {
            const currentEntry = this.cache.get(key)
            if (currentEntry && currentEntry.status === 'rejected') {
              this.cache.delete(key)
            }
          }, 1000) // 1 second delay before eviction
        }
      })

    this.cache.set(key, entry)
    this.enforceMaxSize()

    return promise
  }

  /**
   * Check if a key exists in cache
   */
  has(key: string): boolean {
    return this.cache.has(key)
  }

  /**
   * Get cache entry status without triggering promise creation
   */
  getStatus(key: string): CacheEntry<unknown>['status'] | null {
    const entry = this.cache.get(key)
    return entry?.status ?? null
  }

  /**
   * Manually evict a cache entry
   */
  evict(key: string): boolean {
    return this.cache.delete(key)
  }

  /**
   * Clear all cache entries
   */
  clear(): void {
    this.cache.clear()
  }

  /**
   * Evict expired entries
   */
  evictExpired(): number {
    const now = Date.now()
    let evictedCount = 0

    for (const [key, entry] of this.cache.entries()) {
      if (now - entry.timestamp > this.options.ttl) {
        this.cache.delete(key)
        evictedCount++
      }
    }

    return evictedCount
  }

  /**
   * Get cache statistics
   */
  getStats() {
    const entries = Array.from(this.cache.values())
    const now = Date.now()

    return {
      totalEntries: this.cache.size,
      pendingEntries: entries.filter((entry) => entry.status === 'pending')
        .length,
      resolvedEntries: entries.filter((entry) => entry.status === 'resolved')
        .length,
      rejectedEntries: entries.filter((entry) => entry.status === 'rejected')
        .length,
      expiredEntries: entries.filter(
        (entry) => now - entry.timestamp > this.options.ttl
      ).length,
    }
  }

  /**
   * Enforce maximum cache size by evicting oldest entries
   */
  private enforceMaxSize(): void {
    if (this.cache.size <= this.options.maxSize) {
      return
    }

    // Sort by timestamp and remove oldest entries
    const entries = Array.from(this.cache.entries()).sort(
      ([, a], [, b]) => a.timestamp - b.timestamp
    )

    const entriesToRemove = this.cache.size - this.options.maxSize

    for (let i = 0; i < entriesToRemove; i++) {
      const [key] = entries[i]
      this.cache.delete(key)
    }
  }
}

// Default cache instance
const defaultCache = new PromiseCache({
  ttl: 5 * 60 * 1000, // 5 minutes
  maxSize: 100,
  retryOnError: true,
})

/**
 * Hook-friendly cache function for React 19 Suspense
 */
export const getCachedPromise = <T>(
  key: string,
  promiseFactory: () => Promise<T>,
  cache: PromiseCache = defaultCache
): Promise<T> => {
  return cache.get(key, promiseFactory)
}

/**
 * Create a new cache instance with custom options
 */
export const createPromiseCache = (options?: CacheOptions): PromiseCache => {
  return new PromiseCache(options)
}

/**
 * Access the default cache instance
 */
export const getDefaultCache = (): PromiseCache => defaultCache

/**
 * Create a cache key from parameters
 */
export const createCacheKey = (
  ...parts: readonly (string | number | boolean | null | undefined)[]
): string => {
  return parts
    .filter((part) => part != null)
    .map((part) => String(part))
    .join(':')
}

/**
 * Preload a promise into cache without suspending
 */
export const preloadPromise = <T>(
  key: string,
  promiseFactory: () => Promise<T>,
  cache: PromiseCache = defaultCache
): void => {
  if (!cache.has(key)) {
    cache.get(key, promiseFactory)
  }
}

export { PromiseCache }
export type { CacheEntry, CacheOptions }

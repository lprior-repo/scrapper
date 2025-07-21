import { Context, Effect, Layer, Duration, Console } from 'effect'
import { ParseResult } from '@effect/schema'
import { UnknownException } from 'effect/Cause'
import {
  GraphResponseSchema,
  validateApiResponseSync,
  type GraphNode,
  type GraphEdge,
  type GraphResponse,
} from '@overseer/shared'
import {
  TracingService,
  type TraceContext,
  type ApiCallTrace,
  type PerformanceMetric,
  TracingError,
} from './tracing'

// Re-export types for compatibility
export type { GraphNode, GraphEdge, GraphResponse }

// Observability types
export interface ApiMetrics {
  readonly successCount: number
  readonly errorCount: number
  readonly totalRequests: number
  readonly averageResponseTime: number
  readonly lastRequestTime: number
}

export interface RequestContext {
  readonly correlationId: string
  readonly timestamp: number
  readonly url: string
  readonly method: string
  readonly headers: Record<string, string>
}

export interface ResponseContext {
  readonly status: number
  readonly duration: number
  readonly size?: number
  readonly headers: Record<string, string>
}

export interface ApiCallMetadata {
  readonly request: RequestContext
  readonly response?: ResponseContext
  readonly error?: string
  readonly retryCount: number
}

// Utility functions
export const extractTopicsFromNodes = (
  nodes: readonly GraphNode[]
): readonly GraphTopic[] => {
  return nodes
    .filter((node) => node.type === 'topic')
    .map((node) => ({
      name: node.data.name as string,
      count: node.data.count as number,
    }))
}

export const extractNodesByType = (
  nodes: readonly GraphNode[],
  type: string
): readonly GraphNode[] => {
  return nodes.filter((node) => node.type === type)
}

// Data sanitization utilities for secure logging
export const sanitizeHeaders = (
  headers: Record<string, string>
): Record<string, string> => {
  const sensitiveHeaders = [
    'authorization',
    'cookie',
    'x-api-key',
    'x-auth-token',
  ]
  const sanitized = { ...headers }

  Object.keys(sanitized).forEach((key) => {
    if (sensitiveHeaders.includes(key.toLowerCase())) {
      sanitized[key] = '[REDACTED]'
    }
  })

  return sanitized
}

export const sanitizeUrl = (url: string): string => {
  try {
    const urlObj = new URL(url)
    // Remove sensitive query parameters
    const sensitiveParams = ['token', 'key', 'secret', 'password', 'auth']
    sensitiveParams.forEach((param) => {
      if (urlObj.searchParams.has(param)) {
        urlObj.searchParams.set(param, '[REDACTED]')
      }
    })
    return urlObj.toString()
  } catch {
    return url // Return original if parsing fails
  }
}

export const truncateData = (data: unknown, maxLength = 1000): string => {
  const serialized = typeof data === 'string' ? data : JSON.stringify(data)
  return serialized.length > maxLength
    ? `${serialized.substring(0, maxLength)}...[TRUNCATED]`
    : serialized
}

// Performance monitoring helpers
export const measureOperation = <T, E>(
  operation: Effect.Effect<T, E>,
  operationName: string
): Effect.Effect<T, E | TracingError> =>
  Effect.gen(function* () {
    const startTime = yield* Effect.sync(() => Date.now())
    const observabilityService = yield* ObservabilityService

    const result = yield* Effect.catchAll(operation, (error) =>
      Effect.gen(function* () {
        const duration = Date.now() - startTime
        yield* observabilityService.recordMetrics(
          operationName,
          duration,
          false
        )
        yield* Effect.fail(error)
      })
    )

    const duration = Date.now() - startTime
    yield* observabilityService.recordMetrics(operationName, duration, true)

    return result
  })

// Retry with exponential backoff and observability
export const retryWithObservability = <T, E>(
  operation: Effect.Effect<T, E>,
  operationName: string,
  maxRetries = 3
): Effect.Effect<T, E | TracingError> =>
  Effect.gen(function* () {
    const observabilityService = yield* ObservabilityService
    let lastError: E | null = null

    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        const result = yield* operation

        if (attempt > 0) {
          yield* Effect.gen(function* () {
            yield* Console.log(
              `Operation ${operationName} succeeded on attempt ${attempt + 1}`
            )
          })
        }

        return result
      } catch (error) {
        lastError = error as E

        if (attempt < maxRetries) {
          const delay = Math.min(1000 * Math.pow(2, attempt), 10000) // Cap at 10 seconds
          yield* Effect.gen(function* () {
            yield* Console.log(
              `Operation ${operationName} failed on attempt ${attempt + 1}, retrying in ${delay}ms`
            )
          })
          yield* Effect.sleep(Duration.millis(delay))
        }
      }
    }

    // All retries exhausted
    yield* Effect.gen(function* () {
      yield* Console.log(
        `Operation ${operationName} failed after ${maxRetries + 1} attempts`
      )
    })

    yield* Effect.fail(lastError!)
  })

// Circuit breaker state management
interface CircuitBreakerState {
  readonly failures: number
  readonly lastFailureTime: number
  readonly state: 'closed' | 'open' | 'half-open'
}

const circuitBreakerStates = new Map<string, CircuitBreakerState>()

export const withCircuitBreaker = <T, E>(
  operation: Effect.Effect<T, E>,
  operationName: string,
  threshold = 5,
  timeout = 30000
): Effect.Effect<T, E | Error> =>
  Effect.gen(function* () {
    const state = circuitBreakerStates.get(operationName) || {
      failures: 0,
      lastFailureTime: 0,
      state: 'closed' as const,
    }

    const now = Date.now()

    // Check if circuit should transition from open to half-open
    if (state.state === 'open' && now - state.lastFailureTime > timeout) {
      state.state = 'half-open'
      circuitBreakerStates.set(operationName, state)
    }

    // Reject if circuit is open
    if (state.state === 'open') {
      yield* Effect.fail(
        new Error(`Circuit breaker is open for ${operationName}`)
      )
    }

    try {
      const result = yield* operation

      // Reset on success
      if (state.failures > 0) {
        circuitBreakerStates.set(operationName, {
          failures: 0,
          lastFailureTime: 0,
          state: 'closed',
        })
      }

      return result
    } catch (error) {
      // Record failure
      const newState: CircuitBreakerState = {
        failures: state.failures + 1,
        lastFailureTime: now,
        state: state.failures + 1 >= threshold ? 'open' : 'closed',
      }

      circuitBreakerStates.set(operationName, newState)

      yield* Effect.fail(error as E)
    }
  })

// Graph topic type for utility functions
export interface GraphTopic {
  readonly name: string
  readonly count: number
}

// Observability Service Interface
export interface ObservabilityService {
  readonly generateCorrelationId: () => Effect.Effect<string, never>
  readonly createRequestContext: (
    url: string,
    method: string,
    headers?: Record<string, string>
  ) => Effect.Effect<RequestContext, never>
  readonly logApiCall: (
    metadata: ApiCallMetadata
  ) => Effect.Effect<void, TracingError>
  readonly recordMetrics: (
    operationName: string,
    duration: number,
    success: boolean
  ) => Effect.Effect<void, TracingError>
  readonly getMetrics: () => Effect.Effect<ApiMetrics, never>
}

// API Client Service
export interface ApiClient {
  readonly getGraph: (
    org: string,
    useTopics?: boolean
  ) => Effect.Effect<
    GraphResponse,
    ParseResult.ParseError | UnknownException | TracingError
  >
}

export const ObservabilityService = Context.GenericTag<ObservabilityService>(
  'ObservabilityService'
)
export const ApiClient = Context.GenericTag<ApiClient>('ApiClient')

// Observability Service Implementation
export const ObservabilityServiceLive = Layer.succeed(
  ObservabilityService,
  ObservabilityService.of({
    generateCorrelationId: () =>
      Effect.sync(() => crypto.randomUUID().replace(/-/g, '')),

    createRequestContext: (url: string, method: string, headers = {}) =>
      Effect.gen(function* () {
        const correlationId = yield* Effect.sync(() =>
          crypto.randomUUID().replace(/-/g, '')
        )
        const timestamp = yield* Effect.sync(() => Date.now())

        return {
          correlationId,
          timestamp,
          url,
          method,
          headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
            'X-Correlation-ID': correlationId,
            'X-Request-Timestamp': timestamp.toString(),
            ...headers,
          },
        }
      }),

    logApiCall: (metadata: ApiCallMetadata) =>
      Effect.gen(function* () {
        const tracingService = yield* TracingService
        const context = yield* tracingService.getCurrentContext()

        if (context) {
          const apiTrace: ApiCallTrace = {
            method: metadata.request.method,
            url: metadata.request.url,
            duration: metadata.response?.duration || 0,
            responseStatus: metadata.response?.status,
            error: metadata.error,
          }

          yield* tracingService.recordApiCall(context, apiTrace)
        }

        // Enhanced structured logging
        const logLevel = metadata.error ? 'error' : 'info'
        const logMessage = `API ${metadata.request.method} ${metadata.request.url}`

        yield* Console.log(`[${logLevel.toUpperCase()}] ${logMessage}`, {
          correlationId: metadata.request.correlationId,
          method: metadata.request.method,
          url: metadata.request.url,
          status: metadata.response?.status,
          duration: metadata.response?.duration,
          retryCount: metadata.retryCount,
          timestamp: new Date().toISOString(),
          ...(metadata.error && { error: metadata.error }),
        })
      }),

    recordMetrics: (
      operationName: string,
      duration: number,
      success: boolean
    ) =>
      Effect.gen(function* () {
        const tracingService = yield* TracingService
        const context = yield* tracingService.getCurrentContext()

        if (context) {
          const metric: PerformanceMetric = {
            name: `api.${operationName}`,
            value: duration,
            unit: 'ms',
            metadata: {
              success,
              operation: operationName,
              timestamp: Date.now(),
            },
          }

          yield* tracingService.recordPerformanceMetric(context, metric)
        }

        // Store metrics for aggregation
        yield* Effect.sync(() => {
          const metricsKey = `api_metrics_${operationName}`
          const existing = JSON.parse(localStorage.getItem(metricsKey) || '[]')
          existing.push({ duration, success, timestamp: Date.now() })

          // Keep only last 100 entries
          const recent = existing.slice(-100)
          localStorage.setItem(metricsKey, JSON.stringify(recent))
        })
      }),

    getMetrics: () =>
      Effect.gen(function* () {
        return yield* Effect.sync(() => {
          const allMetrics = Object.keys(localStorage)
            .filter((key) => key.startsWith('api_metrics_'))
            .flatMap((key) => {
              const metrics = JSON.parse(localStorage.getItem(key) || '[]')
              return metrics
            })

          if (allMetrics.length === 0) {
            return {
              successCount: 0,
              errorCount: 0,
              totalRequests: 0,
              averageResponseTime: 0,
              lastRequestTime: 0,
            }
          }

          const successCount = allMetrics.filter((m) => m.success).length
          const errorCount = allMetrics.filter((m) => !m.success).length
          const totalRequests = allMetrics.length
          const averageResponseTime =
            allMetrics.reduce((sum, m) => sum + m.duration, 0) / totalRequests
          const lastRequestTime = Math.max(
            ...allMetrics.map((m) => m.timestamp)
          )

          return {
            successCount,
            errorCount,
            totalRequests,
            averageResponseTime,
            lastRequestTime,
          }
        })
      }),
  })
)

// Enhanced API Client with Comprehensive Observability
export const ApiClientLive = Layer.effect(
  ApiClient,
  Effect.gen(function* () {
    const observabilityService = yield* ObservabilityService
    const tracingService = yield* TracingService

    return ApiClient.of({
      getGraph: (org: string, useTopics?: boolean) =>
        Effect.gen(function* () {
          const baseUrl = 'http://localhost:8081'
          const url = `${baseUrl}/api/graph/${org}${useTopics ? '?useTopics=true' : ''}`
          const operationName = 'getGraph'
          const startTime = Date.now()

          // Create request context with correlation ID
          const requestContext =
            yield* observabilityService.createRequestContext(url, 'GET')

          // Start tracing span
          const traceContext = yield* tracingService.getCurrentContext()
          let spanContext: TraceContext | null = null

          if (traceContext) {
            spanContext = yield* tracingService.createSpan(
              traceContext,
              `api.${operationName}`,
              {
                'http.method': 'GET',
                'http.url': url,
                'api.operation': operationName,
                'graph.organization': org,
                'graph.use_topics': useTopics || false,
              }
            )
          }

          // Enhanced error handling with observability
          const executeRequest = Effect.gen(function* () {
            // Make HTTP request with enhanced headers
            const httpResponse = yield* Effect.tryPromise({
              try: () =>
                fetch(url, {
                  method: 'GET',
                  headers: requestContext.headers,
                }),
              catch: (error) => {
                const errorMessage = `Network error: ${error instanceof Error ? error.message : String(error)}`

                // Record error in tracing
                if (spanContext) {
                  Effect.runSync(
                    tracingService.recordError(
                      spanContext,
                      new Error(errorMessage),
                      {
                        'error.type': 'network',
                        'http.url': url,
                        'api.operation': operationName,
                      }
                    )
                  )
                }

                return new Error(errorMessage)
              },
            })

            // Check response status
            if (!httpResponse.ok) {
              const errorMessage = `HTTP error! status: ${httpResponse.status}`

              // Record HTTP error
              if (spanContext) {
                yield* tracingService.recordError(
                  spanContext,
                  new Error(errorMessage),
                  {
                    'error.type': 'http',
                    'http.status_code': httpResponse.status,
                    'http.url': url,
                    'api.operation': operationName,
                  }
                )
              }

              yield* Effect.fail(new Error(errorMessage))
            }

            // Parse JSON response
            const apiJsonData = yield* Effect.tryPromise({
              try: () => httpResponse.json(),
              catch: (error) => {
                const errorMessage = `JSON parsing error: ${error instanceof Error ? error.message : String(error)}`

                // Record parsing error
                if (spanContext) {
                  Effect.runSync(
                    tracingService.recordError(
                      spanContext,
                      new Error(errorMessage),
                      {
                        'error.type': 'parsing',
                        'http.url': url,
                        'api.operation': operationName,
                      }
                    )
                  )
                }

                return new Error(errorMessage)
              },
            })

            // Create response context
            const responseContext: ResponseContext = {
              status: httpResponse.status,
              duration: Date.now() - startTime,
              size: JSON.stringify(apiJsonData).length,
              headers: Object.fromEntries(httpResponse.headers.entries()),
            }

            // Validate response using shared schema
            const validatedResponse = validateApiResponseSync(
              GraphResponseSchema,
              apiJsonData,
              `Graph API response for ${org}`
            )

            // Log successful API call
            yield* observabilityService.logApiCall({
              request: requestContext,
              response: responseContext,
              retryCount: 0,
            })

            // Record success metrics
            yield* observabilityService.recordMetrics(
              operationName,
              responseContext.duration,
              true
            )

            // Record additional performance metrics
            if (spanContext) {
              yield* tracingService.recordPerformanceMetric(spanContext, {
                name: 'api.response_size',
                value: responseContext.size || 0,
                unit: 'bytes',
                metadata: {
                  operation: operationName,
                  organization: org,
                },
              })

              yield* tracingService.recordPerformanceMetric(spanContext, {
                name: 'api.nodes_count',
                value: validatedResponse.nodes.length,
                unit: 'count',
                metadata: {
                  operation: operationName,
                  organization: org,
                },
              })

              yield* tracingService.recordPerformanceMetric(spanContext, {
                name: 'api.edges_count',
                value: validatedResponse.edges.length,
                unit: 'count',
                metadata: {
                  operation: operationName,
                  organization: org,
                },
              })
            }

            return validatedResponse
          })

          // Execute with error recovery
          const result = yield* Effect.catchAll(executeRequest, (error) =>
            Effect.gen(function* () {
              const duration = Date.now() - startTime

              // Log failed API call
              yield* observabilityService.logApiCall({
                request: requestContext,
                error: error instanceof Error ? error.message : String(error),
                retryCount: 0,
              })

              // Record failure metrics
              yield* observabilityService.recordMetrics(
                operationName,
                duration,
                false
              )

              // Re-throw the error
              yield* Effect.fail(error)
            })
          )

          return result
        }),
    })
  })
).pipe(Layer.provide(ObservabilityServiceLive))

// Enhanced API Client with all observability features
export const ObservableApiClientLive = Layer.mergeAll(
  ObservabilityServiceLive,
  ApiClientLive
)

// Convenience function for creating observable API operations
export const createObservableOperation = <T, E>(
  operation: Effect.Effect<T, E>,
  operationName: string,
  options: {
    readonly enableRetry?: boolean
    readonly enableCircuitBreaker?: boolean
    readonly enablePerformanceTracking?: boolean
    readonly maxRetries?: number
    readonly circuitBreakerThreshold?: number
    readonly circuitBreakerTimeout?: number
  } = {}
): Effect.Effect<T, E | TracingError | Error> => {
  const {
    enableRetry = true,
    enableCircuitBreaker = true,
    enablePerformanceTracking = true,
    maxRetries = 3,
    circuitBreakerThreshold = 5,
    circuitBreakerTimeout = 30000,
  } = options

  let enhancedOperation = operation

  // Add performance tracking
  if (enablePerformanceTracking) {
    enhancedOperation = measureOperation(enhancedOperation, operationName)
  }

  // Add retry logic
  if (enableRetry) {
    enhancedOperation = retryWithObservability(
      enhancedOperation,
      operationName,
      maxRetries
    )
  }

  // Add circuit breaker
  if (enableCircuitBreaker) {
    enhancedOperation = withCircuitBreaker(
      enhancedOperation,
      operationName,
      circuitBreakerThreshold,
      circuitBreakerTimeout
    )
  }

  return enhancedOperation
}

// Frontend performance monitoring integration
export const recordFrontendMetrics = (): Effect.Effect<void, TracingError> =>
  Effect.gen(function* () {
    const tracingService = yield* TracingService
    const context = yield* tracingService.getCurrentContext()

    if (!context || typeof window === 'undefined') {
      return
    }

    // Record browser performance metrics
    yield* tracingService.recordBrowserMetrics(context)

    // Record performance navigation timing
    if ('performance' in window && window.performance.getEntriesByType) {
      const navigation = window.performance.getEntriesByType(
        'navigation'
      )[0] as PerformanceNavigationTiming

      if (navigation) {
        yield* tracingService.recordPerformanceMetric(context, {
          name: 'frontend.page_load_time',
          value: navigation.loadEventEnd - navigation.navigationStart,
          unit: 'ms',
          metadata: { source: 'navigation_timing' },
        })

        yield* tracingService.recordPerformanceMetric(context, {
          name: 'frontend.dom_content_loaded',
          value:
            navigation.domContentLoadedEventEnd - navigation.navigationStart,
          unit: 'ms',
          metadata: { source: 'navigation_timing' },
        })

        yield* tracingService.recordPerformanceMetric(context, {
          name: 'frontend.first_contentful_paint',
          value: navigation.loadEventStart - navigation.navigationStart,
          unit: 'ms',
          metadata: { source: 'navigation_timing' },
        })
      }
    }

    // Record resource timing if available
    if ('performance' in window && window.performance.getEntriesByType) {
      const resources = window.performance.getEntriesByType(
        'resource'
      ) as readonly PerformanceResourceTiming[]

      if (resources.length > 0) {
        const totalResourceTime = resources.reduce(
          (sum, resource) => sum + resource.duration,
          0
        )
        const averageResourceTime = totalResourceTime / resources.length

        yield* tracingService.recordPerformanceMetric(context, {
          name: 'frontend.average_resource_load_time',
          value: averageResourceTime,
          unit: 'ms',
          metadata: {
            source: 'resource_timing',
            resource_count: resources.length,
          },
        })
      }
    }
  })

// Utility for creating correlation ID headers for external requests
export const createCorrelationHeaders = (): Effect.Effect<
  Record<string, string>,
  never
> =>
  Effect.gen(function* () {
    const observabilityService = yield* ObservabilityService
    const correlationId = yield* observabilityService.generateCorrelationId()
    const timestamp = yield* Effect.sync(() => Date.now())

    return {
      'X-Correlation-ID': correlationId,
      'X-Request-Timestamp': timestamp.toString(),
      'X-Source': 'webapp-frontend',
      'X-User-Agent':
        typeof navigator !== 'undefined' ? navigator.userAgent : 'unknown',
    }
  })

// Note: All exports are already declared above with individual export statements

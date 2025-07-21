import { Context, Effect, Layer, Data, Runtime, Scope } from 'effect'
import { UnknownException } from 'effect/Cause'
import { z } from 'zod'

// ============================================================================
// Types and Schemas
// ============================================================================

export const TraceContextSchema = z.object({
  traceId: z.string(),
  spanId: z.string(),
  parentSpanId: z.string().optional(),
  correlationId: z.string(),
  timestamp: z.number(),
  userId: z.string().optional(),
  sessionId: z.string(),
})

export const ComponentEventSchema = z.object({
  componentName: z.string(),
  eventType: z.enum(['mount', 'unmount', 'render', 'error']),
  props: z.record(z.unknown()).optional(),
  renderTime: z.number().optional(),
  error: z.string().optional(),
})

export const UserInteractionSchema = z.object({
  eventType: z.enum(['click', 'submit', 'navigation', 'input']),
  target: z.string(),
  value: z.string().optional(),
  metadata: z.record(z.unknown()).optional(),
})

export const ApiCallTraceSchema = z.object({
  method: z.string(),
  url: z.string(),
  requestBody: z.string().optional(),
  responseStatus: z.number().optional(),
  responseBody: z.string().optional(),
  duration: z.number(),
  error: z.string().optional(),
})

export const PerformanceMetricSchema = z.object({
  name: z.string(),
  value: z.number(),
  unit: z.string(),
  metadata: z.record(z.unknown()).optional(),
})

export const BrowserMetricSchema = z.object({
  memoryUsage: z.number(),
  networkStatus: z.enum(['online', 'offline']),
  connectionType: z.string().optional(),
  viewportSize: z.object({ width: z.number(), height: z.number() }),
  userAgent: z.string(),
})

export type TraceContext = z.infer<typeof TraceContextSchema>
export type ComponentEvent = z.infer<typeof ComponentEventSchema>
export type UserInteraction = z.infer<typeof UserInteractionSchema>
export type ApiCallTrace = z.infer<typeof ApiCallTraceSchema>
export type PerformanceMetric = z.infer<typeof PerformanceMetricSchema>
export type BrowserMetric = z.infer<typeof BrowserMetricSchema>

// ============================================================================
// Errors
// ============================================================================

export const TracingError = Data.TaggedError('TracingError')<{
  readonly message: string
  readonly context?: string
  readonly originalError?: unknown
}>()

// ============================================================================
// Core Tracing Service Interface
// ============================================================================

export interface TracingService {
  readonly startTrace: (name: string, metadata?: Record<string, unknown>) => Effect.Effect<TraceContext, TracingError>
  readonly finishTrace: (context: TraceContext) => Effect.Effect<void, TracingError>
  readonly createSpan: (parentContext: TraceContext, name: string, metadata?: Record<string, unknown>) => Effect.Effect<TraceContext, TracingError>
  readonly recordComponentEvent: (context: TraceContext, event: ComponentEvent) => Effect.Effect<void, TracingError>
  readonly recordUserInteraction: (context: TraceContext, interaction: UserInteraction) => Effect.Effect<void, TracingError>
  readonly recordApiCall: (context: TraceContext, apiCall: ApiCallTrace) => Effect.Effect<void, TracingError>
  readonly recordPerformanceMetric: (context: TraceContext, metric: PerformanceMetric) => Effect.Effect<void, TracingError>
  readonly recordBrowserMetrics: (context: TraceContext) => Effect.Effect<void, TracingError>
  readonly recordError: (context: TraceContext, error: Error, metadata?: Record<string, unknown>) => Effect.Effect<void, TracingError>
  readonly generateCorrelationId: () => Effect.Effect<string, never>
  readonly getCurrentContext: () => Effect.Effect<TraceContext | null, never>
  readonly setCurrentContext: (context: TraceContext | null) => Effect.Effect<void, never>
}

export const TracingService = Context.GenericTag<TracingService>('TracingService')

// ============================================================================
// Utility Functions
// ============================================================================

const generateId = (): string => {
  return crypto.randomUUID().replace(/-/g, '')
}

const generateSessionId = (): string => {
  const stored = sessionStorage.getItem('tracing-session-id')
  if (stored) return stored
  
  const newId = generateId()
  sessionStorage.setItem('tracing-session-id', newId)
  return newId
}

const getBrowserMetrics = (): BrowserMetric => {
  const memory = (performance as any).memory
  const memoryUsage = memory ? memory.usedJSHeapSize : 0
  
  const connection = (navigator as any).connection
  const connectionType = connection ? connection.effectiveType : 'unknown'
  
  return {
    memoryUsage,
    networkStatus: navigator.onLine ? 'online' : 'offline',
    connectionType,
    viewportSize: {
      width: window.innerWidth,
      height: window.innerHeight,
    },
    userAgent: navigator.userAgent,
  }
}

const formatTraceForJaeger = (
  context: TraceContext,
  operationName: string,
  tags: Record<string, unknown> = {},
  logs: ReadonlyArray<{ readonly timestamp: number; readonly fields: Record<string, unknown> }> = []
) => {
  return {
    traceID: context.traceId,
    spanID: context.spanId,
    parentSpanID: context.parentSpanId,
    operationName,
    startTime: context.timestamp * 1000, // Jaeger expects microseconds
    duration: Date.now() * 1000 - context.timestamp * 1000,
    tags: [
      { key: 'component', value: 'frontend' },
      { key: 'correlation.id', value: context.correlationId },
      { key: 'session.id', value: context.sessionId },
      ...(context.userId ? [{ key: 'user.id', value: context.userId }] : []),
      ...Object.entries(tags).map(([key, value]) => ({ key, value: String(value) })),
    ],
    logs,
    process: {
      serviceName: 'webapp-frontend',
      tags: [
        { key: 'jaeger.version', value: 'browser-client' },
        { key: 'hostname', value: window.location.hostname },
      ],
    },
  }
}

// ============================================================================
// Live Implementation
// ============================================================================

export const TracingServiceLive = Layer.succeed(
  TracingService,
  TracingService.of({
    startTrace: (name: string, metadata?: Record<string, unknown>) =>
      Effect.gen(function* () {
        const traceId = generateId()
        const spanId = generateId()
        const correlationId = generateId()
        const timestamp = Date.now()
        const sessionId = generateSessionId()

        const context: TraceContext = {
          traceId,
          spanId,
          correlationId,
          timestamp,
          sessionId,
        }

        // Send to Jaeger if available
        yield* Effect.tryPromise({
          try: async () => {
            const jaegerSpan = formatTraceForJaeger(context, name, metadata)
            
            // In a real implementation, you would send this to your Jaeger collector
            // For now, we'll log it to console in development
            if (import.meta.env.DEV) {
              console.log('[Tracing] Started trace:', { name, context, metadata, jaegerSpan })
            }
            
            // Store in localStorage for debugging
            const traces = JSON.parse(localStorage.getItem('tracing-spans') || '[]')
            traces.push({ ...jaegerSpan, type: 'start' })
            localStorage.setItem('tracing-spans', JSON.stringify(traces.slice(-100))) // Keep last 100
          },
          catch: (error) => TracingError({
            message: 'Failed to start trace',
            context: name,
            originalError: error,
          }),
        })

        return context
      }),

    finishTrace: (context: TraceContext) =>
      Effect.gen(function* () {
        yield* Effect.tryPromise({
          try: async () => {
            const jaegerSpan = formatTraceForJaeger(context, 'trace-finished')
            
            if (import.meta.env.DEV) {
              console.log('[Tracing] Finished trace:', { context, jaegerSpan })
            }
            
            const traces = JSON.parse(localStorage.getItem('tracing-spans') || '[]')
            traces.push({ ...jaegerSpan, type: 'finish' })
            localStorage.setItem('tracing-spans', JSON.stringify(traces.slice(-100)))
          },
          catch: (error) => TracingError({
            message: 'Failed to finish trace',
            context: context.traceId,
            originalError: error,
          }),
        })
      }),

    createSpan: (parentContext: TraceContext, name: string, metadata?: Record<string, unknown>) =>
      Effect.gen(function* () {
        const spanId = generateId()
        const timestamp = Date.now()

        const spanContext: TraceContext = {
          ...parentContext,
          spanId,
          parentSpanId: parentContext.spanId,
          timestamp,
        }

        yield* Effect.tryPromise({
          try: async () => {
            const jaegerSpan = formatTraceForJaeger(spanContext, name, metadata)
            
            if (import.meta.env.DEV) {
              console.log('[Tracing] Created span:', { name, spanContext, metadata, jaegerSpan })
            }
            
            const traces = JSON.parse(localStorage.getItem('tracing-spans') || '[]')
            traces.push({ ...jaegerSpan, type: 'span' })
            localStorage.setItem('tracing-spans', JSON.stringify(traces.slice(-100)))
          },
          catch: (error) => TracingError({
            message: 'Failed to create span',
            context: name,
            originalError: error,
          }),
        })

        return spanContext
      }),

    recordComponentEvent: (context: TraceContext, event: ComponentEvent) =>
      Effect.gen(function* () {
        yield* Effect.tryPromise({
          try: async () => {
            const tags = {
              'component.name': event.componentName,
              'event.type': event.eventType,
              ...(event.renderTime && { 'render.time': event.renderTime }),
            }

            const logs = [{
              timestamp: Date.now() * 1000,
              fields: {
                event: 'component_event',
                ...event,
              },
            }]

            const jaegerSpan = formatTraceForJaeger(context, `component.${event.eventType}`, tags, logs)
            
            if (import.meta.env.DEV) {
              console.log('[Tracing] Component event:', { event, context, jaegerSpan })
            }

            const traces = JSON.parse(localStorage.getItem('tracing-spans') || '[]')
            traces.push({ ...jaegerSpan, type: 'component_event' })
            localStorage.setItem('tracing-spans', JSON.stringify(traces.slice(-100)))
          },
          catch: (error) => TracingError({
            message: 'Failed to record component event',
            context: context.traceId,
            originalError: error,
          }),
        })
      }),

    recordUserInteraction: (context: TraceContext, interaction: UserInteraction) =>
      Effect.gen(function* () {
        yield* Effect.tryPromise({
          try: async () => {
            const tags = {
              'interaction.type': interaction.eventType,
              'interaction.target': interaction.target,
              ...(interaction.value && { 'interaction.value': interaction.value }),
            }

            const logs = [{
              timestamp: Date.now() * 1000,
              fields: {
                event: 'user_interaction',
                ...interaction,
              },
            }]

            const jaegerSpan = formatTraceForJaeger(context, `user.${interaction.eventType}`, tags, logs)
            
            if (import.meta.env.DEV) {
              console.log('[Tracing] User interaction:', { interaction, context, jaegerSpan })
            }

            const traces = JSON.parse(localStorage.getItem('tracing-spans') || '[]')
            traces.push({ ...jaegerSpan, type: 'user_interaction' })
            localStorage.setItem('tracing-spans', JSON.stringify(traces.slice(-100)))
          },
          catch: (error) => TracingError({
            message: 'Failed to record user interaction',
            context: context.traceId,
            originalError: error,
          }),
        })
      }),

    recordApiCall: (context: TraceContext, apiCall: ApiCallTrace) =>
      Effect.gen(function* () {
        yield* Effect.tryPromise({
          try: async () => {
            const tags = {
              'http.method': apiCall.method,
              'http.url': apiCall.url,
              'http.status_code': apiCall.responseStatus || 0,
              'api.duration': apiCall.duration,
              ...(apiCall.error && { 'error': true, 'error.message': apiCall.error }),
            }

            const logs = [{
              timestamp: Date.now() * 1000,
              fields: {
                event: 'api_call',
                ...apiCall,
              },
            }]

            const jaegerSpan = formatTraceForJaeger(context, `http.${apiCall.method.toLowerCase()}`, tags, logs)
            
            if (import.meta.env.DEV) {
              console.log('[Tracing] API call:', { apiCall, context, jaegerSpan })
            }

            const traces = JSON.parse(localStorage.getItem('tracing-spans') || '[]')
            traces.push({ ...jaegerSpan, type: 'api_call' })
            localStorage.setItem('tracing-spans', JSON.stringify(traces.slice(-100)))
          },
          catch: (error) => TracingError({
            message: 'Failed to record API call',
            context: context.traceId,
            originalError: error,
          }),
        })
      }),

    recordPerformanceMetric: (context: TraceContext, metric: PerformanceMetric) =>
      Effect.gen(function* () {
        yield* Effect.tryPromise({
          try: async () => {
            const tags = {
              'metric.name': metric.name,
              'metric.value': metric.value,
              'metric.unit': metric.unit,
            }

            const logs = [{
              timestamp: Date.now() * 1000,
              fields: {
                event: 'performance_metric',
                ...metric,
              },
            }]

            const jaegerSpan = formatTraceForJaeger(context, `performance.${metric.name}`, tags, logs)
            
            if (import.meta.env.DEV) {
              console.log('[Tracing] Performance metric:', { metric, context, jaegerSpan })
            }

            const traces = JSON.parse(localStorage.getItem('tracing-spans') || '[]')
            traces.push({ ...jaegerSpan, type: 'performance_metric' })
            localStorage.setItem('tracing-spans', JSON.stringify(traces.slice(-100)))
          },
          catch: (error) => TracingError({
            message: 'Failed to record performance metric',
            context: context.traceId,
            originalError: error,
          }),
        })
      }),

    recordBrowserMetrics: (context: TraceContext) =>
      Effect.gen(function* () {
        yield* Effect.tryPromise({
          try: async () => {
            const browserMetric = getBrowserMetrics()
            
            const tags = {
              'browser.memory_usage': browserMetric.memoryUsage,
              'browser.network_status': browserMetric.networkStatus,
              'browser.connection_type': browserMetric.connectionType || 'unknown',
              'browser.viewport_width': browserMetric.viewportSize.width,
              'browser.viewport_height': browserMetric.viewportSize.height,
            }

            const logs = [{
              timestamp: Date.now() * 1000,
              fields: {
                event: 'browser_metrics',
                ...browserMetric,
              },
            }]

            const jaegerSpan = formatTraceForJaeger(context, 'browser.metrics', tags, logs)
            
            if (import.meta.env.DEV) {
              console.log('[Tracing] Browser metrics:', { browserMetric, context, jaegerSpan })
            }

            const traces = JSON.parse(localStorage.getItem('tracing-spans') || '[]')
            traces.push({ ...jaegerSpan, type: 'browser_metrics' })
            localStorage.setItem('tracing-spans', JSON.stringify(traces.slice(-100)))
          },
          catch: (error) => TracingError({
            message: 'Failed to record browser metrics',
            context: context.traceId,
            originalError: error,
          }),
        })
      }),

    recordError: (context: TraceContext, error: Error, metadata?: Record<string, unknown>) =>
      Effect.gen(function* () {
        yield* Effect.tryPromise({
          try: async () => {
            const tags = {
              'error': true,
              'error.kind': error.name,
              'error.message': error.message,
              'error.stack': error.stack || 'No stack trace available',
            }

            const logs = [{
              timestamp: Date.now() * 1000,
              fields: {
                event: 'error',
                'error.kind': error.name,
                'error.message': error.message,
                'error.stack': error.stack,
                ...metadata,
              },
            }]

            const jaegerSpan = formatTraceForJaeger(context, 'error', tags, logs)
            
            if (import.meta.env.DEV) {
              console.error('[Tracing] Error recorded:', { error, metadata, context, jaegerSpan })
            }

            const traces = JSON.parse(localStorage.getItem('tracing-spans') || '[]')
            traces.push({ ...jaegerSpan, type: 'error' })
            localStorage.setItem('tracing-spans', JSON.stringify(traces.slice(-100)))
          },
          catch: (tracingError) => TracingError({
            message: 'Failed to record error',
            context: context.traceId,
            originalError: tracingError,
          }),
        })
      }),

    generateCorrelationId: () =>
      Effect.succeed(generateId()),

    getCurrentContext: () =>
      Effect.succeed(
        JSON.parse(sessionStorage.getItem('current-trace-context') || 'null')
      ),

    setCurrentContext: (context: TraceContext | null) =>
      Effect.sync(() => {
        if (context) {
          sessionStorage.setItem('current-trace-context', JSON.stringify(context))
        } else {
          sessionStorage.removeItem('current-trace-context')
        }
      }),
  })
)

// ============================================================================
// React Integration Hooks and HOCs
// ============================================================================

export const withTracing = <P extends object>(
  WrappedComponent: React.ComponentType<P>,
  componentName: string
) => {
  return React.forwardRef<any, P>((props, ref) => {
    const [traceContext, setTraceContext] = React.useState<TraceContext | null>(null)
    const renderStartTime = React.useRef<number>(Date.now())

    React.useEffect(() => {
      const runtime = Runtime.defaultRuntime

      // Component mount tracing
      const mountEffect = Effect.gen(function* () {
        const tracingService = yield* TracingService
        const context = yield* tracingService.getCurrentContext()
        
        if (context) {
          const spanContext = yield* tracingService.createSpan(context, `component.${componentName}`)
          setTraceContext(spanContext)
          
          yield* tracingService.recordComponentEvent(spanContext, {
            componentName,
            eventType: 'mount',
            props: props as Record<string, unknown>,
          })
        }
      })

      Runtime.runPromise(runtime)(mountEffect).catch(console.error)

      return () => {
        // Component unmount tracing
        if (traceContext) {
          const unmountEffect = Effect.gen(function* () {
            const tracingService = yield* TracingService
            yield* tracingService.recordComponentEvent(traceContext, {
              componentName,
              eventType: 'unmount',
            })
          })

          Runtime.runPromise(runtime)(unmountEffect).catch(console.error)
        }
      }
    }, [])

    // Render time tracking
    React.useEffect(() => {
      if (traceContext) {
        const renderTime = Date.now() - renderStartTime.current
        const runtime = Runtime.defaultRuntime

        const renderEffect = Effect.gen(function* () {
          const tracingService = yield* TracingService
          yield* tracingService.recordComponentEvent(traceContext, {
            componentName,
            eventType: 'render',
            renderTime,
          })
        })

        Runtime.runPromise(runtime)(renderEffect).catch(console.error)
      }
      renderStartTime.current = Date.now()
    })

    return React.createElement(WrappedComponent, { ...props, ref })
  })
}

// ============================================================================
// API Integration
// ============================================================================

export const withApiTracing = <T,>(
  apiCall: () => Promise<T>,
  method: string,
  url: string
): Effect.Effect<T, TracingError | UnknownException> =>
  Effect.gen(function* () {
    const tracingService = yield* TracingService
    const context = yield* tracingService.getCurrentContext()
    
    if (!context) {
      // If no context, just run the API call without tracing
      return yield* Effect.tryPromise({
        try: apiCall,
        catch: (error) => new UnknownException(error),
      })
    }

    const startTime = Date.now()
    let apiTrace: ApiCallTrace = {
      method,
      url,
      duration: 0,
    }

    try {
      const result = yield* Effect.tryPromise({
        try: apiCall,
        catch: (error) => {
          apiTrace = {
            ...apiTrace,
            duration: Date.now() - startTime,
            error: error instanceof Error ? error.message : String(error),
          }
          return new UnknownException(error)
        },
      })

      apiTrace = {
        ...apiTrace,
        duration: Date.now() - startTime,
        responseStatus: 200, // Assume success if no error
      }

      yield* tracingService.recordApiCall(context, apiTrace)
      return result
    } catch (error) {
      apiTrace = {
        ...apiTrace,
        duration: Date.now() - startTime,
        error: error instanceof Error ? error.message : String(error),
      }

      yield* tracingService.recordApiCall(context, apiTrace)
      throw error
    }
  })

// ============================================================================
// Error Boundary Integration
// ============================================================================

export interface TracingErrorBoundaryProps {
  readonly children: React.ReactNode
  readonly fallback?: React.ComponentType<{ readonly error: Error; readonly resetError: () => void }>
}

export const TracingErrorBoundary: React.FC<TracingErrorBoundaryProps> = ({
  children,
  fallback: Fallback,
}) => {
  return (
    <ErrorBoundary
      fallback={Fallback}
      onError={(error, errorInfo) => {
        const runtime = Runtime.defaultRuntime
        
        const errorEffect = Effect.gen(function* () {
          const tracingService = yield* TracingService
          const context = yield* tracingService.getCurrentContext()
          
          if (context) {
            yield* tracingService.recordError(context, error, {
              componentStack: errorInfo.componentStack,
              errorBoundary: true,
            })
          }
        })

        Runtime.runPromise(runtime)(errorEffect).catch(console.error)
      }}
    >
      {children}
    </ErrorBoundary>
  )
}

// Simple error boundary implementation
class ErrorBoundary extends React.Component<
  TracingErrorBoundaryProps & { readonly onError?: (error: Error, errorInfo: React.ErrorInfo) => void },
  { readonly hasError: boolean; readonly error?: Error }
> {
  constructor(props: TracingErrorBoundaryProps & { readonly onError?: (error: Error, errorInfo: React.ErrorInfo) => void }) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    this.props.onError?.(error, errorInfo)
  }

  render() {
    if (this.state.hasError && this.state.error) {
      if (this.props.fallback) {
        const Fallback = this.props.fallback
        return <Fallback error={this.state.error} resetError={() => this.setState({ hasError: false, error: undefined })} />
      }
      
      return (
        <div style={{ padding: '2rem', border: '1px solid red', borderRadius: '4px', margin: '1rem' }}>
          <h2>Something went wrong</h2>
          <details>
            <summary>Error details</summary>
            <pre style={{ whiteSpace: 'pre-wrap', color: 'red' }}>
              {this.state.error.stack}
            </pre>
          </details>
          <button onClick={() => this.setState({ hasError: false, error: undefined })}>
            Try again
          </button>
        </div>
      )
    }

    return this.props.children
  }
}

// ============================================================================
// Performance Monitoring Hook
// ============================================================================

export const usePerformanceMonitoring = (name: string) => {
  React.useEffect(() => {
    const runtime = Runtime.defaultRuntime
    
    // Record initial performance metrics
    const performanceEffect = Effect.gen(function* () {
      const tracingService = yield* TracingService
      const context = yield* tracingService.getCurrentContext()
      
      if (context) {
        // Record browser metrics
        yield* tracingService.recordBrowserMetrics(context)
        
        // Record performance metrics
        if (typeof window !== 'undefined' && 'performance' in window) {
          const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming
          
          if (navigation) {
            yield* tracingService.recordPerformanceMetric(context, {
              name: 'page_load_time',
              value: navigation.loadEventEnd - navigation.navigationStart,
              unit: 'ms',
              metadata: { component: name },
            })
            
            yield* tracingService.recordPerformanceMetric(context, {
              name: 'dom_content_loaded',
              value: navigation.domContentLoadedEventEnd - navigation.navigationStart,
              unit: 'ms',
              metadata: { component: name },
            })
          }
        }
      }
    })

    Runtime.runPromise(runtime)(performanceEffect).catch(console.error)
  }, [name])
}

// ============================================================================
// User Interaction Tracking Hook
// ============================================================================

export const useUserInteractionTracking = () => {
  React.useEffect(() => {
    const runtime = Runtime.defaultRuntime
    
    const handleClick = (event: MouseEvent) => {
      const target = event.target as HTMLElement
      const selector = target.tagName.toLowerCase() + (target.id ? `#${target.id}` : '') + 
                     (target.className ? `.${target.className.split(' ').join('.')}` : '')
      
      const interactionEffect = Effect.gen(function* () {
        const tracingService = yield* TracingService
        const context = yield* tracingService.getCurrentContext()
        
        if (context) {
          yield* tracingService.recordUserInteraction(context, {
            eventType: 'click',
            target: selector,
            metadata: {
              x: event.clientX,
              y: event.clientY,
              timestamp: Date.now(),
            },
          })
        }
      })

      Runtime.runPromise(runtime)(interactionEffect).catch(console.error)
    }

    const handleSubmit = (event: SubmitEvent) => {
      const form = event.target as HTMLFormElement
      const formData = new FormData(form)
      const formFields = Object.fromEntries(formData.entries())
      
      const interactionEffect = Effect.gen(function* () {
        const tracingService = yield* TracingService
        const context = yield* tracingService.getCurrentContext()
        
        if (context) {
          yield* tracingService.recordUserInteraction(context, {
            eventType: 'submit',
            target: form.id || form.className || 'form',
            metadata: {
              fieldCount: Object.keys(formFields).length,
              timestamp: Date.now(),
            },
          })
        }
      })

      Runtime.runPromise(runtime)(interactionEffect).catch(console.error)
    }

    document.addEventListener('click', handleClick, true)
    document.addEventListener('submit', handleSubmit, true)

    return () => {
      document.removeEventListener('click', handleClick, true)
      document.removeEventListener('submit', handleSubmit, true)
    }
  }, [])
}

// ============================================================================
// Main Layer for Integration
// ============================================================================

export const TracingLayer = Layer.mergeAll(TracingServiceLive)

// ============================================================================
// Exports
// ============================================================================

export {
  TraceContextSchema,
  ComponentEventSchema,
  UserInteractionSchema,
  ApiCallTraceSchema,
  PerformanceMetricSchema,
  BrowserMetricSchema,
  TracingError,
}
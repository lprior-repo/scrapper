import React from 'react'
import { createRoot } from 'react-dom/client'
import { Effect, Runtime } from 'effect'
import { App } from './App'
import { TracingLayer } from './tracing'
import './index.css'

const container = document.getElementById('root')
const root = !container
  ? Effect.runSync(Effect.fail(new Error('Root element not found')))
  : createRoot(container)

// Initialize tracing runtime
const tracingRuntime = Runtime.defaultRuntime.pipe(
  Runtime.provide(TracingLayer)
)

// Start global tracing context
Runtime.runPromise(tracingRuntime)(
  Effect.gen(function* () {
    const { TracingService } = yield* import('./tracing')
    const tracingService = yield* TracingService
    const appTrace = yield* tracingService.startTrace('app-init')
    yield* tracingService.setCurrentContext(appTrace)
  })
).catch(console.error)

root.render(<App />)

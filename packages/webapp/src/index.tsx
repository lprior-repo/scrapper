import React from 'react'
import { createRoot } from 'react-dom/client'
import { Effect } from 'effect'
import { App } from './App'

const container = document.getElementById('root')
const root = !container
  ? Effect.runSync(Effect.fail(new Error('Root element not found')))
  : createRoot(container)
root.render(<App />)

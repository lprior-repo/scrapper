# AI-PATTERNS.md

This file contains code patterns and examples for AI assistants to follow.
These patterns ensure consistency across the codebase.

## Component Patterns

### Basic Component Template

```jsx
import React from 'react';

interface ComponentNameProps {
  // Define props here
}

export function ComponentName({ ...props }: ComponentNameProps) {
  return (
    <div>
      {/* Component content */}
    </div>
  );
}
```

### Component with State

```jsx
import React, { useState } from 'react'

export function StatefulComponent() {
  const [state, setState] = useState(initialValue)

  return <div>{/* Component content */}</div>
}
```

### Component with Effects

```jsx
import React, { useEffect } from 'react'

export function EffectComponent() {
  useEffect(() => {
    // Effect logic
    return () => {
      // Cleanup
    }
  }, [dependencies])

  return <div />
}
```

## Hook Patterns

### Custom Hook Template

```js
import { useState, useEffect } from 'react'

export function useCustomHook(param) {
  const [data, setData] = useState(null)

  useEffect(() => {
    // Hook logic
  }, [param])

  return { data }
}
```

## API Patterns

### API Call Pattern

```js
export async function fetchData(endpoint: string) {
  try {
    const response = await fetch(`/api/${endpoint}`);
    if (!response.ok) {
      throw new Error('API request failed');
    }
    return await response.json();
  } catch (error) {
    console.error('API Error:', error);
    throw error;
  }
}
```

### API Hook Pattern

```js
export function useApiData(endpoint: string) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchData(endpoint)
      .then(setData)
      .catch(setError)
      .finally(() => setLoading(false));
  }, [endpoint]);

  return { data, loading, error };
}
```

## State Management Patterns

### Context Pattern

```jsx
import React, { createContext, useContext, useState } from 'react'

const StateContext = createContext()

export function StateProvider({ children }) {
  const [state, setState] = useState(initialState)

  return (
    <StateContext.Provider value={{ state, setState }}>
      {children}
    </StateContext.Provider>
  )
}

export function useAppState() {
  const context = useContext(StateContext)
  if (!context) {
    throw new Error('useAppState must be used within StateProvider')
  }
  return context
}
```

## Error Handling Patterns

### Try-Catch Pattern

```js
export async function safeOperation() {
  try {
    const result = await riskyOperation()
    return { success: true, data: result }
  } catch (error) {
    console.error('Operation failed:', error)
    return { success: false, error: error.message }
  }
}
```

### Error Boundary Pattern

```jsx
import React, { Component, ErrorInfo } from 'react';

export class ErrorBoundary extends Component {
  state = { hasError: false };

  static getDerivedStateFromError(error: Error) {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Error caught:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return <div>Something went wrong.</div>;
    }
    return this.props.children;
  }
}
```

## Testing Patterns

### Component Test Pattern

```js
import { render, screen } from '@testing-library/react'
import { ComponentName } from './ComponentName'

describe('ComponentName', () => {
  it('renders correctly', () => {
    render(<ComponentName />)
    expect(screen.getByText('expected text')).toBeInTheDocument()
  })

  it('handles user interaction', async () => {
    render(<ComponentName />)
    const button = screen.getByRole('button')
    await userEvent.click(button)
    expect(screen.getByText('updated text')).toBeInTheDocument()
  })
})
```

### Hook Test Pattern

```js
import { renderHook } from '@testing-library/react'
import { useCustomHook } from './useCustomHook'

describe('useCustomHook', () => {
  it('returns expected data', () => {
    const { result } = renderHook(() => useCustomHook('param'))
    expect(result.current.data).toBe(expectedValue)
  })
})
```

## Style Patterns

### CSS Module Pattern

```css
/* ComponentName.module.css */
.container {
  display: flex;
  align-items: center;
  padding: 1rem;
}

.title {
  font-size: 1.5rem;
  font-weight: bold;
}
```

### Styled Component Pattern

```js
import styled from 'styled-components'

export const Container = styled.div`
  display: flex;
  align-items: center;
  padding: 1rem;
`

export const Title = styled.h1`
  font-size: 1.5rem;
  font-weight: bold;
`
```

## Type Patterns (TypeScript)

### Interface Pattern

```typescript
export interface User {
  id: string
  name: string
  email: string
  role: 'admin' | 'user'
}

export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
}
```

### Type Guard Pattern

```typescript
export function isUser(obj: any): obj is User {
  return (
    typeof obj === 'object' &&
    typeof obj.id === 'string' &&
    typeof obj.name === 'string' &&
    typeof obj.email === 'string'
  )
}
```

## File Organization Patterns

### Component File Structure

```
ComponentName/
├── ComponentName.tsx      # Main component
├── ComponentName.test.tsx # Tests
├── ComponentName.module.css # Styles
├── types.ts              # Local types
└── index.ts              # Export
```

### Feature File Structure

```
features/
└── user/
    ├── components/       # Feature components
    ├── hooks/           # Feature hooks
    ├── utils/           # Feature utilities
    ├── types.ts         # Feature types
    └── index.ts         # Feature exports
```

---

_This file is part of the AI Context system. Update it when patterns change._

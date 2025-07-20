import React from 'react'

interface IGraphErrorDisplayProps {
  readonly error: unknown
}

export const GraphErrorDisplay: React.FC<IGraphErrorDisplayProps> = ({
  error,
}) => (
  <div
    style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
      color: '#f85149',
      padding: '2rem',
      textAlign: 'center',
    }}
  >
    <div>
      <h2>Error loading graph</h2>
      <p>{error instanceof Error ? error.message : String(error)}</p>
    </div>
  </div>
)

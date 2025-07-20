import React from 'react'

export const GraphLoadingSpinner: React.FC = () => (
  <div
    style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh',
      fontSize: '1.5rem',
      color: '#58a6ff',
    }}
  >
    <div>Loading graph data...</div>
  </div>
)

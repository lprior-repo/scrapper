import React, { useState } from 'react'
import { GraphCanvas } from './components/GraphCanvas'

const AppHeader: React.FC<{
  organization: string
  onOrganizationChange: (org: string) => void
  onScan: () => void
  useTopics: boolean
  onUseTopicsChange: (useTopics: boolean) => void
}> = ({
  organization,
  onOrganizationChange,
  onScan,
  useTopics,
  onUseTopicsChange,
}) => (
  <header
    style={{
      position: 'absolute',
      top: 0,
      left: 0,
      right: 0,
      zIndex: 1000,
      backgroundColor: 'rgba(13, 17, 23, 0.95)',
      borderBottom: '1px solid #30363d',
      padding: '1rem 2rem',
      display: 'flex',
      alignItems: 'center',
      gap: '1rem',
      backdropFilter: 'blur(10px)',
    }}
  >
    <h1
      style={{
        margin: 0,
        fontSize: '1.5rem',
        fontWeight: 600,
        color: '#f0f6fc',
      }}
    >
      GitHub Codeowners Visualization
    </h1>

    <div
      style={{
        marginLeft: 'auto',
        display: 'flex',
        gap: '0.5rem',
        alignItems: 'center',
      }}
    >
      <input
        type="text"
        value={organization}
        onChange={(e) => onOrganizationChange(e.target.value)}
        placeholder="Enter organization name"
        style={{
          padding: '0.5rem 1rem',
          borderRadius: '6px',
          border: '1px solid #30363d',
          backgroundColor: '#0d1117',
          color: '#c9d1d9',
          fontSize: '14px',
          width: '200px',
        }}
      />

      <label
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: '0.5rem',
          color: '#c9d1d9',
          fontSize: '14px',
          cursor: 'pointer',
          userSelect: 'none',
        }}
      >
        <input
          type="checkbox"
          checked={useTopics}
          onChange={(e) => onUseTopicsChange(e.target.checked)}
          style={{
            width: '16px',
            height: '16px',
            accentColor: '#238636',
            cursor: 'pointer',
          }}
        />
        Use Topics instead of Teams
      </label>

      <button
        onClick={onScan}
        disabled={!organization}
        style={{
          padding: '0.5rem 1.5rem',
          borderRadius: '6px',
          border: 'none',
          backgroundColor: organization ? '#238636' : '#21262d',
          color: organization ? '#ffffff' : '#8b949e',
          fontSize: '14px',
          fontWeight: 500,
          cursor: organization ? 'pointer' : 'not-allowed',
          transition: 'all 0.2s',
        }}
        onMouseEnter={(e) => {
          if (organization) {
            e.currentTarget.style.backgroundColor = '#2ea043'
          }
        }}
        onMouseLeave={(e) => {
          if (organization) {
            e.currentTarget.style.backgroundColor = '#238636'
          }
        }}
      >
        Load Graph
      </button>
    </div>
  </header>
)

export const App: React.FC = () => {
  const [organization, setOrganization] = useState('')
  const [displayedOrg, setDisplayedOrg] = useState<string | null>(null)
  const [useTopics, setUseTopics] = useState(false)
  const [displayedUseTopics, setDisplayedUseTopics] = useState(false)

  const handleScan = () => {
    if (organization) {
      setDisplayedOrg(organization)
      setDisplayedUseTopics(useTopics)
    }
  }

  return (
    <div style={{ width: '100%', height: '100vh', position: 'relative' }}>
      <AppHeader
        organization={organization}
        onOrganizationChange={setOrganization}
        onScan={handleScan}
        useTopics={useTopics}
        onUseTopicsChange={setUseTopics}
      />

      {displayedOrg ? (
        <GraphCanvas
          organization={displayedOrg}
          useTopics={displayedUseTopics}
        />
      ) : (
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            height: '100vh',
            color: '#8b949e',
            fontSize: '1.2rem',
          }}
        >
          Enter an organization name and click &quot;Load Graph&quot; to
          visualize the codeowners
        </div>
      )}
    </div>
  )
}

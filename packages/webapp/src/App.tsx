import React, { useState } from 'react'
import { GraphCanvas } from './components/GraphCanvas'

interface IAppHeaderProps {
  readonly organization: string
  readonly onOrganizationChange: (org: string) => void
  readonly onScan: () => void
  readonly useTopics: boolean
  readonly onUseTopicsChange: (useTopics: boolean) => void
}

const OrganizationInput: React.FC<{
  readonly value: string
  readonly onChange: (value: string) => void
}> = ({ value, onChange }) => (
  <input
    type="text"
    value={value}
    onChange={(e) => onChange(e.target.value)}
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
)

const TopicsToggle: React.FC<{
  readonly checked: boolean
  readonly onChange: (checked: boolean) => void
}> = ({ checked, onChange }) => (
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
      checked={checked}
      onChange={(e) => onChange(e.target.checked)}
      style={{
        width: '16px',
        height: '16px',
        accentColor: '#238636',
        cursor: 'pointer',
      }}
    />
    Use Topics instead of Teams
  </label>
)

const LoadButton: React.FC<{
  readonly onClick: () => void
  readonly organization: string
}> = ({ onClick, organization }) => {
  const [isHovered, setIsHovered] = useState(false)
  
  const getBackgroundColor = () => {
    if (!organization) return '#21262d'
    return isHovered ? '#2ea043' : '#238636'
  }

  return (
    <button
      onClick={onClick}
      disabled={!organization}
      style={{
        padding: '0.5rem 1.5rem',
        borderRadius: '6px',
        border: 'none',
        backgroundColor: getBackgroundColor(),
        color: organization ? '#ffffff' : '#8b949e',
        fontSize: '14px',
        fontWeight: 500,
        cursor: organization ? 'pointer' : 'not-allowed',
        transition: 'all 0.2s',
      }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      Load Graph
    </button>
  )
}

const AppHeader: React.FC<IAppHeaderProps> = ({
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
      <OrganizationInput value={organization} onChange={onOrganizationChange} />
      <TopicsToggle checked={useTopics} onChange={onUseTopicsChange} />
      <LoadButton onClick={onScan} organization={organization} />
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

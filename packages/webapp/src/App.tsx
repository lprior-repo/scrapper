import React, { useState } from 'react'
import { AppHeader, GraphCanvas } from './components'

export const App: React.FC = () => {
  const [organization, setOrganization] = useState('')
  const [displayedOrg, setDisplayedOrg] = useState<string | null>(null)
  const [useTopics, setUseTopics] = useState(false)
  const [displayedUseTopics, setDisplayedUseTopics] = useState(false)

  const handleScan = () =>
    organization
      ? (setDisplayedOrg(organization),
        setDisplayedUseTopics(useTopics),
        undefined)
      : undefined

  return (
    <div className="w-full h-screen relative">
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
        <div className="flex justify-center items-center h-screen text-[#8b949e] text-xl">
          Enter an organization name and click &quot;Load Graph&quot; to
          visualize the codeowners
        </div>
      )}
    </div>
  )
}

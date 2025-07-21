import React from 'react'
import { OrganizationInput } from './OrganizationInput'
import { TopicsToggle } from './TopicsToggle'
import { LoadButton } from './LoadButton'

interface AppHeaderProps {
  readonly organization: string
  readonly onOrganizationChange: (org: string) => void
  readonly onScan: () => void
  readonly useTopics: boolean
  readonly onUseTopicsChange: (useTopics: boolean) => void
}

export const AppHeader: React.FC<AppHeaderProps> = ({
  organization,
  onOrganizationChange,
  onScan,
  useTopics,
  onUseTopicsChange,
}) => (
  <header className="absolute top-0 left-0 right-0 z-[1000] bg-[rgba(13,17,23,0.95)] border-b border-[#30363d] p-4 px-8 flex items-center gap-4 backdrop-blur-[10px]">
    <h1 className="m-0 text-2xl font-semibold text-[#f0f6fc]">
      GitHub Codeowners Visualization
    </h1>

    <div className="ml-auto flex gap-2 items-center">
      <OrganizationInput value={organization} onChange={onOrganizationChange} />
      <TopicsToggle checked={useTopics} onChange={onUseTopicsChange} />
      <LoadButton onClick={onScan} organization={organization} />
    </div>
  </header>
)

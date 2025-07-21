import React, { useState } from 'react'

interface LoadButtonProps {
  readonly onClick: () => void
  readonly organization: string
}

export const LoadButton: React.FC<LoadButtonProps> = ({
  onClick,
  organization,
}) => {
  const [isHovered, setIsHovered] = useState(false)

  return (
    <button
      onClick={onClick}
      disabled={!organization}
      className={`py-2 px-6 rounded-md border-none text-sm font-medium transition-all duration-200 ${
        !organization
          ? 'bg-[#21262d] text-[#8b949e] cursor-not-allowed'
          : isHovered
            ? 'bg-[#2ea043] text-white cursor-pointer'
            : 'bg-[#238636] text-white cursor-pointer'
      }`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      Load Graph
    </button>
  )
}

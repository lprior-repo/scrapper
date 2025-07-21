import React from 'react'

interface OrganizationInputProps {
  readonly value: string
  readonly onChange: (value: string) => void
}

export const OrganizationInput: React.FC<OrganizationInputProps> = ({
  value,
  onChange,
}) => (
  <input
    type="text"
    value={value}
    onChange={(e) => onChange(e.target.value)}
    placeholder="Enter organization name"
    className="px-4 py-2 rounded-md border border-[#30363d] bg-[#0d1117] text-[#c9d1d9] text-sm w-50"
  />
)

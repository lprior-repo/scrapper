import React from 'react'

interface TopicsToggleProps {
  readonly checked: boolean
  readonly onChange: (checked: boolean) => void
}

export const TopicsToggle: React.FC<TopicsToggleProps> = ({
  checked,
  onChange,
}) => (
  <label className="flex items-center gap-2 text-[#c9d1d9] text-sm cursor-pointer select-none">
    <input
      type="checkbox"
      checked={checked}
      onChange={(e) => onChange(e.target.checked)}
      className="w-4 h-4 cursor-pointer"
      style={{ accentColor: '#238636' }}
    />
    Use Topics instead of Teams
  </label>
)

import React, { useEffect, useRef } from 'react'
import type { ContextMenuConfig } from './types/context-menu'

interface ContextMenuProps {
  readonly config: ContextMenuConfig | null
  readonly onClose: () => void
}

/**
 * Context menu container styles matching GitHub dark theme
 */
const menuContainerStyle = (position: {
  readonly x: number
  readonly y: number
}) => ({
  position: 'fixed' as const,
  left: `${position.x}px`,
  top: `${position.y}px`,
  background: '#1c2128',
  border: '1px solid #30363d',
  borderRadius: '6px',
  boxShadow: '0 8px 24px rgba(0, 0, 0, 0.12)',
  zIndex: 10000,
  minWidth: '180px',
  padding: '4px 0',
  fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  fontSize: '14px',
  color: '#f0f6fc',
  userSelect: 'none' as const,
  animation: 'contextMenuSlideIn 0.15s ease-out',
})

/**
 * Menu item styles
 */
const menuItemStyle = (disabled: boolean = false) => ({
  padding: '8px 16px',
  cursor: disabled ? 'default' : 'pointer',
  display: 'flex',
  alignItems: 'center',
  gap: '8px',
  color: disabled ? '#7d8590' : '#f0f6fc',
  backgroundColor: 'transparent',
  border: 'none',
  width: '100%',
  textAlign: 'left' as const,
  fontSize: '14px',
  fontFamily: 'inherit',
  transition: 'background-color 0.1s ease',
})

/**
 * Menu item hover styles
 */
const menuItemHoverStyle = {
  backgroundColor: '#21262d',
}

/**
 * Separator line styles
 */
const separatorStyle = {
  height: '1px',
  backgroundColor: '#30363d',
  margin: '4px 0',
}

/**
 * CSS animation for context menu appearance
 */
const menuAnimationCSS = `
  @keyframes contextMenuSlideIn {
    from {
      opacity: 0;
      transform: scale(0.95) translateY(-5px);
    }
    to {
      opacity: 1;
      transform: scale(1) translateY(0);
    }
  }
`

/**
 * Injects CSS animation styles into document head
 */
const injectMenuAnimations = (): void => {
  const existingStyle = document.getElementById('context-menu-animations')
  if (existingStyle) return

  const style = document.createElement('style')
  style.id = 'context-menu-animations'
  style.textContent = menuAnimationCSS
  document.head.appendChild(style)
}

/**
 * Adjusts menu position to stay within viewport bounds
 */
const adjustMenuPosition = (
  position: { readonly x: number; readonly y: number },
  menuRef: React.RefObject<HTMLDivElement>
): { readonly x: number; readonly y: number } => {
  if (!menuRef.current) return position

  const rect = menuRef.current.getBoundingClientRect()
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight

  const adjustedX =
    position.x + rect.width > viewportWidth
      ? viewportWidth - rect.width - 10
      : position.x

  const adjustedY =
    position.y + rect.height > viewportHeight
      ? viewportHeight - rect.height - 10
      : position.y

  return { x: adjustedX, y: adjustedY }
}

/**
 * Menu item component with hover effects
 */
const MenuItem: React.FC<{
  readonly item: {
    readonly id: string
    readonly label: string
    readonly icon?: string
    readonly disabled?: boolean
    readonly action: () => void
  }
  readonly onItemClick: (action: () => void) => void
}> = ({ item, onItemClick }) => {
  const [isHovered, setIsHovered] = React.useState(false)

  const handleClick = (): void => {
    if (!item.disabled) {
      onItemClick(item.action)
    }
  }

  const combinedStyle = {
    ...menuItemStyle(item.disabled),
    ...(isHovered && !item.disabled ? menuItemHoverStyle : {}),
  }

  return (
    <button
      type="button"
      style={combinedStyle}
      onClick={handleClick}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      disabled={item.disabled}
    >
      {item.icon && <span style={{ fontSize: '16px' }}>{item.icon}</span>}
      <span>{item.label}</span>
    </button>
  )
}

/**
 * Separator component
 */
const Separator: React.FC = () => <div style={separatorStyle} />

/**
 * Handles clicks outside the menu to close it
 */
const useClickOutside = (
  ref: React.RefObject<HTMLDivElement>,
  onClose: () => void
): void => {
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent): void => {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        onClose()
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [ref, onClose])
}

/**
 * Context menu component with proper positioning and theme styling
 */
export const ContextMenu: React.FC<ContextMenuProps> = ({
  config,
  onClose,
}) => {
  const menuRef = useRef<HTMLDivElement>(null)

  // Close menu on outside clicks
  useClickOutside(menuRef, onClose)

  // Inject animations on mount
  useEffect(() => {
    injectMenuAnimations()
  }, [])

  // Close menu on Escape key
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent): void => {
      if (event.key === 'Escape') {
        onClose()
      }
    }

    document.addEventListener('keydown', handleEscape)
    return () => document.removeEventListener('keydown', handleEscape)
  }, [onClose])

  // Adjust position to stay within viewport
  const [adjustedPosition, setAdjustedPosition] = React.useState(
    config?.position ?? { x: 0, y: 0 }
  )

  useEffect(() => {
    if (config && menuRef.current) {
      const adjusted = adjustMenuPosition(config.position, menuRef)
      setAdjustedPosition(adjusted)
    }
  }, [config])

  if (!config) return null

  const handleItemClick = (action: () => void): void => {
    action()
    onClose()
  }

  return (
    <div
      ref={menuRef}
      style={menuContainerStyle(adjustedPosition)}
      data-testid={`context-menu-${config.type}`}
    >
      {config.items.map((item, index) =>
        item.separator ? (
          <Separator key={`separator-${index}`} />
        ) : (
          <MenuItem key={item.id} item={item} onItemClick={handleItemClick} />
        )
      )}
    </div>
  )
}

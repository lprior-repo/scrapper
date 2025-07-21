/**
 * Context menu types and interfaces for graph interactions
 */

/**
 * Position coordinates for context menu placement
 */
export interface MenuPosition {
  readonly x: number
  readonly y: number
}

/**
 * Context menu item configuration
 */
export interface ContextMenuItem {
  readonly id: string
  readonly label: string
  readonly icon?: string
  readonly disabled?: boolean
  readonly separator?: boolean
  readonly action: () => void
}

/**
 * Context menu configuration based on target type
 */
export interface ContextMenuConfig {
  readonly type: 'node' | 'edge' | 'background'
  readonly targetId?: string
  readonly position: MenuPosition
  readonly items: readonly ContextMenuItem[]
}

/**
 * Context menu state for component management
 */
export interface ContextMenuState {
  readonly isVisible: boolean
  readonly config: ContextMenuConfig | null
}

/**
 * Element data for context menu actions
 */
export interface ElementData {
  readonly id: string
  readonly label: string
  readonly type: string
  readonly data?: Record<string, unknown>
}

/**
 * Hidden elements tracking state
 */
export interface HiddenElementsState {
  readonly nodes: readonly string[]
  readonly edges: readonly string[]
}

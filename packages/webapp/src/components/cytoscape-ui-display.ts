/**
 * Edge information interface
 */
export interface EdgeInfo {
  readonly id: string
  readonly relationship: string
  readonly label: string
  readonly source: {
    readonly id: string
    readonly label: string
    readonly type: string
  }
  readonly target: {
    readonly id: string
    readonly label: string
    readonly type: string
  }
}

/**
 * Node information interface
 */
export interface NodeInfo {
  readonly id: string
  readonly type: string
  readonly label: string
  readonly data?: Record<string, unknown>
}

/**
 * Formats relationship type with icons
 */
const formatRelationshipType = (type: string): string => {
  const typeMap: Record<string, string> = {
    owns: '👑 Owns',
    member_of: '👥 Member Of',
    codeowner: '🛡️ Code Owner',
    maintained_by: '🔧 Maintained By',
    has_topic: '🏷️ Has Topic',
    has: '📁 Has',
  }
  return typeMap[type] || `🔗 ${type}`
}

/**
 * Formats node type with icons
 */
const formatNodeType = (type: string): string => {
  const typeMap: Record<string, string> = {
    organization: '🏢 Organization',
    repository: '📁 Repository',
    team: '👥 Team',
    user: '👤 User',
    topic: '🏷️ Topic',
  }
  return typeMap[type] || `🔵 ${type}`
}

/**
 * Formats data value for display
 */
const formatDataValue = (value: unknown): string =>
  typeof value === 'object' ? JSON.stringify(value, null, 2) : String(value)

/**
 * Creates HTML for data row
 */
const createDataRow = ([key, value]: readonly [string, unknown]): string => `
  <div style="margin-bottom: 6px; padding: 6px; background: #21262d; border-radius: 3px;">
    <strong style="color: #79c0ff;">${key}:</strong>
    <div style="margin-top: 2px; word-break: break-all;">${formatDataValue(value)}</div>
  </div>
`

/**
 * Creates HTML for data table
 */
const createDataTable = (data: Record<string, unknown>): string =>
  !data || Object.keys(data).length === 0
    ? '<div style="color: #7d8590; font-style: italic;">No additional data</div>'
    : Object.entries(data).map(createDataRow).join('')

/**
 * Creates removal function for display element
 */
const createRemovalFunction =
  (element: HTMLElement, style: HTMLElement) => (): void => {
    element.remove()
    style.remove()
  }

/**
 * Adds slideIn animation CSS
 */
const addSlideInAnimation = (): HTMLElement => {
  const style = document.createElement('style')
  style.textContent = `
    @keyframes slideIn {
      from {
        transform: translateX(100%);
        opacity: 0;
      }
      to {
        transform: translateX(0);
        opacity: 1;
      }
    }
  `
  document.head.appendChild(style)
  return style
}

/**
 * Adds slideInLeft animation CSS
 */
const addSlideInLeftAnimation = (): HTMLElement => {
  const style = document.createElement('style')
  style.textContent = `
    @keyframes slideInLeft {
      from {
        transform: translateX(-100%);
        opacity: 0;
      }
      to {
        transform: translateX(0);
        opacity: 1;
      }
    }
  `
  document.head.appendChild(style)
  return style
}

/**
 * Creates edge info HTML content
 */
const createEdgeInfoContent = (edgeInfo: EdgeInfo): string => `
  <div style="margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid #30363d;">
    <strong style="color: #58a6ff;">Relationship Details</strong>
  </div>
  <div style="margin-bottom: 8px;">
    <strong>Type:</strong> ${formatRelationshipType(edgeInfo.relationship)}
  </div>
  ${
    edgeInfo.label !== 'No label'
      ? `
  <div style="margin-bottom: 8px;">
    <strong>Label:</strong> ${edgeInfo.label}
  </div>
  `
      : ''
  }
  <div style="margin-bottom: 8px;">
    <strong>Source:</strong> ${edgeInfo.source.label} (${edgeInfo.source.type})
  </div>
  <div style="margin-bottom: 12px;">
    <strong>Target:</strong> ${edgeInfo.target.label} (${edgeInfo.target.type})
  </div>
  <div style="font-size: 12px; color: #7d8590;">
    Click anywhere to dismiss
  </div>
`

/**
 * Creates node info HTML content
 */
const createNodeInfoContent = (nodeInfo: NodeInfo): string => `
  <div style="margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid #30363d;">
    <strong style="color: #58a6ff;">Node Details</strong>
  </div>
  <div style="margin-bottom: 8px;">
    <strong>Type:</strong> ${formatNodeType(nodeInfo.type)}
  </div>
  <div style="margin-bottom: 8px;">
    <strong>ID:</strong> <code style="background: #21262d; padding: 2px 4px; border-radius: 3px;">${nodeInfo.id}</code>
  </div>
  <div style="margin-bottom: 12px;">
    <strong>Label:</strong> ${nodeInfo.label}
  </div>
  ${
    Object.keys(nodeInfo.data || {}).length > 0
      ? `
  <div style="margin-bottom: 8px;">
    <strong>Additional Data:</strong>
  </div>
  <div style="margin-bottom: 12px;">
    ${createDataTable(nodeInfo.data)}
  </div>
  `
      : ''
  }
  <div style="font-size: 12px; color: #7d8590;">
    Click anywhere to dismiss
  </div>
`

/**
 * Creates a temporary UI element to display edge information
 */
export const createEdgeInfoDisplay = (edgeInfo: EdgeInfo): void => {
  // Remove existing edge info display
  const existingDisplay = document.getElementById('edge-info-display')
  existingDisplay?.remove()

  const infoElement = document.createElement('div')
  infoElement.id = 'edge-info-display'
  infoElement.style.cssText = `
    position: fixed;
    top: 20px;
    right: 20px;
    background: #1c2128;
    border: 1px solid #30363d;
    border-radius: 6px;
    padding: 16px;
    color: #f0f6fc;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    font-size: 14px;
    max-width: 350px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    z-index: 1000;
    animation: slideIn 0.3s ease-out;
  `

  infoElement.innerHTML = createEdgeInfoContent(edgeInfo)

  // Add CSS animation
  const style = addSlideInAnimation()
  document.body.appendChild(infoElement)

  // Auto-remove after 10 seconds or on click
  const removeDisplay = createRemovalFunction(infoElement, style)
  setTimeout(removeDisplay, 10000)
  infoElement.addEventListener('click', removeDisplay)
  document.addEventListener('click', removeDisplay, { once: true })
}

/**
 * Creates a temporary UI element to display node information
 */
export const createNodeInfoDisplay = (nodeInfo: NodeInfo): void => {
  // Remove existing node info display
  const existingDisplay = document.getElementById('node-info-display')
  existingDisplay?.remove()

  const infoElement = document.createElement('div')
  infoElement.id = 'node-info-display'
  infoElement.style.cssText = `
    position: fixed;
    top: 20px;
    left: 20px;
    background: #1c2128;
    border: 1px solid #30363d;
    border-radius: 6px;
    padding: 16px;
    color: #f0f6fc;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    font-size: 14px;
    max-width: 400px;
    max-height: 80vh;
    overflow-y: auto;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    z-index: 1000;
    animation: slideInLeft 0.3s ease-out;
  `

  infoElement.innerHTML = createNodeInfoContent(nodeInfo)

  // Add CSS animation for left slide-in
  const style = addSlideInLeftAnimation()
  document.body.appendChild(infoElement)

  // Auto-remove after 10 seconds or on click
  const removeDisplay = createRemovalFunction(infoElement, style)
  setTimeout(removeDisplay, 10000)
  infoElement.addEventListener('click', removeDisplay)
  document.addEventListener('click', removeDisplay, { once: true })
}

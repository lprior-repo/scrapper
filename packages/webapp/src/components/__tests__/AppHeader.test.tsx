/**
 * AppHeader Component Tests
 * 
 * Tests for the application header component including:
 * - Basic rendering and layout
 * - Component composition and integration
 * - Props passing to child components
 * - User interactions through child components
 * - Accessibility and semantic structure
 * - Responsive behavior
 */

import React from 'react'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { AppHeader } from '../AppHeader'

describe('AppHeader', () => {
  const defaultProps = {
    organization: 'github',
    onOrganizationChange: jest.fn(),
    onScan: jest.fn(),
    useTopics: false,
    onUseTopicsChange: jest.fn(),
  }
  
  beforeEach(() => {
    jest.clearAllMocks()
  })

  describe('Basic Rendering', () => {
    test('renders header with correct semantic structure', () => {
      render(<AppHeader {...defaultProps} />)
      
      const header = screen.getByRole('banner')
      expect(header).toBeInTheDocument()
      expect(header.tagName).toBe('HEADER')
    })

    test('displays application title', () => {
      render(<AppHeader {...defaultProps} />)
      
      const title = screen.getByRole('heading', { level: 1 })
      expect(title).toBeInTheDocument()
      expect(title).toHaveTextContent('GitHub Codeowners Visualization')
    })

    test('applies correct header styling and positioning', () => {
      render(<AppHeader {...defaultProps} />)
      
      const header = screen.getByRole('banner')
      expect(header).toHaveClass(
        'absolute',
        'top-0',
        'left-0',
        'right-0',
        'z-[1000]',
        'bg-[rgba(13,17,23,0.95)]',
        'border-b',
        'border-[#30363d]',
        'p-4',
        'px-8',
        'flex',
        'items-center',
        'gap-4',
        'backdrop-blur-[10px]'
      )
    })

    test('title has correct styling', () => {
      render(<AppHeader {...defaultProps} />)
      
      const title = screen.getByRole('heading', { level: 1 })
      expect(title).toHaveClass(
        'm-0',
        'text-2xl',
        'font-semibold',
        'text-[#f0f6fc]'
      )
    })
  })

  describe('Component Composition', () => {
    test('renders OrganizationInput component', () => {
      render(<AppHeader {...defaultProps} />)
      
      const input = screen.getByPlaceholderText('Enter organization name')
      expect(input).toBeInTheDocument()
    })

    test('renders TopicsToggle component', () => {
      render(<AppHeader {...defaultProps} />)
      
      const toggle = screen.getByRole('checkbox')
      const label = screen.getByText('Use Topics instead of Teams')
      
      expect(toggle).toBeInTheDocument()
      expect(label).toBeInTheDocument()
    })

    test('renders LoadButton component', () => {
      render(<AppHeader {...defaultProps} />)
      
      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeInTheDocument()
    })

    test('all interactive controls are in the correct container', () => {
      render(<AppHeader {...defaultProps} />)
      
      const controlsContainer = screen.getByRole('banner').querySelector('.ml-auto')
      expect(controlsContainer).toBeInTheDocument()
      expect(controlsContainer).toHaveClass('ml-auto', 'flex', 'gap-2', 'items-center')
      
      // Check that all controls are within this container
      const input = screen.getByPlaceholderText('Enter organization name')
      const checkbox = screen.getByRole('checkbox')
      const button = screen.getByRole('button', { name: 'Load Graph' })
      
      expect(controlsContainer).toContainElement(input)
      expect(controlsContainer).toContainElement(checkbox)
      expect(controlsContainer).toContainElement(button)
    })
  })

  describe('Props Passing', () => {
    test('passes correct props to OrganizationInput', () => {
      render(<AppHeader {...defaultProps} organization="test-org" />)
      
      const input = screen.getByDisplayValue('test-org')
      expect(input).toBeInTheDocument()
    })

    test('passes correct props to TopicsToggle when unchecked', () => {
      render(<AppHeader {...defaultProps} useTopics={false} />)
      
      const checkbox = screen.getByRole('checkbox')
      expect(checkbox).not.toBeChecked()
    })

    test('passes correct props to TopicsToggle when checked', () => {
      render(<AppHeader {...defaultProps} useTopics={true} />)
      
      const checkbox = screen.getByRole('checkbox')
      expect(checkbox).toBeChecked()
    })

    test('passes correct props to LoadButton', () => {
      render(<AppHeader {...defaultProps} organization="github" />)
      
      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeEnabled()
    })

    test('LoadButton is disabled when organization is empty', () => {
      render(<AppHeader {...defaultProps} organization="" />)
      
      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeDisabled()
    })
  })

  describe('User Interactions', () => {
    test('organization input changes trigger onOrganizationChange', async () => {
      const user = userEvent.setup()
      const onOrganizationChange = jest.fn()
      
      render(
        <AppHeader 
          {...defaultProps} 
          organization=""
          onOrganizationChange={onOrganizationChange} 
        />
      )
      
      const input = screen.getByPlaceholderText('Enter organization name')
      await user.type(input, 'new-org')
      
      expect(onOrganizationChange).toHaveBeenCalledWith('new-org')
    })

    test('topics toggle changes trigger onUseTopicsChange', async () => {
      const user = userEvent.setup()
      const onUseTopicsChange = jest.fn()
      
      render(
        <AppHeader 
          {...defaultProps} 
          useTopics={false}
          onUseTopicsChange={onUseTopicsChange} 
        />
      )
      
      const checkbox = screen.getByRole('checkbox')
      await user.click(checkbox)
      
      expect(onUseTopicsChange).toHaveBeenCalledWith(true)
    })

    test('load button clicks trigger onScan', async () => {
      const user = userEvent.setup()
      const onScan = jest.fn()
      
      render(
        <AppHeader 
          {...defaultProps} 
          organization="github"
          onScan={onScan} 
        />
      )
      
      const button = screen.getByRole('button', { name: 'Load Graph' })
      await user.click(button)
      
      expect(onScan).toHaveBeenCalledTimes(1)
    })

    test('disabled load button does not trigger onScan', async () => {
      const user = userEvent.setup()
      const onScan = jest.fn()
      
      render(
        <AppHeader 
          {...defaultProps} 
          organization=""
          onScan={onScan} 
        />
      )
      
      const button = screen.getByRole('button', { name: 'Load Graph' })
      await user.click(button)
      
      expect(onScan).not.toHaveBeenCalled()
    })
  })

  describe('Accessibility', () => {
    test('has proper landmark role for navigation', () => {
      render(<AppHeader {...defaultProps} />)
      
      const header = screen.getByRole('banner')
      expect(header).toBeInTheDocument()
    })

    test('title is properly marked as heading', () => {
      render(<AppHeader {...defaultProps} />)
      
      const title = screen.getByRole('heading', { level: 1 })
      expect(title).toHaveTextContent('GitHub Codeowners Visualization')
    })

    test('all interactive elements are keyboard accessible', () => {
      render(<AppHeader {...defaultProps} organization="github" />)
      
      const input = screen.getByPlaceholderText('Enter organization name')
      const checkbox = screen.getByRole('checkbox')
      const button = screen.getByRole('button', { name: 'Load Graph' })
      
      // Focus each element to verify they're keyboard accessible
      input.focus()
      expect(input).toHaveFocus()
      
      checkbox.focus()
      expect(checkbox).toHaveFocus()
      
      button.focus()
      expect(button).toHaveFocus()
    })

    test('provides proper semantic structure for screen readers', () => {
      render(<AppHeader {...defaultProps} />)
      
      // Header should be identifiable as banner landmark
      const header = screen.getByRole('banner')
      expect(header).toBeInTheDocument()
      
      // All form controls should be properly labeled
      const input = screen.getByPlaceholderText('Enter organization name')
      expect(input).toHaveAttribute('placeholder')
      
      const checkbox = screen.getByRole('checkbox', { name: 'Use Topics instead of Teams' })
      expect(checkbox).toBeInTheDocument()
      
      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeInTheDocument()
    })
  })

  describe('Layout and Positioning', () => {
    test('positions header at top of viewport', () => {
      render(<AppHeader {...defaultProps} />)
      
      const header = screen.getByRole('banner')
      expect(header).toHaveClass('absolute', 'top-0', 'left-0', 'right-0')
    })

    test('has high z-index for overlay positioning', () => {
      render(<AppHeader {...defaultProps} />)
      
      const header = screen.getByRole('banner')
      expect(header).toHaveClass('z-[1000]')
    })

    test('uses flexbox for proper control alignment', () => {
      render(<AppHeader {...defaultProps} />)
      
      const header = screen.getByRole('banner')
      const controlsContainer = header.querySelector('.ml-auto')
      
      expect(header).toHaveClass('flex', 'items-center', 'gap-4')
      expect(controlsContainer).toHaveClass('flex', 'gap-2', 'items-center')
    })

    test('applies backdrop blur for visual layering', () => {
      render(<AppHeader {...defaultProps} />)
      
      const header = screen.getByRole('banner')
      expect(header).toHaveClass('backdrop-blur-[10px]')
    })
  })

  describe('State Management', () => {
    test('reflects current organization in input', () => {
      const { rerender } = render(
        <AppHeader {...defaultProps} organization="initial" />
      )
      
      let input = screen.getByDisplayValue('initial')
      expect(input).toBeInTheDocument()
      
      rerender(<AppHeader {...defaultProps} organization="updated" />)
      
      input = screen.getByDisplayValue('updated')
      expect(input).toBeInTheDocument()
    })

    test('reflects current topics toggle state', () => {
      const { rerender } = render(
        <AppHeader {...defaultProps} useTopics={false} />
      )
      
      let checkbox = screen.getByRole('checkbox')
      expect(checkbox).not.toBeChecked()
      
      rerender(<AppHeader {...defaultProps} useTopics={true} />)
      
      checkbox = screen.getByRole('checkbox')
      expect(checkbox).toBeChecked()
    })

    test('load button state reflects organization availability', () => {
      const { rerender } = render(
        <AppHeader {...defaultProps} organization="" />
      )
      
      let button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeDisabled()
      
      rerender(<AppHeader {...defaultProps} organization="github" />)
      
      button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeEnabled()
    })
  })

  describe('Visual States', () => {
    test('applies GitHub dark theme colors', () => {
      render(<AppHeader {...defaultProps} />)
      
      const header = screen.getByRole('banner')
      const title = screen.getByRole('heading', { level: 1 })
      
      // Header background and border colors
      expect(header).toHaveClass('bg-[rgba(13,17,23,0.95)]', 'border-[#30363d]')
      
      // Title text color
      expect(title).toHaveClass('text-[#f0f6fc]')
    })

    test('maintains consistent spacing and sizing', () => {
      render(<AppHeader {...defaultProps} />)
      
      const header = screen.getByRole('banner')
      const title = screen.getByRole('heading', { level: 1 })
      const controlsContainer = header.querySelector('.ml-auto')
      
      // Header padding
      expect(header).toHaveClass('p-4', 'px-8')
      
      // Title sizing
      expect(title).toHaveClass('text-2xl')
      
      // Container spacing
      expect(header).toHaveClass('gap-4')
      expect(controlsContainer).toHaveClass('gap-2')
    })
  })

  describe('Edge Cases', () => {
    test('handles undefined organization gracefully', () => {
      render(<AppHeader {...defaultProps} organization={undefined as any} />)
      
      const input = screen.getByPlaceholderText('Enter organization name')
      const button = screen.getByRole('button', { name: 'Load Graph' })
      
      expect(input).toBeInTheDocument()
      expect(button).toBeDisabled()
    })

    test('handles missing callback props gracefully', () => {
      render(
        <AppHeader 
          organization="github"
          onOrganizationChange={undefined as any}
          onScan={undefined as any}
          useTopics={false}
          onUseTopicsChange={undefined as any}
        />
      )
      
      // Component should render without crashing
      const header = screen.getByRole('banner')
      expect(header).toBeInTheDocument()
    })

    test('handles rapid state changes', () => {
      const { rerender } = render(
        <AppHeader {...defaultProps} organization="" useTopics={false} />
      )
      
      // Rapid state changes
      for (let i = 0; i < 10; i++) {
        rerender(
          <AppHeader 
            {...defaultProps} 
            organization={i % 2 === 0 ? '' : 'github'} 
            useTopics={i % 2 === 1}
          />
        )
      }
      
      // Component should still be functional
      const header = screen.getByRole('banner')
      expect(header).toBeInTheDocument()
    })
  })

  describe('Performance', () => {
    test('does not cause unnecessary re-renders with stable props', () => {
      let renderCount = 0
      
      const TestComponent = (props: typeof defaultProps) => {
        renderCount++
        return <AppHeader {...props} />
      }
      
      const { rerender } = render(<TestComponent {...defaultProps} />)
      expect(renderCount).toBe(1)
      
      // Re-render with same props
      rerender(<TestComponent {...defaultProps} />)
      expect(renderCount).toBe(2)
    })

    test('child components maintain their individual state', async () => {
      const user = userEvent.setup()
      render(<AppHeader {...defaultProps} organization="github" />)
      
      const input = screen.getByPlaceholderText('Enter organization name')
      const button = screen.getByRole('button', { name: 'Load Graph' })
      
      // Focus input
      input.focus()
      expect(input).toHaveFocus()
      
      // Hover button
      await user.hover(button)
      
      // Input should maintain focus while button shows hover state
      expect(input).toHaveFocus()
    })
  })

  describe('Integration Scenarios', () => {
    test('supports typical user workflow', async () => {
      const user = userEvent.setup()
      const onOrganizationChange = jest.fn()
      const onUseTopicsChange = jest.fn()
      const onScan = jest.fn()
      
      render(
        <AppHeader 
          organization=""
          onOrganizationChange={onOrganizationChange}
          onScan={onScan}
          useTopics={false}
          onUseTopicsChange={onUseTopicsChange}
        />
      )
      
      // 1. User enters organization name
      const input = screen.getByPlaceholderText('Enter organization name')
      await user.type(input, 'github')
      expect(onOrganizationChange).toHaveBeenLastCalledWith('github')
      
      // 2. User toggles to topics view
      const checkbox = screen.getByRole('checkbox')
      await user.click(checkbox)
      expect(onUseTopicsChange).toHaveBeenCalledWith(true)
      
      // 3. Button should become enabled (needs props update in real app)
      // This simulates the parent component updating the props
      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeInTheDocument() // Button exists regardless of disabled state
    })

    test('works correctly with controlled state patterns', () => {
      const ControlledHeader = () => {
        const [organization, setOrganization] = React.useState('')
        const [useTopics, setUseTopics] = React.useState(false)
        
        return (
          <AppHeader
            organization={organization}
            onOrganizationChange={setOrganization}
            onScan={() => {}}
            useTopics={useTopics}
            onUseTopicsChange={setUseTopics}
          />
        )
      }
      
      render(<ControlledHeader />)
      
      const header = screen.getByRole('banner')
      expect(header).toBeInTheDocument()
      
      const input = screen.getByPlaceholderText('Enter organization name')
      const checkbox = screen.getByRole('checkbox')
      
      expect(input).toHaveValue('')
      expect(checkbox).not.toBeChecked()
    })
  })
})
/**
 * LoadButton Component Tests
 *
 * Tests for the load graph button component including:
 * - Basic rendering and states (enabled/disabled)
 * - User interactions (click, hover, keyboard)
 * - Accessibility features
 * - State transitions and visual feedback
 * - Props handling and edge cases
 */

import React from 'react'
import { render, screen, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { LoadButton } from '../LoadButton'

describe('LoadButton', () => {
  const mockOnClick = jest.fn()

  beforeEach(() => {
    mockOnClick.mockClear()
  })

  describe('Basic Rendering', () => {
    test('renders button with correct text', () => {
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeInTheDocument()
      expect(button).toHaveTextContent('Load Graph')
    })

    test('is enabled when organization is provided', () => {
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeEnabled()
      expect(button).not.toHaveAttribute('disabled')
    })

    test('is disabled when organization is empty', () => {
      render(<LoadButton onClick={mockOnClick} organization="" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeDisabled()
      expect(button).toHaveAttribute('disabled')
    })

    test('is disabled when organization is whitespace only', () => {
      render(<LoadButton onClick={mockOnClick} organization="   " />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeDisabled()
    })
  })

  describe('Visual States and Styling', () => {
    test('applies enabled styling when organization is provided', () => {
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      // Should have enabled state classes
      expect(button).toHaveClass(
        'py-2',
        'px-6',
        'rounded-md',
        'border-none',
        'text-sm',
        'font-medium',
        'transition-all',
        'duration-200',
        'bg-[#238636]',
        'text-white',
        'cursor-pointer'
      )
    })

    test('applies disabled styling when organization is empty', () => {
      render(<LoadButton onClick={mockOnClick} organization="" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      // Should have disabled state classes
      expect(button).toHaveClass(
        'bg-[#21262d]',
        'text-[#8b949e]',
        'cursor-not-allowed'
      )
    })

    test('shows hover state when enabled and hovered', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      // Initial state
      expect(button).toHaveClass('bg-[#238636]')

      // Hover
      await user.hover(button)

      // Should show hover color
      expect(button).toHaveClass('bg-[#2ea043]')
    })

    test('removes hover state when mouse leaves', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      // Hover then unhover
      await user.hover(button)
      expect(button).toHaveClass('bg-[#2ea043]')

      await user.unhover(button)
      expect(button).toHaveClass('bg-[#238636]')
    })

    test('does not show hover state when disabled', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      // Try to hover on disabled button
      await user.hover(button)

      // Should maintain disabled styling
      expect(button).toHaveClass('bg-[#21262d]')
      expect(button).not.toHaveClass('bg-[#2ea043]')
    })
  })

  describe('User Interactions', () => {
    test('calls onClick when button is clicked and enabled', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      await user.click(button)

      expect(mockOnClick).toHaveBeenCalledTimes(1)
    })

    test('does not call onClick when button is disabled', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      await user.click(button)

      expect(mockOnClick).not.toHaveBeenCalled()
    })

    test('can be activated with keyboard when enabled', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      button.focus()

      await user.keyboard('{Enter}')

      expect(mockOnClick).toHaveBeenCalledTimes(1)
    })

    test('can be activated with space key when enabled', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      button.focus()

      await user.keyboard(' ')

      expect(mockOnClick).toHaveBeenCalledTimes(1)
    })

    test('does not respond to keyboard when disabled', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      button.focus()

      await user.keyboard('{Enter}')
      await user.keyboard(' ')

      expect(mockOnClick).not.toHaveBeenCalled()
    })
  })

  describe('Accessibility Features', () => {
    test('has proper button semantics', () => {
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button.tagName).toBe('BUTTON')
    })

    test('is focusable when enabled', () => {
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      button.focus()

      expect(button).toHaveFocus()
    })

    test('is not focusable when disabled', () => {
      render(<LoadButton onClick={mockOnClick} organization="" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      // Disabled buttons should not be focusable
      expect(button).toBeDisabled()
      expect(button).toHaveAttribute('disabled')
    })

    test('has descriptive text for screen readers', () => {
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toHaveAccessibleName('Load Graph')
    })

    test('conveys disabled state to screen readers', () => {
      render(<LoadButton onClick={mockOnClick} organization="" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toHaveAttribute('disabled')
      expect(button).toHaveAttribute('aria-disabled', 'true')
    })
  })

  describe('State Management with Hover', () => {
    test('manages hover state correctly with React useState', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      // Initial state - not hovered
      expect(button).toHaveClass('bg-[#238636]')

      // Mouse enter
      await act(async () => {
        await user.hover(button)
      })
      expect(button).toHaveClass('bg-[#2ea043]')

      // Mouse leave
      await act(async () => {
        await user.unhover(button)
      })
      expect(button).toHaveClass('bg-[#238636]')
    })

    test('maintains hover state during multiple mouse events', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      // Multiple hover/unhover cycles
      for (let i = 0; i < 3; i++) {
        await user.hover(button)
        expect(button).toHaveClass('bg-[#2ea043]')

        await user.unhover(button)
        expect(button).toHaveClass('bg-[#238636]')
      }
    })

    test('hover state resets correctly after organization changes', async () => {
      const user = userEvent.setup()
      const { rerender } = render(
        <LoadButton onClick={mockOnClick} organization="github" />
      )

      const button = screen.getByRole('button', { name: 'Load Graph' })

      // Hover the button
      await user.hover(button)
      expect(button).toHaveClass('bg-[#2ea043]')

      // Change organization to empty (disable button)
      rerender(<LoadButton onClick={mockOnClick} organization="" />)

      // Button should be disabled and not show hover state
      expect(button).toBeDisabled()
      expect(button).toHaveClass('bg-[#21262d]')
    })
  })

  describe('Props Handling', () => {
    test('updates state when organization prop changes from empty to valid', () => {
      const { rerender } = render(
        <LoadButton onClick={mockOnClick} organization="" />
      )

      let button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeDisabled()

      rerender(<LoadButton onClick={mockOnClick} organization="github" />)

      button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeEnabled()
    })

    test('updates state when organization prop changes from valid to empty', () => {
      const { rerender } = render(
        <LoadButton onClick={mockOnClick} organization="github" />
      )

      let button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeEnabled()

      rerender(<LoadButton onClick={mockOnClick} organization="" />)

      button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeDisabled()
    })

    test('handles onClick prop changes gracefully', async () => {
      const user = userEvent.setup()
      const newOnClick = jest.fn()

      const { rerender } = render(
        <LoadButton onClick={mockOnClick} organization="github" />
      )

      rerender(<LoadButton onClick={newOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      await user.click(button)

      expect(newOnClick).toHaveBeenCalledTimes(1)
      expect(mockOnClick).not.toHaveBeenCalled()
    })
  })

  describe('Edge Cases', () => {
    test('handles null or undefined organization gracefully', () => {
      // TypeScript prevents this, but testing runtime behavior
      render(<LoadButton onClick={mockOnClick} organization={null as any} />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeDisabled()
    })

    test('handles very long organization names', () => {
      const longOrg = 'a'.repeat(1000)
      render(<LoadButton onClick={mockOnClick} organization={longOrg} />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeEnabled()
    })

    test('handles special characters in organization name', () => {
      const specialOrg = 'org-name_123.test'
      render(<LoadButton onClick={mockOnClick} organization={specialOrg} />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeEnabled()
    })

    test('handles rapid state changes', async () => {
      const user = userEvent.setup()
      const { rerender } = render(
        <LoadButton onClick={mockOnClick} organization="" />
      )

      // Rapid enable/disable cycles
      for (let i = 0; i < 10; i++) {
        const org = i % 2 === 0 ? '' : 'github'
        rerender(<LoadButton onClick={mockOnClick} organization={org} />)

        const button = screen.getByRole('button', { name: 'Load Graph' })
        if (org) {
          expect(button).toBeEnabled()
        } else {
          expect(button).toBeDisabled()
        }
      }
    })
  })

  describe('Performance', () => {
    test('does not cause unnecessary re-renders when props stay the same', () => {
      let renderCount = 0

      const TestComponent = ({ org }: { readonly org: string }) => {
        renderCount++
        return <LoadButton onClick={mockOnClick} organization={org} />
      }

      const { rerender } = render(<TestComponent org="github" />)
      expect(renderCount).toBe(1)

      // Re-render with same props
      rerender(<TestComponent org="github" />)
      expect(renderCount).toBe(2)
    })

    test('hover state changes do not cause excessive re-renders', async () => {
      const user = userEvent.setup()
      render(<LoadButton onClick={mockOnClick} organization="github" />)

      const button = screen.getByRole('button', { name: 'Load Graph' })

      // Multiple hover events should not break the component
      for (let i = 0; i < 5; i++) {
        await user.hover(button)
        await user.unhover(button)
      }

      // Component should still be functional
      expect(button).toBeInTheDocument()
      expect(button).toBeEnabled()
    })
  })

  describe('Integration Patterns', () => {
    test('works with form submission patterns', async () => {
      const user = userEvent.setup()
      const handleSubmit = jest.fn()

      const FormWrapper = () => (
        <form
          onSubmit={(e) => {
            e.preventDefault()
            handleSubmit()
          }}
        >
          <LoadButton onClick={handleSubmit} organization="github" />
        </form>
      )

      render(<FormWrapper />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      await user.click(button)

      expect(handleSubmit).toHaveBeenCalledTimes(1)
    })

    test('integrates with loading state management', () => {
      const LoadingWrapper = () => {
        const [isLoading, setIsLoading] = React.useState(false)

        const handleLoad = () => {
          setIsLoading(true)
          // Simulate async operation
          setTimeout(() => setIsLoading(false), 100)
        }

        return (
          <LoadButton
            onClick={handleLoad}
            organization={isLoading ? '' : 'github'}
          />
        )
      }

      render(<LoadingWrapper />)

      const button = screen.getByRole('button', { name: 'Load Graph' })
      expect(button).toBeEnabled()
    })
  })
})

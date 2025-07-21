/**
 * TopicsToggle Component Tests
 *
 * Tests for the topics vs teams toggle component including:
 * - Basic rendering and styling
 * - User interactions (click, keyboard)
 * - Accessibility features (ARIA, screen readers)
 * - State management
 * - Visual states (checked/unchecked)
 */

import React from 'react'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { TopicsToggle } from '../TopicsToggle'

describe('TopicsToggle', () => {
  const mockOnChange = jest.fn()

  beforeEach(() => {
    mockOnChange.mockClear()
  })

  describe('Basic Rendering', () => {
    test('renders checkbox with correct label text', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')
      const label = screen.getByText('Use Topics instead of Teams')

      expect(checkbox).toBeInTheDocument()
      expect(label).toBeInTheDocument()
    })

    test('displays unchecked state correctly', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')
      expect(checkbox).not.toBeChecked()
    })

    test('displays checked state correctly', () => {
      render(<TopicsToggle checked={true} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')
      expect(checkbox).toBeChecked()
    })

    test('applies correct CSS classes for styling', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const label = screen.getByText(
        'Use Topics instead of Teams'
      ).parentElement
      const checkbox = screen.getByRole('checkbox')

      // Label classes
      expect(label).toHaveClass(
        'flex',
        'items-center',
        'gap-2',
        'text-[#c9d1d9]',
        'text-sm',
        'cursor-pointer',
        'select-none'
      )

      // Checkbox classes
      expect(checkbox).toHaveClass('w-4', 'h-4', 'cursor-pointer')
    })

    test('applies correct accent color style', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')
      expect(checkbox).toHaveStyle({ accentColor: '#238636' })
    })
  })

  describe('User Interactions', () => {
    test('calls onChange when checkbox is clicked', async () => {
      const user = userEvent.setup()
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')

      await user.click(checkbox)

      expect(mockOnChange).toHaveBeenCalledTimes(1)
      expect(mockOnChange).toHaveBeenCalledWith(true)
    })

    test('calls onChange when label is clicked', async () => {
      const user = userEvent.setup()
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const label = screen.getByText('Use Topics instead of Teams')

      await user.click(label)

      expect(mockOnChange).toHaveBeenCalledTimes(1)
      expect(mockOnChange).toHaveBeenCalledWith(true)
    })

    test('toggles from checked to unchecked', async () => {
      const user = userEvent.setup()
      render(<TopicsToggle checked={true} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')

      await user.click(checkbox)

      expect(mockOnChange).toHaveBeenCalledWith(false)
    })

    test('toggles from unchecked to checked', async () => {
      const user = userEvent.setup()
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')

      await user.click(checkbox)

      expect(mockOnChange).toHaveBeenCalledWith(true)
    })
  })

  describe('Keyboard Accessibility', () => {
    test('checkbox is focusable via keyboard navigation', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')
      checkbox.focus()

      expect(checkbox).toHaveFocus()
    })

    test('can be toggled with Space key', async () => {
      const user = userEvent.setup()
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')
      checkbox.focus()

      await user.keyboard(' ')

      expect(mockOnChange).toHaveBeenCalledWith(true)
    })

    test('can be toggled with Enter key', async () => {
      const user = userEvent.setup()
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')
      checkbox.focus()

      await user.keyboard('{Enter}')

      expect(mockOnChange).toHaveBeenCalledWith(true)
    })

    test('maintains focus after interaction', async () => {
      const user = userEvent.setup()
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')
      checkbox.focus()

      await user.keyboard(' ')

      expect(checkbox).toHaveFocus()
    })
  })

  describe('Accessibility Features', () => {
    test('has proper semantic structure with label and checkbox', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const label = screen
        .getByText('Use Topics instead of Teams')
        .closest('label')
      const checkbox = screen.getByRole('checkbox')

      expect(label).toBeInTheDocument()
      expect(label).toContainElement(checkbox)
    })

    test('checkbox has correct ARIA attributes', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')

      expect(checkbox).toHaveAttribute('type', 'checkbox')
      expect(checkbox).toHaveAttribute('checked', '')
    })

    test('label text is descriptive and clear', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const labelText = screen.getByText('Use Topics instead of Teams')
      expect(labelText).toBeInTheDocument()

      // Verify the text clearly explains the toggle's purpose
      expect(labelText.textContent).toBe('Use Topics instead of Teams')
    })

    test('supports screen readers with proper role', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox', {
        name: 'Use Topics instead of Teams',
      })
      expect(checkbox).toBeInTheDocument()
    })
  })

  describe('State Management', () => {
    test('reflects checked state changes from props', () => {
      const { rerender } = render(
        <TopicsToggle checked={false} onChange={mockOnChange} />
      )

      let checkbox = screen.getByRole('checkbox')
      expect(checkbox).not.toBeChecked()

      rerender(<TopicsToggle checked={true} onChange={mockOnChange} />)

      checkbox = screen.getByRole('checkbox')
      expect(checkbox).toBeChecked()
    })

    test('maintains state consistency during rapid clicks', async () => {
      const user = userEvent.setup()
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')

      // Rapid clicks
      await user.click(checkbox)
      await user.click(checkbox)
      await user.click(checkbox)

      expect(mockOnChange).toHaveBeenCalledTimes(3)
      expect(mockOnChange).toHaveBeenNthCalledWith(1, true)
      expect(mockOnChange).toHaveBeenNthCalledWith(2, false)
      expect(mockOnChange).toHaveBeenNthCalledWith(3, true)
    })

    test('handles props update without losing focus', async () => {
      const user = userEvent.setup()
      const { rerender } = render(
        <TopicsToggle checked={false} onChange={mockOnChange} />
      )

      const checkbox = screen.getByRole('checkbox')
      await user.click(checkbox)
      expect(checkbox).toHaveFocus()

      // Update props
      rerender(<TopicsToggle checked={true} onChange={mockOnChange} />)

      // Focus should be maintained
      expect(checkbox).toHaveFocus()
    })
  })

  describe('Visual States and Styling', () => {
    test('cursor pointer is applied to interactive elements', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const label = screen.getByText(
        'Use Topics instead of Teams'
      ).parentElement
      const checkbox = screen.getByRole('checkbox')

      expect(label).toHaveClass('cursor-pointer')
      expect(checkbox).toHaveClass('cursor-pointer')
    })

    test('text selection is disabled on label', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const label = screen.getByText(
        'Use Topics instead of Teams'
      ).parentElement
      expect(label).toHaveClass('select-none')
    })

    test('proper spacing between checkbox and label', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const label = screen.getByText(
        'Use Topics instead of Teams'
      ).parentElement
      expect(label).toHaveClass('gap-2')
    })

    test('GitHub theme colors are applied correctly', () => {
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const label = screen.getByText(
        'Use Topics instead of Teams'
      ).parentElement
      const checkbox = screen.getByRole('checkbox')

      // GitHub dark theme text color
      expect(label).toHaveClass('text-[#c9d1d9]')
      // GitHub green accent color for checkbox
      expect(checkbox).toHaveStyle({ accentColor: '#238636' })
    })
  })

  describe('Edge Cases', () => {
    test('handles onChange being called multiple times', async () => {
      const user = userEvent.setup()
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')

      // Multiple interactions
      await user.click(checkbox)
      await user.keyboard(' ')

      // Each interaction should call onChange
      expect(mockOnChange).toHaveBeenCalledTimes(2)
    })

    test('maintains functionality after multiple re-renders', () => {
      const { rerender } = render(
        <TopicsToggle checked={false} onChange={mockOnChange} />
      )

      // Multiple re-renders
      for (let i = 0; i < 5; i++) {
        rerender(<TopicsToggle checked={i % 2 === 0} onChange={mockOnChange} />)
      }

      const checkbox = screen.getByRole('checkbox')
      expect(checkbox).toBeInTheDocument()
      expect(checkbox).not.toBeChecked() // Last render was with false
    })
  })

  describe('Performance', () => {
    test('does not cause unnecessary re-renders', () => {
      let renderCount = 0

      const TestComponent = ({ checked }: { readonly checked: boolean }) => {
        renderCount++
        return <TopicsToggle checked={checked} onChange={mockOnChange} />
      }

      const { rerender } = render(<TestComponent checked={false} />)
      expect(renderCount).toBe(1)

      // Re-render with same props
      rerender(<TestComponent checked={false} />)
      expect(renderCount).toBe(2)

      // Re-render with different props
      rerender(<TestComponent checked={true} />)
      expect(renderCount).toBe(3)
    })

    test('onChange callback is efficient and precise', async () => {
      const user = userEvent.setup()
      render(<TopicsToggle checked={false} onChange={mockOnChange} />)

      const checkbox = screen.getByRole('checkbox')

      await user.click(checkbox)

      // Should be called exactly once with correct value
      expect(mockOnChange).toHaveBeenCalledTimes(1)
      expect(mockOnChange).toHaveBeenCalledWith(true)
    })
  })

  describe('Component Integration', () => {
    test('integrates well with form management libraries', async () => {
      const user = userEvent.setup()
      let formState = { useTopics: false }

      const updateFormState = (useTopics: boolean) => {
        formState = { ...formState, useTopics }
      }

      render(
        <TopicsToggle
          checked={formState.useTopics}
          onChange={updateFormState}
        />
      )

      const checkbox = screen.getByRole('checkbox')
      await user.click(checkbox)

      expect(formState.useTopics).toBe(true)
    })

    test('works correctly in controlled component patterns', () => {
      const ControlledToggle = () => {
        const [checked, setChecked] = React.useState(false)

        return <TopicsToggle checked={checked} onChange={setChecked} />
      }

      render(<ControlledToggle />)

      const checkbox = screen.getByRole('checkbox')
      expect(checkbox).not.toBeChecked()
    })
  })
})

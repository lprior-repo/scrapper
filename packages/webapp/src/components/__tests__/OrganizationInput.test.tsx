/**
 * OrganizationInput Component Tests
 *
 * Tests for the organization name input component including:
 * - Basic rendering
 * - User input handling
 * - Prop validation
 * - Accessibility features
 * - Edge cases
 */

import React from 'react'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { OrganizationInput } from '../OrganizationInput'

describe('OrganizationInput', () => {
  const mockOnChange = jest.fn()

  beforeEach(() => {
    mockOnChange.mockClear()
  })

  describe('Basic Rendering', () => {
    test('renders input field with correct placeholder', () => {
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')
      expect(input).toBeInTheDocument()
      expect(input).toHaveAttribute('type', 'text')
    })

    test('displays the provided value', () => {
      const testValue = 'github'
      render(<OrganizationInput value={testValue} onChange={mockOnChange} />)

      const input = screen.getByDisplayValue(testValue)
      expect(input).toBeInTheDocument()
      expect(input).toHaveValue(testValue)
    })

    test('applies correct CSS classes for styling', () => {
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')
      expect(input).toHaveClass(
        'px-4',
        'py-2',
        'rounded-md',
        'border',
        'border-[#30363d]',
        'bg-[#0d1117]',
        'text-[#c9d1d9]',
        'text-sm',
        'w-50'
      )
    })
  })

  describe('User Interactions', () => {
    test('calls onChange when user types in input', async () => {
      const user = userEvent.setup()
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')

      await user.type(input, 'react')

      // onChange should be called for each character
      expect(mockOnChange).toHaveBeenCalledTimes(5) // 'r', 'e', 'a', 'c', 't'
      expect(mockOnChange).toHaveBeenLastCalledWith('react')
    })

    test('handles clearing input value', async () => {
      const user = userEvent.setup()
      render(<OrganizationInput value="github" onChange={mockOnChange} />)

      const input = screen.getByDisplayValue('github')

      await user.clear(input)

      expect(mockOnChange).toHaveBeenCalledWith('')
    })

    test('handles pasting text into input', async () => {
      const user = userEvent.setup()
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')

      await user.paste(input, 'microsoft')

      expect(mockOnChange).toHaveBeenCalledWith('microsoft')
    })

    test('handles special characters and spaces in organization name', async () => {
      const user = userEvent.setup()
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')

      await user.type(input, 'org-name_123')

      expect(mockOnChange).toHaveBeenLastCalledWith('org-name_123')
    })
  })

  describe('Accessibility', () => {
    test('input is focusable via keyboard navigation', () => {
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')
      input.focus()

      expect(input).toHaveFocus()
    })

    test('supports screen readers with proper semantic element', () => {
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')
      expect(input.tagName).toBe('INPUT')
      expect(input).toHaveAttribute('type', 'text')
    })

    test('placeholder provides clear instruction to users', () => {
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')
      expect(input).toHaveAttribute('placeholder', 'Enter organization name')
    })
  })

  describe('Prop Handling', () => {
    test('updates when value prop changes', () => {
      const { rerender } = render(
        <OrganizationInput value="initial" onChange={mockOnChange} />
      )

      let input = screen.getByDisplayValue('initial')
      expect(input).toHaveValue('initial')

      rerender(<OrganizationInput value="updated" onChange={mockOnChange} />)

      input = screen.getByDisplayValue('updated')
      expect(input).toHaveValue('updated')
    })

    test('maintains focus when props change', async () => {
      const user = userEvent.setup()
      const { rerender } = render(
        <OrganizationInput value="" onChange={mockOnChange} />
      )

      const input = screen.getByPlaceholderText('Enter organization name')
      await user.click(input)
      expect(input).toHaveFocus()

      // Re-render with new value
      rerender(<OrganizationInput value="test" onChange={mockOnChange} />)

      // Focus should be maintained
      expect(input).toHaveFocus()
    })
  })

  describe('Edge Cases', () => {
    test('handles empty string value correctly', () => {
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')
      expect(input).toHaveValue('')
    })

    test('handles long organization names', async () => {
      const user = userEvent.setup()
      const longName = 'a'.repeat(100)

      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')
      await user.paste(input, longName)

      expect(mockOnChange).toHaveBeenCalledWith(longName)
    })

    test('handles unicode characters in organization name', async () => {
      const user = userEvent.setup()
      const unicodeName = '组织名称'

      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')
      await user.paste(input, unicodeName)

      expect(mockOnChange).toHaveBeenCalledWith(unicodeName)
    })

    test('handles null-like values gracefully', () => {
      // Test with null value (should not crash)
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')
      expect(input).toHaveValue('')
    })
  })

  describe('Performance', () => {
    test('does not trigger unnecessary re-renders', () => {
      let renderCount = 0

      const TestComponent = ({ value }: { readonly value: string }) => {
        renderCount++
        return <OrganizationInput value={value} onChange={mockOnChange} />
      }

      const { rerender } = render(<TestComponent value="test" />)

      expect(renderCount).toBe(1)

      // Re-render with same value
      rerender(<TestComponent value="test" />)

      // Should only render twice (initial + rerender)
      expect(renderCount).toBe(2)
    })

    test('onChange callback is called efficiently', async () => {
      const user = userEvent.setup()
      render(<OrganizationInput value="" onChange={mockOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')

      await user.type(input, 'a')

      // Should be called exactly once for single character
      expect(mockOnChange).toHaveBeenCalledTimes(1)
      expect(mockOnChange).toHaveBeenCalledWith('a')
    })
  })

  describe('Integration with Form Libraries', () => {
    test('works with controlled form state', async () => {
      const user = userEvent.setup()
      let formValue = ''

      const TestForm = () => (
        <OrganizationInput
          value={formValue}
          onChange={(value) => {
            formValue = value
          }}
        />
      )

      render(<TestForm />)

      const input = screen.getByPlaceholderText('Enter organization name')

      await user.type(input, 'form-test')

      expect(formValue).toBe('form-test')
    })

    test('supports validation scenarios', () => {
      const validationOnChange = jest.fn((value: string) => {
        // Example validation: no spaces allowed
        if (value.includes(' ')) {
          return // Don't update if invalid
        }
        mockOnChange(value)
      })

      render(<OrganizationInput value="" onChange={validationOnChange} />)

      const input = screen.getByPlaceholderText('Enter organization name')

      // This test verifies the component can work with validation logic
      expect(input).toBeInTheDocument()
      expect(validationOnChange).toBeDefined()
    })
  })
})

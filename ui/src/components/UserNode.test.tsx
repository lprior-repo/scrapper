import React from 'react';
import { render, screen } from '@testing-library/react';
import { jest, describe, test, expect } from '@jest/globals';
import UserNode from './UserNode';

// Mock @xyflow/react
jest.mock('@xyflow/react', () => ({
  Handle: ({ type, position }: { type: string; position: string }) => (
    <div data-testid={`handle-${type}-${position}`} />
  ),
  Position: {
    Left: 'left',
    Right: 'right',
    Top: 'top',
    Bottom: 'bottom',
  },
}));

// Mock lucide-react
jest.mock('lucide-react', () => ({
  User: ({ size, style }: { size: number; style?: React.CSSProperties }) => (
    <div data-testid="user-icon" data-size={size} style={style} />
  ),
  Mail: ({ size, style }: { size: number; style?: React.CSSProperties }) => (
    <div data-testid="mail-icon" data-size={size} style={style} />
  ),
}));

describe('UserNode', () => {
  const mockUserData = {
    name: 'john-doe',
    type: 'user',
    platform: 'github',
  };

  const mockEmailData = {
    name: 'john.doe@company.com',
    email: 'john.doe@company.com',
    type: 'user',
    platform: 'github',
  };

  const defaultProps = {
    id: '1',
    data: mockUserData,
    selected: false,
    type: 'user',
    position: { x: 0, y: 0 },
    dragging: false,
    isConnectable: true,
    zIndex: 1,
    width: 160,
    height: 100,
    positionAbsolute: { x: 0, y: 0 },
  };

  test('renders user information correctly', () => {
    render(<UserNode {...defaultProps} />);
    
    expect(screen.getByText('User')).toBeInTheDocument();
    expect(screen.getByText('john-doe')).toBeInTheDocument();
    expect(screen.getByText('user')).toBeInTheDocument();
  });

  test('renders handles correctly', () => {
    render(<UserNode {...defaultProps} />);
    
    expect(screen.getByTestId('handle-target-left')).toBeInTheDocument();
    expect(screen.getByTestId('handle-source-right')).toBeInTheDocument();
  });

  test('renders user icon for regular users', () => {
    render(<UserNode {...defaultProps} />);
    
    expect(screen.getByTestId('user-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('mail-icon')).not.toBeInTheDocument();
  });

  test('renders mail icon for email users', () => {
    const emailProps = { ...defaultProps, data: mockEmailData };
    render(<UserNode {...emailProps} />);
    
    expect(screen.getByTestId('mail-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('user-icon')).not.toBeInTheDocument();
    expect(screen.getByText('Email')).toBeInTheDocument();
  });

  test('detects email by @ symbol in name', () => {
    const emailBySymbolData = {
      name: 'jane@example.com',
      type: 'user',
      platform: 'github',
    };
    const emailBySymbolProps = { ...defaultProps, data: emailBySymbolData };
    
    render(<UserNode {...emailBySymbolProps} />);
    
    expect(screen.getByTestId('mail-icon')).toBeInTheDocument();
    expect(screen.getByText('Email')).toBeInTheDocument();
  });

  test('applies selected styles when selected', () => {
    const selectedProps = { ...defaultProps, selected: true };
    const { container } = render(<UserNode {...selectedProps} />);
    
    const nodeElement = container.firstChild as HTMLElement;
    expect(nodeElement).toHaveStyle({
      background: '#fef3c7',
      border: '2px solid #f59e0b',
      boxShadow: '0 4px 12px rgba(245, 158, 11, 0.3)',
    });
  });

  test('applies default styles when not selected', () => {
    const { container } = render(<UserNode {...defaultProps} />);
    
    const nodeElement = container.firstChild as HTMLElement;
    expect(nodeElement).toHaveStyle({
      background: '#fffbeb',
      border: '1px solid #f59e0b',
      boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
    });
  });

  test('shows separate email when different from name', () => {
    const userWithSeparateEmail = {
      name: 'john-doe',
      email: 'john.doe@company.com',
      type: 'user',
      platform: 'github',
    };
    const propsWithSeparateEmail = { ...defaultProps, data: userWithSeparateEmail };
    
    render(<UserNode {...propsWithSeparateEmail} />);
    
    expect(screen.getByText('john-doe')).toBeInTheDocument();
    expect(screen.getByText('john.doe@company.com')).toBeInTheDocument();
  });

  test('does not show duplicate email when same as name', () => {
    const emailProps = { ...defaultProps, data: mockEmailData };
    render(<UserNode {...emailProps} />);
    
    const emailTexts = screen.getAllByText('john.doe@company.com');
    expect(emailTexts).toHaveLength(1); // Only shown once as name, not as separate email
  });

  test('handles long usernames with word break', () => {
    const longNameData = {
      name: 'very-long-username-that-might-overflow-the-container',
      type: 'user',
      platform: 'github',
    };
    const longNameProps = { ...defaultProps, data: longNameData };
    
    render(<UserNode {...longNameProps} />);
    
    expect(screen.getByText('very-long-username-that-might-overflow-the-container')).toBeInTheDocument();
  });

  test('handles long email addresses with word break', () => {
    const longEmailData = {
      name: 'user-with-very-long-email-address@very-long-domain-name.com',
      email: 'user-with-very-long-email-address@very-long-domain-name.com',
      type: 'user',
      platform: 'github',
    };
    const longEmailProps = { ...defaultProps, data: longEmailData };
    
    render(<UserNode {...longEmailProps} />);
    
    expect(screen.getByText('user-with-very-long-email-address@very-long-domain-name.com')).toBeInTheDocument();
  });

  test('renders different user types correctly', () => {
    const adminData = { ...mockUserData, type: 'admin' };
    const adminProps = { ...defaultProps, data: adminData };
    
    render(<UserNode {...adminProps} />);
    
    expect(screen.getByText('admin')).toBeInTheDocument();
  });
});
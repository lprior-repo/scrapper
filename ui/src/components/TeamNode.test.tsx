import React from 'react';
import { render, screen } from '@testing-library/react';
import { jest, describe, test, expect } from '@jest/globals';
import TeamNode from './TeamNode';

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
  Users: ({ size, style }: { size: number; style?: React.CSSProperties }) => (
    <div data-testid="users-icon" data-size={size} style={style} />
  ),
  Shield: ({ size, style }: { size: number; style?: React.CSSProperties }) => (
    <div data-testid="shield-icon" data-size={size} style={style} />
  ),
}));

describe('TeamNode', () => {
  const mockData = {
    name: 'backend-team',
    slug: 'backend-team',
    description: 'Backend development team',
    privacy: 'closed',
    memberCount: 5,
    platform: 'github',
  };

  const defaultProps = {
    id: '1',
    data: mockData,
    selected: false,
    type: 'team',
    position: { x: 0, y: 0 },
    dragging: false,
    isConnectable: true,
    zIndex: 1,
    width: 180,
    height: 120,
    positionAbsolute: { x: 0, y: 0 },
  };

  test('renders team information correctly', () => {
    render(<TeamNode {...defaultProps} />);
    
    expect(screen.getByText('Team')).toBeInTheDocument();
    expect(screen.getByText('backend-team')).toBeInTheDocument();
    expect(screen.getByText('Backend development team')).toBeInTheDocument();
    expect(screen.getByText('5 members')).toBeInTheDocument();
    expect(screen.getByText('closed')).toBeInTheDocument();
  });

  test('renders handles correctly', () => {
    render(<TeamNode {...defaultProps} />);
    
    expect(screen.getByTestId('handle-target-left')).toBeInTheDocument();
    expect(screen.getByTestId('handle-source-right')).toBeInTheDocument();
  });

  test('renders icons correctly', () => {
    render(<TeamNode {...defaultProps} />);
    
    const usersIcons = screen.getAllByTestId('users-icon');
    expect(usersIcons).toHaveLength(2); // One for Team label, one for member count
  });

  test('shows shield icon for private teams', () => {
    render(<TeamNode {...defaultProps} />);
    
    expect(screen.getByTestId('shield-icon')).toBeInTheDocument();
  });

  test('does not show shield icon for public teams', () => {
    const publicTeamData = { ...mockData, privacy: 'open' };
    const publicTeamProps = { ...defaultProps, data: publicTeamData };
    
    render(<TeamNode {...publicTeamProps} />);
    
    expect(screen.queryByTestId('shield-icon')).not.toBeInTheDocument();
  });

  test('applies selected styles when selected', () => {
    const selectedProps = { ...defaultProps, selected: true };
    const { container } = render(<TeamNode {...selectedProps} />);
    
    const nodeElement = container.firstChild as HTMLElement;
    expect(nodeElement).toHaveStyle({
      background: '#dbeafe',
      border: '2px solid #3b82f6',
      boxShadow: '0 4px 12px rgba(59, 130, 246, 0.3)',
    });
  });

  test('applies default styles when not selected', () => {
    const { container } = render(<TeamNode {...defaultProps} />);
    
    const nodeElement = container.firstChild as HTMLElement;
    expect(nodeElement).toHaveStyle({
      background: '#eff6ff',
      border: '1px solid #3b82f6',
      boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
    });
  });

  test('handles team without description', () => {
    const dataWithoutDescription = { ...mockData, description: '' };
    const propsWithoutDescription = { ...defaultProps, data: dataWithoutDescription };
    
    render(<TeamNode {...propsWithoutDescription} />);
    
    expect(screen.getByText('backend-team')).toBeInTheDocument();
    expect(screen.queryByText('Backend development team')).not.toBeInTheDocument();
  });

  test('renders different privacy levels correctly', () => {
    const publicTeamData = { ...mockData, privacy: 'open' };
    const publicTeamProps = { ...defaultProps, data: publicTeamData };
    
    render(<TeamNode {...publicTeamProps} />);
    
    const privacyBadge = screen.getByText('open');
    expect(privacyBadge).toHaveStyle({
      backgroundColor: '#ecfdf5',
      color: '#16a34a',
    });
  });

  test('renders private team privacy badge correctly', () => {
    render(<TeamNode {...defaultProps} />);
    
    const privacyBadge = screen.getByText('closed');
    expect(privacyBadge).toHaveStyle({
      backgroundColor: '#fee2e2',
      color: '#dc2626',
    });
  });

  test('handles secret privacy level', () => {
    const secretTeamData = { ...mockData, privacy: 'secret' };
    const secretTeamProps = { ...defaultProps, data: secretTeamData };
    
    render(<TeamNode {...secretTeamProps} />);
    
    expect(screen.getByTestId('shield-icon')).toBeInTheDocument();
    const privacyBadge = screen.getByText('secret');
    expect(privacyBadge).toHaveStyle({
      backgroundColor: '#fee2e2',
      color: '#dc2626',
    });
  });

  test('handles different member counts', () => {
    const singleMemberData = { ...mockData, memberCount: 1 };
    const singleMemberProps = { ...defaultProps, data: singleMemberData };
    
    render(<TeamNode {...singleMemberProps} />);
    
    expect(screen.getByText('1 members')).toBeInTheDocument();
  });

  test('handles zero member count', () => {
    const zeroMemberData = { ...mockData, memberCount: 0 };
    const zeroMemberProps = { ...defaultProps, data: zeroMemberData };
    
    render(<TeamNode {...zeroMemberProps} />);
    
    expect(screen.getByText('0 members')).toBeInTheDocument();
  });
});
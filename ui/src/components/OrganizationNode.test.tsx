import React from 'react';
import { render, screen } from '@testing-library/react';
import { jest, describe, test, expect } from '@jest/globals';
import OrganizationNode from './OrganizationNode';

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
  Building2: ({ size, style }: { size: number; style?: React.CSSProperties }) => (
    <div data-testid="building-icon" data-size={size} style={style} />
  ),
  Github: ({ size, style }: { size: number; style?: React.CSSProperties }) => (
    <div data-testid="github-icon" data-size={size} style={style} />
  ),
}));

describe('OrganizationNode', () => {
  const mockData = {
    name: 'test-org',
    type: 'github',
    platform: 'github.com',
  };

  const defaultProps = {
    id: '1',
    data: mockData,
    selected: false,
    type: 'organization',
    position: { x: 0, y: 0 },
    dragging: false,
    isConnectable: true,
    zIndex: 1,
    width: 220,
    height: 120,
    positionAbsolute: { x: 0, y: 0 },
  };

  test('renders organization information correctly', () => {
    render(<OrganizationNode {...defaultProps} />);
    
    expect(screen.getByText('Organization')).toBeInTheDocument();
    expect(screen.getByText('test-org')).toBeInTheDocument();
    expect(screen.getByText('github.com • github')).toBeInTheDocument();
  });

  test('renders handle correctly', () => {
    render(<OrganizationNode {...defaultProps} />);
    
    expect(screen.getByTestId('handle-source-bottom')).toBeInTheDocument();
  });

  test('renders icons correctly', () => {
    render(<OrganizationNode {...defaultProps} />);
    
    expect(screen.getByTestId('building-icon')).toBeInTheDocument();
    expect(screen.getByTestId('github-icon')).toBeInTheDocument();
  });

  test('applies selected styles when selected', () => {
    const selectedProps = { ...defaultProps, selected: true };
    const { container } = render(<OrganizationNode {...selectedProps} />);
    
    const nodeElement = container.firstChild as HTMLElement;
    expect(nodeElement).toHaveStyle({
      background: '#ede9fe',
      border: '2px solid #8b5cf6',
      boxShadow: '0 4px 12px rgba(139, 92, 246, 0.3)',
    });
  });

  test('applies default styles when not selected', () => {
    const { container } = render(<OrganizationNode {...defaultProps} />);
    
    const nodeElement = container.firstChild as HTMLElement;
    expect(nodeElement).toHaveStyle({
      background: '#f5f3ff',
      border: '1px solid #8b5cf6',
      boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
    });
  });

  test('renders organization name with center alignment', () => {
    render(<OrganizationNode {...defaultProps} />);
    
    const orgName = screen.getByText('test-org');
    expect(orgName).toHaveStyle({
      textAlign: 'center',
      fontSize: '20px',
      fontWeight: '700',
    });
  });

  test('renders platform and type with proper formatting', () => {
    render(<OrganizationNode {...defaultProps} />);
    
    const platformType = screen.getByText('github.com • github');
    expect(platformType).toHaveStyle({
      textAlign: 'center',
      textTransform: 'uppercase',
      letterSpacing: '0.5px',
      fontSize: '12px',
    });
  });

  test('handles different organization names', () => {
    const customOrgData = { ...mockData, name: 'my-awesome-org' };
    const customOrgProps = { ...defaultProps, data: customOrgData };
    
    render(<OrganizationNode {...customOrgProps} />);
    
    expect(screen.getByText('my-awesome-org')).toBeInTheDocument();
    expect(screen.queryByText('test-org')).not.toBeInTheDocument();
  });

  test('handles different platforms', () => {
    const gitlabData = { 
      ...mockData, 
      platform: 'gitlab.com',
      type: 'gitlab' 
    };
    const gitlabProps = { ...defaultProps, data: gitlabData };
    
    render(<OrganizationNode {...gitlabProps} />);
    
    expect(screen.getByText('gitlab.com • gitlab')).toBeInTheDocument();
  });

  test('has decorative elements positioned correctly', () => {
    const { container } = render(<OrganizationNode {...defaultProps} />);
    
    const decorativeElements = container.querySelectorAll('div[style*="position: absolute"]');
    expect(decorativeElements).toHaveLength(2);
    
    // Check top-right decorative element
    const topRightElement = decorativeElements[0] as HTMLElement;
    expect(topRightElement).toHaveStyle({
      top: '8px',
      right: '8px',
      width: '6px',
      height: '6px',
      borderRadius: '50%',
    });
    
    // Check bottom-left decorative element
    const bottomLeftElement = decorativeElements[1] as HTMLElement;
    expect(bottomLeftElement).toHaveStyle({
      bottom: '8px',
      left: '8px',
      width: '4px',
      height: '4px',
      borderRadius: '50%',
    });
  });

  test('handles long organization names', () => {
    const longNameData = { 
      ...mockData, 
      name: 'very-long-organization-name-that-might-overflow' 
    };
    const longNameProps = { ...defaultProps, data: longNameData };
    
    render(<OrganizationNode {...longNameProps} />);
    
    expect(screen.getByText('very-long-organization-name-that-might-overflow')).toBeInTheDocument();
  });
});
import React from 'react';
import { render, screen } from '@testing-library/react';
import { jest, describe, test, expect } from '@jest/globals';
import RepositoryNode from './RepositoryNode';

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
  GitBranch: ({ size, style }: { size: number; style?: React.CSSProperties }) => (
    <div data-testid="git-branch-icon" data-size={size} style={style} />
  ),
  FileText: ({ size, style }: { size: number; style?: React.CSSProperties }) => (
    <div data-testid="file-text-icon" data-size={size} style={style} />
  ),
}));

describe('RepositoryNode', () => {
  const mockData = {
    name: 'test-repo',
    fullName: 'test-org/test-repo',
    defaultBranch: 'main',
    hasCodeowners: true,
    platform: 'github',
  };

  const defaultProps = {
    id: '1',
    data: mockData,
    selected: false,
    type: 'repository',
    position: { x: 0, y: 0 },
    dragging: false,
    isConnectable: true,
    zIndex: 1,
    width: 200,
    height: 100,
    positionAbsolute: { x: 0, y: 0 },
  };

  test('renders repository information correctly', () => {
    render(<RepositoryNode {...defaultProps} />);
    
    expect(screen.getByText('Repository')).toBeInTheDocument();
    expect(screen.getByText('test-repo')).toBeInTheDocument();
    expect(screen.getByText('test-org/test-repo')).toBeInTheDocument();
    expect(screen.getByText('main')).toBeInTheDocument();
    expect(screen.getByText('CODEOWNERS')).toBeInTheDocument();
  });

  test('renders handles correctly', () => {
    render(<RepositoryNode {...defaultProps} />);
    
    expect(screen.getByTestId('handle-target-left')).toBeInTheDocument();
    expect(screen.getByTestId('handle-source-right')).toBeInTheDocument();
  });

  test('renders icons correctly', () => {
    render(<RepositoryNode {...defaultProps} />);
    
    const gitBranchIcons = screen.getAllByTestId('git-branch-icon');
    expect(gitBranchIcons).toHaveLength(2); // One for Repository label, one for branch
    
    const fileTextIcon = screen.getByTestId('file-text-icon');
    expect(fileTextIcon).toBeInTheDocument();
  });

  test('applies selected styles when selected', () => {
    const selectedProps = { ...defaultProps, selected: true };
    const { container } = render(<RepositoryNode {...selectedProps} />);
    
    const nodeElement = container.firstChild as HTMLElement;
    expect(nodeElement).toHaveStyle({
      background: '#dcfce7',
      border: '2px solid #10b981',
      boxShadow: '0 4px 12px rgba(16, 185, 129, 0.3)',
    });
  });

  test('applies default styles when not selected', () => {
    const { container } = render(<RepositoryNode {...defaultProps} />);
    
    const nodeElement = container.firstChild as HTMLElement;
    expect(nodeElement).toHaveStyle({
      background: '#f0fdf4',
      border: '1px solid #10b981',
      boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
    });
  });

  test('hides CODEOWNERS when not present', () => {
    const dataWithoutCodeowners = { ...mockData, hasCodeowners: false };
    const propsWithoutCodeowners = { ...defaultProps, data: dataWithoutCodeowners };
    
    render(<RepositoryNode {...propsWithoutCodeowners} />);
    
    expect(screen.queryByText('CODEOWNERS')).not.toBeInTheDocument();
    expect(screen.queryByTestId('file-text-icon')).not.toBeInTheDocument();
  });

  test('renders different branch names correctly', () => {
    const dataWithMasterBranch = { ...mockData, defaultBranch: 'master' };
    const propsWithMasterBranch = { ...defaultProps, data: dataWithMasterBranch };
    
    render(<RepositoryNode {...propsWithMasterBranch} />);
    
    expect(screen.getByText('master')).toBeInTheDocument();
    expect(screen.queryByText('main')).not.toBeInTheDocument();
  });

  test('renders repository with long names', () => {
    const dataWithLongName = {
      ...mockData,
      name: 'very-long-repository-name-that-might-overflow',
      fullName: 'organization-with-long-name/very-long-repository-name-that-might-overflow',
    };
    const propsWithLongName = { ...defaultProps, data: dataWithLongName };
    
    render(<RepositoryNode {...propsWithLongName} />);
    
    expect(screen.getByText('very-long-repository-name-that-might-overflow')).toBeInTheDocument();
    expect(screen.getByText('organization-with-long-name/very-long-repository-name-that-might-overflow')).toBeInTheDocument();
  });
});
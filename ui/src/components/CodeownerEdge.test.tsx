import React from 'react';
import { render } from '@testing-library/react';
import { jest, describe, test, expect } from '@jest/globals';
import CodeownerEdge from './CodeownerEdge';

// Mock @xyflow/react
jest.mock('@xyflow/react', () => ({
  getStraightPath: jest.fn(() => ['M 0 0 L 100 100', 50, 50]),
  EdgeLabelRenderer: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="edge-label-renderer">{children}</div>
  ),
}));

describe('CodeownerEdge', () => {
  const defaultProps = {
    id: 'edge-1',
    source: 'node-1',
    target: 'node-2',
    sourceX: 0,
    sourceY: 0,
    targetX: 100,
    targetY: 100,
    sourcePosition: 'right' as const,
    targetPosition: 'left' as const,
    selected: false,
    animated: false,
    style: {},
    markerStart: '',
    markerEnd: '',
    selectable: true,
    deletable: true,
  };

  test('renders edge path correctly', () => {
    const props = {
      ...defaultProps,
      data: { relationshipType: 'repository_codeowner' },
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toBeInTheDocument();
    expect(path).toHaveAttribute('d', 'M 0 0 L 100 100');
  });

  test('applies correct color for repository_codeowner relationship', () => {
    const props = {
      ...defaultProps,
      data: { relationshipType: 'repository_codeowner' },
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toHaveStyle({ stroke: '#ef4444' });
  });

  test('applies correct color for organization_repository relationship', () => {
    const props = {
      ...defaultProps,
      data: { relationshipType: 'organization_repository' },
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toHaveStyle({ stroke: '#8b5cf6' });
  });

  test('applies correct color for organization_team relationship', () => {
    const props = {
      ...defaultProps,
      data: { relationshipType: 'organization_team' },
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toHaveStyle({ stroke: '#3b82f6' });
  });

  test('applies default color for unknown relationship', () => {
    const props = {
      ...defaultProps,
      data: { relationshipType: 'unknown_type' },
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toHaveStyle({ stroke: '#6b7280' });
  });

  test('applies dashed line for repository_codeowner relationship', () => {
    const props = {
      ...defaultProps,
      data: { relationshipType: 'repository_codeowner' },
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toHaveAttribute('stroke-dasharray', '5,5');
  });

  test('applies solid line for non-codeowner relationships', () => {
    const props = {
      ...defaultProps,
      data: { relationshipType: 'organization_repository' },
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toHaveAttribute('stroke-dasharray', 'none');
  });

  test('applies selected styles when selected', () => {
    const props = {
      ...defaultProps,
      selected: true,
      data: { relationshipType: 'repository_codeowner' },
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toHaveStyle({ 
      strokeWidth: '3',
      opacity: '1' 
    });
  });

  test('applies default styles when not selected', () => {
    const props = {
      ...defaultProps,
      selected: false,
      data: { relationshipType: 'repository_codeowner' },
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toHaveStyle({ 
      strokeWidth: '2',
      opacity: '0.8' 
    });
  });

  test('does not render label when not selected', () => {
    const props = {
      ...defaultProps,
      selected: false,
      data: { relationshipType: 'repository_codeowner' },
    };
    
    const { queryByTestId } = render(<CodeownerEdge {...props} />);
    
    expect(queryByTestId('edge-label-renderer')).not.toBeInTheDocument();
  });

  test('renders label when selected', () => {
    const props = {
      ...defaultProps,
      selected: true,
      data: { relationshipType: 'repository_codeowner' },
    };
    
    const { getByTestId, getByText } = render(<CodeownerEdge {...props} />);
    
    expect(getByTestId('edge-label-renderer')).toBeInTheDocument();
    expect(getByText('CODEOWNER')).toBeInTheDocument();
  });

  test('renders correct labels for different relationship types', () => {
    const testCases = [
      { type: 'repository_codeowner', label: 'CODEOWNER' },
      { type: 'organization_repository', label: 'OWNS' },
      { type: 'organization_team', label: 'CONTAINS' },
      { type: 'custom_type', label: 'CUSTOM_TYPE' },
    ];

    testCases.forEach(({ type, label }) => {
      const props = {
        ...defaultProps,
        selected: true,
        data: { relationshipType: type },
      };
      
      const { getByText } = render(<CodeownerEdge {...props} />);
      expect(getByText(label)).toBeInTheDocument();
    });
  });

  test('renders pattern when provided in data', () => {
    const props = {
      ...defaultProps,
      selected: true,
      data: { 
        relationshipType: 'repository_codeowner',
        pattern: '*.js'
      },
    };
    
    const { getByText } = render(<CodeownerEdge {...props} />);
    
    expect(getByText('CODEOWNER')).toBeInTheDocument();
    expect(getByText('*.js')).toBeInTheDocument();
  });

  test('does not render pattern when not provided', () => {
    const props = {
      ...defaultProps,
      selected: true,
      data: { relationshipType: 'repository_codeowner' },
    };
    
    const { getByText, queryByText } = render(<CodeownerEdge {...props} />);
    
    expect(getByText('CODEOWNER')).toBeInTheDocument();
    expect(queryByText('*.js')).not.toBeInTheDocument();
  });

  test('handles missing data gracefully', () => {
    const props = {
      ...defaultProps,
      selected: false,
      data: undefined,
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toBeInTheDocument();
    expect(path).toHaveStyle({ stroke: '#6b7280' }); // default color
  });

  test('handles missing relationship type gracefully', () => {
    const props = {
      ...defaultProps,
      selected: true,
      data: { pattern: '*.js' },
    };
    
    const { container } = render(<CodeownerEdge {...props} />);
    
    const path = container.querySelector('path');
    expect(path).toHaveStyle({ stroke: '#6b7280' }); // default color
  });
});
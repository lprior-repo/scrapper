import React from 'react';
import { EdgeProps, getStraightPath, EdgeLabelRenderer } from '@xyflow/react';

interface CodeownerEdgeData {
  pattern?: string;
  ownerReference?: string;
  relationshipType: string;
}

const CodeownerEdge: React.FC<EdgeProps<CodeownerEdgeData>> = ({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  data,
  selected,
}) => {
  const [edgePath, labelX, labelY] = getStraightPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  });

  const getEdgeColor = (relationshipType: string) => {
    switch (relationshipType) {
      case 'repository_codeowner':
        return '#ef4444'; // red for codeowner relationships
      case 'organization_repository':
        return '#8b5cf6'; // purple for ownership
      case 'organization_team':
        return '#3b82f6'; // blue for contains
      default:
        return '#6b7280'; // gray for others
    }
  };

  const getEdgeLabel = (relationshipType: string) => {
    switch (relationshipType) {
      case 'repository_codeowner':
        return 'CODEOWNER';
      case 'organization_repository':
        return 'OWNS';
      case 'organization_team':
        return 'CONTAINS';
      default:
        return relationshipType.toUpperCase();
    }
  };

  const edgeColor = getEdgeColor(data?.relationshipType || '');
  const strokeWidth = selected ? 3 : 2;
  const opacity = selected ? 1 : 0.8;

  return (
    <>
      <path
        id={id}
        d={edgePath}
        stroke={edgeColor}
        strokeWidth={strokeWidth}
        fill="none"
        opacity={opacity}
        strokeDasharray={data?.relationshipType === 'repository_codeowner' ? '5,5' : 'none'}
        style={{
          transition: 'all 0.2s ease',
        }}
      />
      
      {selected && data && (
        <EdgeLabelRenderer>
          <div
            style={{
              position: 'absolute',
              transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
              background: 'white',
              padding: '4px 8px',
              borderRadius: '4px',
              fontSize: '10px',
              fontWeight: '600',
              color: edgeColor,
              border: `1px solid ${edgeColor}`,
              boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
              pointerEvents: 'none',
              maxWidth: '120px',
              textAlign: 'center',
              lineHeight: '1.2',
            }}
            className="nodrag nopan"
          >
            <div>{getEdgeLabel(data.relationshipType)}</div>
            {data.pattern && (
              <div style={{ fontSize: '9px', opacity: 0.8, marginTop: '2px' }}>
                {data.pattern}
              </div>
            )}
          </div>
        </EdgeLabelRenderer>
      )}
    </>
  );
};

export default CodeownerEdge;
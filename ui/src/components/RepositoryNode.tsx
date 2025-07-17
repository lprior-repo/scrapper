import React from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { GitBranch, FileText } from 'lucide-react';

interface RepositoryData {
  name: string;
  fullName: string;
  defaultBranch: string;
  hasCodeowners: boolean;
  platform: string;
}

const RepositoryNode: React.FC<NodeProps<RepositoryData>> = ({ data, selected }) => {
  return (
    <div
      style={{
        background: selected ? '#dcfce7' : '#f0fdf4',
        border: selected ? '2px solid #10b981' : '1px solid #10b981',
        borderRadius: '8px',
        padding: '12px',
        minWidth: '200px',
        boxShadow: selected ? '0 4px 12px rgba(16, 185, 129, 0.3)' : '0 2px 8px rgba(0, 0, 0, 0.1)',
        transition: 'all 0.2s ease',
      }}
    >
      <Handle
        type="target"
        position={Position.Left}
        style={{ background: '#10b981', width: '8px', height: '8px' }}
      />
      
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
        <GitBranch size={16} style={{ color: '#10b981', marginRight: '8px' }} />
        <div style={{ fontSize: '14px', fontWeight: 'bold', color: '#065f46' }}>
          Repository
        </div>
      </div>
      
      <div style={{ fontSize: '16px', fontWeight: '600', color: '#1f2937', marginBottom: '4px' }}>
        {data.name}
      </div>
      
      <div style={{ fontSize: '12px', color: '#6b7280', marginBottom: '8px' }}>
        {data.fullName}
      </div>
      
      <div style={{ display: 'flex', alignItems: 'center', gap: '12px', fontSize: '11px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
          <GitBranch size={12} style={{ color: '#6b7280' }} />
          <span style={{ color: '#6b7280' }}>{data.defaultBranch}</span>
        </div>
        
        {data.hasCodeowners && (
          <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
            <FileText size={12} style={{ color: '#10b981' }} />
            <span style={{ color: '#10b981', fontWeight: '500' }}>CODEOWNERS</span>
          </div>
        )}
      </div>
      
      <Handle
        type="source"
        position={Position.Right}
        style={{ background: '#10b981', width: '8px', height: '8px' }}
      />
    </div>
  );
};

export default RepositoryNode;
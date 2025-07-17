import React from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Users, Shield } from 'lucide-react';

interface TeamData {
  name: string;
  slug: string;
  description: string;
  privacy: string;
  memberCount: number;
  platform: string;
}

const TeamNode: React.FC<NodeProps<TeamData>> = ({ data, selected }) => {
  const isPrivate = data.privacy === 'closed' || data.privacy === 'secret';
  
  return (
    <div
      style={{
        background: selected ? '#dbeafe' : '#eff6ff',
        border: selected ? '2px solid #3b82f6' : '1px solid #3b82f6',
        borderRadius: '8px',
        padding: '12px',
        minWidth: '180px',
        boxShadow: selected ? '0 4px 12px rgba(59, 130, 246, 0.3)' : '0 2px 8px rgba(0, 0, 0, 0.1)',
        transition: 'all 0.2s ease',
      }}
    >
      <Handle
        type="target"
        position={Position.Left}
        style={{ background: '#3b82f6', width: '8px', height: '8px' }}
      />
      
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
        <Users size={16} style={{ color: '#3b82f6', marginRight: '8px' }} />
        <div style={{ fontSize: '14px', fontWeight: 'bold', color: '#1e40af' }}>
          Team
        </div>
        {isPrivate && (
          <Shield size={12} style={{ color: '#ef4444', marginLeft: '4px' }} />
        )}
      </div>
      
      <div style={{ fontSize: '16px', fontWeight: '600', color: '#1f2937', marginBottom: '4px' }}>
        {data.name}
      </div>
      
      {data.description && (
        <div style={{ 
          fontSize: '12px', 
          color: '#6b7280', 
          marginBottom: '8px',
          lineHeight: '1.3',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          display: '-webkit-box',
          WebkitLineClamp: 2,
          WebkitBoxOrient: 'vertical',
        }}>
          {data.description}
        </div>
      )}
      
      <div style={{ display: 'flex', alignItems: 'center', gap: '12px', fontSize: '11px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
          <Users size={12} style={{ color: '#6b7280' }} />
          <span style={{ color: '#6b7280' }}>{data.memberCount} members</span>
        </div>
        
        <div style={{ 
          padding: '2px 6px', 
          borderRadius: '4px', 
          backgroundColor: isPrivate ? '#fee2e2' : '#ecfdf5',
          color: isPrivate ? '#dc2626' : '#16a34a',
          fontSize: '10px',
          fontWeight: '500',
        }}>
          {data.privacy}
        </div>
      </div>
      
      <Handle
        type="source"
        position={Position.Right}
        style={{ background: '#3b82f6', width: '8px', height: '8px' }}
      />
    </div>
  );
};

export default TeamNode;
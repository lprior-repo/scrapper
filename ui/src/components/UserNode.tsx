import React from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { User, Mail } from 'lucide-react';

interface UserData {
  name: string;
  email?: string;
  type: string;
  platform: string;
}

const UserNode: React.FC<NodeProps<UserData>> = ({ data, selected }) => {
  const isEmail = data.email || data.name.includes('@');
  
  return (
    <div
      style={{
        background: selected ? '#fef3c7' : '#fffbeb',
        border: selected ? '2px solid #f59e0b' : '1px solid #f59e0b',
        borderRadius: '8px',
        padding: '12px',
        minWidth: '160px',
        boxShadow: selected ? '0 4px 12px rgba(245, 158, 11, 0.3)' : '0 2px 8px rgba(0, 0, 0, 0.1)',
        transition: 'all 0.2s ease',
      }}
    >
      <Handle
        type="target"
        position={Position.Left}
        style={{ background: '#f59e0b', width: '8px', height: '8px' }}
      />
      
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
        {isEmail ? (
          <Mail size={16} style={{ color: '#f59e0b', marginRight: '8px' }} />
        ) : (
          <User size={16} style={{ color: '#f59e0b', marginRight: '8px' }} />
        )}
        <div style={{ fontSize: '14px', fontWeight: 'bold', color: '#92400e' }}>
          {isEmail ? 'Email' : 'User'}
        </div>
      </div>
      
      <div style={{ 
        fontSize: '16px', 
        fontWeight: '600', 
        color: '#1f2937', 
        marginBottom: '4px',
        wordBreak: 'break-word',
      }}>
        {data.name}
      </div>
      
      {data.email && data.email !== data.name && (
        <div style={{ 
          fontSize: '12px', 
          color: '#6b7280', 
          marginBottom: '8px',
          wordBreak: 'break-word',
        }}>
          {data.email}
        </div>
      )}
      
      <div style={{ 
        padding: '2px 6px', 
        borderRadius: '4px', 
        backgroundColor: '#fef3c7',
        color: '#92400e',
        fontSize: '10px',
        fontWeight: '500',
        display: 'inline-block',
      }}>
        {data.type}
      </div>
      
      <Handle
        type="source"
        position={Position.Right}
        style={{ background: '#f59e0b', width: '8px', height: '8px' }}
      />
    </div>
  );
};

export default UserNode;
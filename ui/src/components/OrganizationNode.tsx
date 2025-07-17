import React from 'react';
import { Handle, Position, NodeProps } from '@xyflow/react';
import { Building2, Github } from 'lucide-react';

interface OrganizationData {
  name: string;
  type: string;
  platform: string;
}

const OrganizationNode: React.FC<NodeProps<OrganizationData>> = ({ data, selected }) => {
  return (
    <div
      style={{
        background: selected ? '#ede9fe' : '#f5f3ff',
        border: selected ? '2px solid #8b5cf6' : '1px solid #8b5cf6',
        borderRadius: '12px',
        padding: '16px',
        minWidth: '220px',
        boxShadow: selected ? '0 4px 12px rgba(139, 92, 246, 0.3)' : '0 2px 8px rgba(0, 0, 0, 0.1)',
        transition: 'all 0.2s ease',
        position: 'relative',
      }}
    >
      <Handle
        type="source"
        position={Position.Bottom}
        style={{ background: '#8b5cf6', width: '10px', height: '10px' }}
      />
      
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '12px' }}>
        <Building2 size={20} style={{ color: '#8b5cf6', marginRight: '8px' }} />
        <div style={{ fontSize: '16px', fontWeight: 'bold', color: '#6b21a8' }}>
          Organization
        </div>
        <Github size={16} style={{ color: '#8b5cf6', marginLeft: 'auto' }} />
      </div>
      
      <div style={{ 
        fontSize: '20px', 
        fontWeight: '700', 
        color: '#1f2937', 
        marginBottom: '8px',
        textAlign: 'center',
      }}>
        {data.name}
      </div>
      
      <div style={{ 
        fontSize: '12px', 
        color: '#6b7280', 
        textAlign: 'center',
        textTransform: 'uppercase',
        letterSpacing: '0.5px',
        fontWeight: '500',
      }}>
        {data.platform} • {data.type}
      </div>
      
      {/* Decorative elements */}
      <div style={{
        position: 'absolute',
        top: '8px',
        right: '8px',
        width: '6px',
        height: '6px',
        borderRadius: '50%',
        background: '#8b5cf6',
        opacity: 0.3,
      }} />
      
      <div style={{
        position: 'absolute',
        bottom: '8px',
        left: '8px',
        width: '4px',
        height: '4px',
        borderRadius: '50%',
        background: '#8b5cf6',
        opacity: 0.3,
      }} />
    </div>
  );
};

export default OrganizationNode;
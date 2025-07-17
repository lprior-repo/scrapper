import React, { useCallback, useState, useEffect } from 'react';
import {
  ReactFlow,
  Node,
  addEdge,
  Background,
  Controls,
  MiniMap,
  useNodesState,
  useEdgesState,
  Connection,
  NodeTypes,
  EdgeTypes,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import axios from 'axios';

import RepositoryNode from './components/RepositoryNode';
import TeamNode from './components/TeamNode';
import UserNode from './components/UserNode';
import OrganizationNode from './components/OrganizationNode';
import CodeownerEdge from './components/CodeownerEdge';
import './App.css';

const API_URL = 'http://localhost:8081/api';

// Define custom node types
const nodeTypes: NodeTypes = {
  repository: RepositoryNode,
  team: TeamNode,
  user: UserNode,
  organization: OrganizationNode,
};

// Define custom edge types
const edgeTypes: EdgeTypes = {
  codeowner: CodeownerEdge,
  owns: CodeownerEdge,
  contains: CodeownerEdge,
};

function App() {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [orgName, setOrgName] = useState('');
  const [githubToken, setGithubToken] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [stats, setStats] = useState<Record<string, number | string> | null>(null);

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  );

  // Load GitHub token from localStorage on mount
  useEffect(() => {
    const savedToken = localStorage.getItem('github_token');
    if (savedToken) {
      setGithubToken(savedToken);
    }
  }, []);

  // Scan GitHub organization
  const scanOrganization = async () => {
    if (!orgName || !githubToken) {
      setError('Please enter organization name and GitHub token');
      return;
    }

    setLoading(true);
    setError('');

    try {
      // Save token to localStorage
      localStorage.setItem('github_token', githubToken);

      // Scan the organization
      await axios.post(
        `${API_URL}/scan/${orgName}`,
        {},
        {
          headers: {
            'X-GitHub-Token': githubToken
          }
        }
      );


      // Load the graph data
      await loadGraphData();
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to scan organization';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // Load graph data for organization
  const loadGraphData = async () => {
    if (!orgName) return;

    setLoading(true);
    setError('');

    try {
      // Get graph data
      const graphResponse = await axios.get(`${API_URL}/graph/${orgName}`);
      const { nodes: newNodes, edges: newEdges } = graphResponse.data;

      // Set nodes and edges
      if (newNodes && newEdges) {
        // Update node and edge states
        const nodesWithIds = newNodes.map((node: Record<string, unknown>) => ({
          ...node,
          id: node.id || String(Math.random())
        }));
        const edgesWithIds = newEdges.map((edge: Record<string, unknown>) => ({
          ...edge,
          id: edge.id || String(Math.random())
        }));

        // Set new data directly
        setNodes(nodesWithIds);
        setEdges(edgesWithIds);
        
        // Fit view after a short delay
        setTimeout(() => {
          // @ts-expect-error reactFlowInstance is added to window
          window.reactFlowInstance?.fitView({ padding: 0.2 });
        }, 100);
      }

      // Get stats
      const statsResponse = await axios.get(`${API_URL}/stats/${orgName}`);
      setStats(statsResponse.data);
    } catch (err: unknown) {
      if (err instanceof Error && 'response' in err && (err as Record<string, unknown>).response && typeof (err as Record<string, unknown>).response === 'object' && (err as Record<string, unknown>).response !== null && 'status' in (err as Record<string, unknown>).response && ((err as Record<string, unknown>).response as Record<string, unknown>).status === 404) {
        setError('Organization not found. Please scan it first.');
      } else {
        const errorMessage = err instanceof Error ? err.message : 'Failed to load graph data';
        setError(errorMessage);
      }
    } finally {
      setLoading(false);
    }
  };

  // Mini map node colors
  const nodeColor = (node: Node) => {
    switch (node.type) {
      case 'repository': return '#10b981'; // green
      case 'team': return '#3b82f6'; // blue
      case 'user': return '#f59e0b'; // amber
      case 'organization': return '#8b5cf6'; // purple
      default: return '#6b7280'; // gray
    }
  };

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <div style={{ 
        position: 'absolute', 
        top: 0, 
        left: 0, 
        right: 0, 
        zIndex: 10, 
        background: 'rgba(255, 255, 255, 0.95)',
        padding: '1rem',
        borderBottom: '1px solid #e5e7eb',
        backdropFilter: 'blur(8px)'
      }}>
        <h1 style={{ 
          margin: 0, 
          fontSize: '1.5rem', 
          fontWeight: 'bold', 
          color: '#1f2937' 
        }}>
          GitHub Codeowner Visualization
        </h1>
        <p style={{ 
          margin: '0.25rem 0 0 0', 
          color: '#6b7280', 
          fontSize: '0.875rem' 
        }}>
          Interactive graph showing repositories, teams, users, and their codeowner relationships
        </p>
        
        <div style={{ 
          marginTop: '1rem', 
          display: 'flex', 
          gap: '0.5rem', 
          alignItems: 'flex-start',
          flexWrap: 'wrap'
        }}>
          <input
            type="text"
            placeholder="GitHub organization name"
            value={orgName}
            onChange={(e) => setOrgName(e.target.value)}
            style={{
              padding: '0.5rem',
              border: '1px solid #d1d5db',
              borderRadius: '0.375rem',
              fontSize: '0.875rem',
              width: '200px'
            }}
          />
          <input
            type="password"
            placeholder="GitHub token"
            value={githubToken}
            onChange={(e) => setGithubToken(e.target.value)}
            style={{
              padding: '0.5rem',
              border: '1px solid #d1d5db',
              borderRadius: '0.375rem',
              fontSize: '0.875rem',
              width: '200px'
            }}
          />
          <button
            onClick={scanOrganization}
            disabled={loading}
            style={{
              padding: '0.5rem 1rem',
              background: loading ? '#9ca3af' : '#3b82f6',
              color: 'white',
              border: 'none',
              borderRadius: '0.375rem',
              fontSize: '0.875rem',
              cursor: loading ? 'not-allowed' : 'pointer',
              fontWeight: '500'
            }}
          >
            {loading ? 'Scanning...' : 'Scan Organization'}
          </button>
          <button
            onClick={loadGraphData}
            disabled={loading || !orgName}
            style={{
              padding: '0.5rem 1rem',
              background: loading || !orgName ? '#9ca3af' : '#10b981',
              color: 'white',
              border: 'none',
              borderRadius: '0.375rem',
              fontSize: '0.875rem',
              cursor: loading || !orgName ? 'not-allowed' : 'pointer',
              fontWeight: '500'
            }}
          >
            {loading ? 'Loading...' : 'Reload Graph'}
          </button>
        </div>
        
        {error && (
          <div style={{
            marginTop: '0.5rem',
            padding: '0.5rem',
            background: '#fee2e2',
            border: '1px solid #fca5a5',
            borderRadius: '0.375rem',
            color: '#dc2626',
            fontSize: '0.875rem'
          }}>
            {error}
          </div>
        )}
        
        {stats && (
          <div style={{
            marginTop: '0.5rem',
            padding: '0.5rem',
            background: '#f3f4f6',
            borderRadius: '0.375rem',
            fontSize: '0.875rem',
            display: 'flex',
            gap: '1rem',
            flexWrap: 'wrap'
          }}>
            <span>Repos: {stats.total_repositories || 0}</span>
            <span>Teams: {stats.total_teams || 0}</span>
            <span>Users: {stats.total_users || 0}</span>
            <span>Codeowners: {stats.total_codeowners || 0}</span>
            <span>Coverage: {stats.codeowner_coverage || '0%'}</span>
          </div>
        )}
      </div>
      
      <div style={{ marginTop: '160px', height: 'calc(100vh - 160px)' }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          nodeTypes={nodeTypes}
          edgeTypes={edgeTypes}
          fitView
          attributionPosition="top-right"
          onInit={(instance) => {
            // @ts-expect-error reactFlowInstance is added to window
            window.reactFlowInstance = instance;
          }}
        >
          <Controls />
          <MiniMap 
            nodeColor={nodeColor}
            style={{
              background: 'rgba(255, 255, 255, 0.8)',
              border: '1px solid #e5e7eb',
              borderRadius: '8px',
            }}
          />
          <Background 
            variant="dots" 
            gap={20} 
            size={1} 
            color="#e5e7eb"
          />
        </ReactFlow>
      </div>
    </div>
  );
}

export default App;
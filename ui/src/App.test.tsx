import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import axios from 'axios';
import { jest, describe, test, beforeEach, expect } from '@jest/globals';
import App from './App';

// Mock axios
jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

// Mock React Flow
jest.mock('@xyflow/react', () => ({
  ReactFlow: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="react-flow">{children}</div>
  ),
  Background: () => <div data-testid="background" />,
  Controls: () => <div data-testid="controls" />,
  MiniMap: () => <div data-testid="minimap" />,
  useNodesState: () => [[], jest.fn(), jest.fn()],
  useEdgesState: () => [[], jest.fn(), jest.fn()],
  addEdge: jest.fn(),
}));

// Mock localStorage
const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
};

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
});

describe('App Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    localStorageMock.getItem.mockReturnValue(null);
  });

  test('renders main components', () => {
    render(<App />);
    
    expect(screen.getByText('GitHub Codeowner Visualization')).toBeInTheDocument();
    expect(screen.getByText('Interactive graph showing repositories, teams, users, and their codeowner relationships')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('GitHub organization name')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('GitHub token')).toBeInTheDocument();
    expect(screen.getByText('Scan Organization')).toBeInTheDocument();
    expect(screen.getByText('Reload Graph')).toBeInTheDocument();
    expect(screen.getByTestId('react-flow')).toBeInTheDocument();
  });

  test('loads saved GitHub token from localStorage', () => {
    const savedToken = 'ghp_test_token';
    localStorageMock.getItem.mockReturnValue(savedToken);
    
    render(<App />);
    
    expect(localStorageMock.getItem).toHaveBeenCalledWith('github_token');
    expect(screen.getByDisplayValue(savedToken)).toBeInTheDocument();
  });

  test('updates organization name input', () => {
    render(<App />);
    
    const orgInput = screen.getByPlaceholderText('GitHub organization name');
    fireEvent.change(orgInput, { target: { value: 'test-org' } });
    
    expect(orgInput).toHaveValue('test-org');
  });

  test('updates GitHub token input', () => {
    render(<App />);
    
    const tokenInput = screen.getByPlaceholderText('GitHub token');
    fireEvent.change(tokenInput, { target: { value: 'test-token' } });
    
    expect(tokenInput).toHaveValue('test-token');
  });

  test('shows error when scanning without organization name or token', async () => {
    render(<App />);
    
    const scanButton = screen.getByText('Scan Organization');
    fireEvent.click(scanButton);
    
    await waitFor(() => {
      expect(screen.getByText('Please enter organization name and GitHub token')).toBeInTheDocument();
    });
  });

  test('successful organization scan', async () => {
    const mockScanResponse = { data: { success: true } };
    const mockGraphResponse = {
      data: {
        nodes: [{ id: '1', type: 'repository', data: { name: 'test-repo' } }],
        edges: [{ id: 'e1', source: '1', target: '2' }]
      }
    };
    const mockStatsResponse = {
      data: {
        total_repositories: 1,
        total_teams: 0,
        total_users: 0,
        total_codeowners: 0,
        codeowner_coverage: '0%'
      }
    };

    mockedAxios.post.mockResolvedValueOnce(mockScanResponse);
    mockedAxios.get.mockResolvedValueOnce(mockGraphResponse);
    mockedAxios.get.mockResolvedValueOnce(mockStatsResponse);

    render(<App />);
    
    // Fill in form
    fireEvent.change(screen.getByPlaceholderText('GitHub organization name'), {
      target: { value: 'test-org' }
    });
    fireEvent.change(screen.getByPlaceholderText('GitHub token'), {
      target: { value: 'test-token' }
    });
    
    // Click scan
    fireEvent.click(screen.getByText('Scan Organization'));
    
    // Check loading state
    expect(screen.getByText('Scanning...')).toBeInTheDocument();
    
    await waitFor(() => {
      expect(mockedAxios.post).toHaveBeenCalledWith(
        'http://localhost:8081/api/scan/test-org',
        {},
        { headers: { 'X-GitHub-Token': 'test-token' } }
      );
      expect(localStorageMock.setItem).toHaveBeenCalledWith('github_token', 'test-token');
      expect(screen.getByText('Repos: 1')).toBeInTheDocument();
    });
  });

  test('handles scan error', async () => {
    const errorMessage = 'API Error';
    mockedAxios.post.mockRejectedValueOnce(new Error(errorMessage));

    render(<App />);
    
    fireEvent.change(screen.getByPlaceholderText('GitHub organization name'), {
      target: { value: 'test-org' }
    });
    fireEvent.change(screen.getByPlaceholderText('GitHub token'), {
      target: { value: 'test-token' }
    });
    
    fireEvent.click(screen.getByText('Scan Organization'));
    
    await waitFor(() => {
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });
  });

  test('reload graph button is disabled without organization name', () => {
    render(<App />);
    
    const reloadButton = screen.getByText('Reload Graph');
    expect(reloadButton).toBeDisabled();
  });

  test('reload graph button works with organization name', async () => {
    const mockGraphResponse = {
      data: {
        nodes: [],
        edges: []
      }
    };
    const mockStatsResponse = { data: {} };

    mockedAxios.get.mockResolvedValueOnce(mockGraphResponse);
    mockedAxios.get.mockResolvedValueOnce(mockStatsResponse);

    render(<App />);
    
    fireEvent.change(screen.getByPlaceholderText('GitHub organization name'), {
      target: { value: 'test-org' }
    });
    
    const reloadButton = screen.getByText('Reload Graph');
    expect(reloadButton).not.toBeDisabled();
    
    fireEvent.click(reloadButton);
    
    await waitFor(() => {
      expect(mockedAxios.get).toHaveBeenCalledWith('http://localhost:8081/api/graph/test-org');
    });
  });

  test('handles 404 error when loading graph data', async () => {
    const error = {
      response: { status: 404 }
    };
    mockedAxios.get.mockRejectedValueOnce(error);

    render(<App />);
    
    fireEvent.change(screen.getByPlaceholderText('GitHub organization name'), {
      target: { value: 'test-org' }
    });
    
    fireEvent.click(screen.getByText('Reload Graph'));
    
    await waitFor(() => {
      expect(screen.getByText('Organization not found. Please scan it first.')).toBeInTheDocument();
    });
  });

  test('displays stats when available', async () => {
    const mockGraphResponse = { data: { nodes: [], edges: [] } };
    const mockStatsResponse = {
      data: {
        total_repositories: 5,
        total_teams: 3,
        total_users: 10,
        total_codeowners: 8,
        codeowner_coverage: '80%'
      }
    };

    mockedAxios.get.mockResolvedValueOnce(mockGraphResponse);
    mockedAxios.get.mockResolvedValueOnce(mockStatsResponse);

    render(<App />);
    
    fireEvent.change(screen.getByPlaceholderText('GitHub organization name'), {
      target: { value: 'test-org' }
    });
    
    fireEvent.click(screen.getByText('Reload Graph'));
    
    await waitFor(() => {
      expect(screen.getByText('Repos: 5')).toBeInTheDocument();
      expect(screen.getByText('Teams: 3')).toBeInTheDocument();
      expect(screen.getByText('Users: 10')).toBeInTheDocument();
      expect(screen.getByText('Codeowners: 8')).toBeInTheDocument();
      expect(screen.getByText('Coverage: 80%')).toBeInTheDocument();
    });
  });
});
package main

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock GraphService for testing
type MockGraphService struct {
	connectErr    error
	healthErr     error
	closeErr      error
	createNodeErr error
	createRelErr  error
	queryErr      error
	connected     bool
	closed        bool
	nodeCallCount int
}

func (m *MockGraphService) Connect(ctx context.Context) error {
	m.connected = true
	return m.connectErr
}

func (m *MockGraphService) Close(ctx context.Context) error {
	m.closed = true
	return m.closeErr
}

func (m *MockGraphService) Health(ctx context.Context) error {
	return m.healthErr
}

func (m *MockGraphService) CreateNode(ctx context.Context, label string, properties map[string]interface{}) (*Node, error) {
	m.nodeCallCount++
	if m.createNodeErr != nil {
		return nil, m.createNodeErr
	}
	return &Node{
		ID:         "test-node-123",
		Labels:     []string{label},
		Properties: properties,
	}, nil
}

func (m *MockGraphService) GetNode(ctx context.Context, id string) (*Node, error) {
	return &Node{ID: id, Labels: []string{"Test"}, Properties: map[string]interface{}{}}, nil
}

func (m *MockGraphService) UpdateNode(ctx context.Context, id string, properties map[string]interface{}) error {
	return nil
}

func (m *MockGraphService) DeleteNode(ctx context.Context, id string) error {
	return nil
}

func (m *MockGraphService) CreateRelationship(ctx context.Context, fromID, toID, relType string, properties map[string]interface{}) (*Relationship, error) {
	if m.createRelErr != nil {
		return nil, m.createRelErr
	}
	return &Relationship{
		ID:         "test-rel-123",
		Type:       relType,
		FromID:     fromID,
		ToID:       toID,
		Properties: properties,
	}, nil
}

func (m *MockGraphService) GetRelationship(ctx context.Context, id string) (*Relationship, error) {
	return &Relationship{ID: id, Type: "TEST", FromID: "1", ToID: "2", Properties: map[string]interface{}{}}, nil
}

func (m *MockGraphService) DeleteRelationship(ctx context.Context, id string) error {
	return nil
}

func (m *MockGraphService) ExecuteQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	return []map[string]interface{}{
		{"name": "test-node", "description": "test description"},
	}, nil
}

func (m *MockGraphService) ExecuteReadQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	return []map[string]interface{}{
		{"name": "demo-node", "description": "A demonstration node created by overseer"},
		{"name": "demo-node-2", "description": "Another demonstration node"},
	}, nil
}

func (m *MockGraphService) ExecuteWriteQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	return []map[string]interface{}{}, nil
}

func (m *MockGraphService) ExecuteBatch(ctx context.Context, operations []BatchOperation) error {
	return nil
}

func (m *MockGraphService) ClearAll(ctx context.Context) error {
	return nil
}

// SecondCallFailMockGraphService is a mock that fails on the second CreateNode call
type SecondCallFailMockGraphService struct {
	callCount int
}

func (m *SecondCallFailMockGraphService) Connect(ctx context.Context) error {
	return nil
}

func (m *SecondCallFailMockGraphService) Close(ctx context.Context) error {
	return nil
}

func (m *SecondCallFailMockGraphService) Health(ctx context.Context) error {
	return nil
}

func (m *SecondCallFailMockGraphService) CreateNode(ctx context.Context, label string, properties map[string]interface{}) (*Node, error) {
	m.callCount++
	if m.callCount == 1 {
		return &Node{
			ID:         "test-node-1",
			Labels:     []string{label},
			Properties: properties,
		}, nil
	}
	return nil, errors.New("second node creation failed")
}

func (m *SecondCallFailMockGraphService) GetNode(ctx context.Context, id string) (*Node, error) {
	return &Node{ID: id, Labels: []string{"Test"}, Properties: map[string]interface{}{}}, nil
}

func (m *SecondCallFailMockGraphService) UpdateNode(ctx context.Context, id string, properties map[string]interface{}) error {
	return nil
}

func (m *SecondCallFailMockGraphService) DeleteNode(ctx context.Context, id string) error {
	return nil
}

func (m *SecondCallFailMockGraphService) CreateRelationship(ctx context.Context, fromID, toID, relType string, properties map[string]interface{}) (*Relationship, error) {
	return &Relationship{
		ID:         "test-rel-123",
		Type:       relType,
		FromID:     fromID,
		ToID:       toID,
		Properties: properties,
	}, nil
}

func (m *SecondCallFailMockGraphService) GetRelationship(ctx context.Context, id string) (*Relationship, error) {
	return &Relationship{ID: id, Type: "TEST", FromID: "1", ToID: "2", Properties: map[string]interface{}{}}, nil
}

func (m *SecondCallFailMockGraphService) DeleteRelationship(ctx context.Context, id string) error {
	return nil
}

func (m *SecondCallFailMockGraphService) ExecuteQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"name": "test-node", "description": "test description"},
	}, nil
}

func (m *SecondCallFailMockGraphService) ExecuteReadQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"name": "demo-node", "description": "A demonstration node created by overseer"},
		{"name": "demo-node-2", "description": "Another demonstration node"},
	}, nil
}

func (m *SecondCallFailMockGraphService) ExecuteWriteQuery(ctx context.Context, query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (m *SecondCallFailMockGraphService) ExecuteBatch(ctx context.Context, operations []BatchOperation) error {
	return nil
}

func (m *SecondCallFailMockGraphService) ClearAll(ctx context.Context) error {
	return nil
}

func TestInitializeConfiguration(t *testing.T) {
	t.Parallel()
	
	// Test successful configuration
	t.Run("successful configuration", func(t *testing.T) {
		config, err := initializeConfiguration()
		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.NotEmpty(t, config.GraphDB.Provider)
	})
	
	// Test with invalid environment to trigger validation error
	t.Run("invalid configuration", func(t *testing.T) {
		// Set an invalid provider
		_ = os.Setenv("GRAPH_DB_PROVIDER", "invalid")
		defer func() { _ = os.Unsetenv("GRAPH_DB_PROVIDER") }()
		
		config, err := initializeConfiguration()
		require.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "invalid configuration")
	})
}

func TestSetupGraphService(t *testing.T) {
	t.Parallel()
	
	// Test connection failure with invalid host
	t.Run("connection failure", func(t *testing.T) {
		ctx := context.Background()
		config := &Config{
			GraphDB: GraphServiceConfig{
				Provider: "neo4j",
				Neo4j: struct {
					URI      string `json:"uri"`
					Username string `json:"username"`
					Password string `json:"password"`
				}{
					URI:      "bolt://invalid-host:7687",
					Username: "neo4j",
					Password: "password",
				},
			},
		}
		
		// This will fail because we use invalid host
		service, err := setupGraphService(ctx, config)
		require.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "failed to connect to graph database")
	})
	
	// Test invalid provider
	t.Run("invalid provider", func(t *testing.T) {
		ctx := context.Background()
		config := &Config{
			GraphDB: GraphServiceConfig{
				Provider: "invalid",
			},
		}
		
		service, err := setupGraphService(ctx, config)
		require.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "failed to create graph service")
	})
}

func TestWaitForShutdown(t *testing.T) {
	t.Parallel()
	
	// Test context cancellation
	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		
		// Cancel after a short delay
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()
		
		err := waitForShutdown(ctx)
		assert.NoError(t, err)
	})
	
	// Test with timeout
	t.Run("context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		
		err := waitForShutdown(ctx)
		assert.NoError(t, err)
	})
}

func TestRunDemoOperations(t *testing.T) {
	t.Parallel()
	
	// Test successful demo operations
	t.Run("successful demo operations", func(t *testing.T) {
		ctx := context.Background()
		mockService := &MockGraphService{}
		
		// This should not panic or return error
		runDemoOperations(ctx, mockService)
		
		// Verify the mock was called (demonstrateGraphOperations should have called CreateNode)
		assert.True(t, mockService.nodeCallCount > 0)
	})
	
	// Test demo operations with error
	t.Run("demo operations with error", func(t *testing.T) {
		ctx := context.Background()
		mockService := &MockGraphService{
			createNodeErr: errors.New("create node failed"),
		}
		
		// This should not panic even with error
		runDemoOperations(ctx, mockService)
	})
}

func TestDemonstrateGraphOperations(t *testing.T) {
	t.Parallel()
	
	// Test successful operations
	t.Run("successful operations", func(t *testing.T) {
		ctx := context.Background()
		mockService := &MockGraphService{}
		
		err := demonstrateGraphOperations(ctx, mockService)
		require.NoError(t, err)
	})
	
	// Test first node creation failure
	t.Run("first node creation failure", func(t *testing.T) {
		ctx := context.Background()
		mockService := &MockGraphService{
			createNodeErr: errors.New("first node creation failed"),
		}
		
		err := demonstrateGraphOperations(ctx, mockService)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create test node")
	})
	
	// Test second node creation failure - using a custom mock
	t.Run("second node creation failure", func(t *testing.T) {
		ctx := context.Background()
		
		// Create a mock that simulates success for first call, failure for second
		mockService := &SecondCallFailMockGraphService{}
		
		err := demonstrateGraphOperations(ctx, mockService)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create second test node")
	})
	
	// Test relationship creation failure
	t.Run("relationship creation failure", func(t *testing.T) {
		ctx := context.Background()
		mockService := &MockGraphService{
			createRelErr: errors.New("relationship creation failed"),
		}
		
		err := demonstrateGraphOperations(ctx, mockService)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create relationship")
	})
	
	// Test query execution failure
	t.Run("query execution failure", func(t *testing.T) {
		ctx := context.Background()
		mockService := &MockGraphService{
			queryErr: errors.New("query execution failed"),
		}
		
		err := demonstrateGraphOperations(ctx, mockService)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to execute query")
	})
}

func TestRunOverseer(t *testing.T) {
	t.Parallel()
	
	// Test with invalid configuration
	t.Run("invalid configuration", func(t *testing.T) {
		ctx := context.Background()
		
		// Set invalid provider
		_ = os.Setenv("GRAPH_DB_PROVIDER", "invalid")
		defer func() { _ = os.Unsetenv("GRAPH_DB_PROVIDER") }()
		
		err := runOverseer(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid configuration")
	})
	
	// Test with valid configuration but connection failure
	t.Run("connection failure", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Set up environment for Neo4j with invalid URI (will fail to connect)
		_ = os.Setenv("GRAPH_DB_PROVIDER", "neo4j")
		_ = os.Setenv("NEO4J_URI", "bolt://nonexistent-host:7687")
		defer func() { 
			_ = os.Unsetenv("GRAPH_DB_PROVIDER")
			_ = os.Unsetenv("NEO4J_URI")
		}()
		
		err := runOverseer(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect to graph database")
	})
}

func TestMainFunctionBehavior(t *testing.T) {
	t.Parallel()
	
	// Test context cancellation behavior
	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		
		// Test that context cancellation works
		cancel()
		
		select {
		case <-ctx.Done():
			// Context was properly cancelled
			assert.True(t, true)
		default:
			assert.Fail(t, "Context should have been cancelled")
		}
	})
}

func TestEnvironmentDependentBehavior(t *testing.T) {
	t.Parallel()
	
	// Test development mode triggers demo operations
	t.Run("development mode", func(t *testing.T) {
		ctx := context.Background()
		mockService := &MockGraphService{}
		
		// Mock config for development
		config := &Config{
			Environment: "development",
		}
		
		// Test that development mode would trigger demo operations
		if checkIsDevelopment(*config) {
			runDemoOperations(ctx, mockService)
		}
		
		assert.True(t, checkIsDevelopment(*config))
	})
	
	// Test non-development mode skips demo operations
	t.Run("production mode", func(t *testing.T) {
		config := &Config{
			Environment: "production",
		}
		
		// Test that production mode would not trigger demo operations
		assert.False(t, checkIsDevelopment(*config))
		assert.True(t, checkIsProduction(*config))
	})
}

func TestConfigurationErrorHandling(t *testing.T) {
	t.Parallel()
	
	// Test error wrapping in initializeConfiguration
	t.Run("configuration loading error", func(t *testing.T) {
		// Set an environment that would cause LoadConfig to fail
		// This is tricky since LoadConfig doesn't really fail in current implementation
		// But we can test the error path by setting invalid provider
		_ = os.Setenv("GRAPH_DB_PROVIDER", "invalid")
		defer func() { _ = os.Unsetenv("GRAPH_DB_PROVIDER") }()
		
		config, err := initializeConfiguration()
		require.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "invalid configuration")
	})
}

func TestGraphServiceErrorHandling(t *testing.T) {
	t.Parallel()
	
	// Test health check failure
	t.Run("health check failure", func(t *testing.T) {
		ctx := context.Background()
		config := &Config{
			GraphDB: GraphServiceConfig{
				Provider: "neo4j",
				Neo4j: struct {
					URI      string `json:"uri"`
					Username string `json:"username"`
					Password string `json:"password"`
				}{
					URI:      "bolt://nonexistent-host:7687",
					Username: "neo4j",
					Password: "password",
				},
			},
		}
		
		// This will fail because we don't have Neo4j running
		service, err := setupGraphService(ctx, config)
		require.Error(t, err)
		assert.Nil(t, service)
		
		// The error should be about connection failure
		assert.Contains(t, err.Error(), "failed to connect to graph database")
	})
}

func TestDeferredCleanup(t *testing.T) {
	t.Parallel()
	
	// Test that deferred cleanup works
	t.Run("deferred cleanup", func(t *testing.T) {
		mockService := &MockGraphService{}
		
		// Simulate the deferred cleanup
		func() {
			ctx := context.Background()
			defer func() {
				if err := mockService.Close(ctx); err != nil {
					// This would normally log the error
					assert.NoError(t, err)
				}
			}()
		}()
		
		assert.True(t, mockService.closed)
	})
	
	// Test deferred cleanup with error
	t.Run("deferred cleanup with error", func(t *testing.T) {
		mockService := &MockGraphService{
			closeErr: errors.New("close error"),
		}
		
		// Simulate the deferred cleanup
		func() {
			ctx := context.Background()
			defer func() {
				if err := mockService.Close(ctx); err != nil {
					// This would normally log the error
					assert.Error(t, err)
				}
			}()
		}()
		
		assert.True(t, mockService.closed)
	})
}
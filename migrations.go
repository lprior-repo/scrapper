package main

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"
)

// Migration represents a database migration
type Migration struct {
	Version     int       `json:"version"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpQuery     string    `json:"up_query"`
	DownQuery   string    `json:"down_query"`
	AppliedAt   time.Time `json:"applied_at,omitempty"`
}

// MigrationState represents the current migration state
type MigrationState struct {
	CurrentVersion int         `json:"current_version"`
	Migrations     []Migration `json:"migrations"`
}

// getMigrations returns all available migrations (Pure Core)
func getMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Name:        "initial_schema",
			Description: "Create initial schema with indexes and constraints",
			UpQuery: `
				// Create constraints
				CREATE CONSTRAINT constraint_organization_name_unique IF NOT EXISTS FOR (n:Organization) REQUIRE n.name IS UNIQUE;
				CREATE CONSTRAINT constraint_repository_full_name_unique IF NOT EXISTS FOR (n:Repository) REQUIRE n.full_name IS UNIQUE;
				CREATE CONSTRAINT constraint_user_username_unique IF NOT EXISTS FOR (n:User) REQUIRE n.username IS UNIQUE;
				CREATE CONSTRAINT constraint_user_email_unique IF NOT EXISTS FOR (n:User) REQUIRE n.email IS UNIQUE;
				CREATE CONSTRAINT constraint_team_name_org_unique IF NOT EXISTS FOR (n:Team) REQUIRE (n.name, n.organization) IS UNIQUE;
				
				// Create indexes
				CREATE INDEX idx_organization_name IF NOT EXISTS FOR (n:Organization) ON (n.name);
				CREATE INDEX idx_repository_full_name IF NOT EXISTS FOR (n:Repository) ON (n.full_name);
				CREATE INDEX idx_repository_name IF NOT EXISTS FOR (n:Repository) ON (n.name);
				CREATE INDEX idx_user_username IF NOT EXISTS FOR (n:User) ON (n.username);
				CREATE INDEX idx_user_email IF NOT EXISTS FOR (n:User) ON (n.email);
				CREATE INDEX idx_team_name IF NOT EXISTS FOR (n:Team) ON (n.name);
				CREATE INDEX idx_org_repo_composite IF NOT EXISTS FOR (n:Repository) ON (n.organization, n.name);
				CREATE INDEX idx_codeowners_text IF NOT EXISTS FOR (n:Repository) ON (n.codeowners_content);
			`,
			DownQuery: `
				// Drop indexes
				DROP INDEX idx_organization_name IF EXISTS;
				DROP INDEX idx_repository_full_name IF EXISTS;
				DROP INDEX idx_repository_name IF EXISTS;
				DROP INDEX idx_user_username IF EXISTS;
				DROP INDEX idx_user_email IF EXISTS;
				DROP INDEX idx_team_name IF EXISTS;
				DROP INDEX idx_org_repo_composite IF EXISTS;
				DROP INDEX idx_codeowners_text IF EXISTS;
				
				// Drop constraints
				DROP CONSTRAINT constraint_organization_name_unique IF EXISTS;
				DROP CONSTRAINT constraint_repository_full_name_unique IF EXISTS;
				DROP CONSTRAINT constraint_user_username_unique IF EXISTS;
				DROP CONSTRAINT constraint_user_email_unique IF EXISTS;
				DROP CONSTRAINT constraint_team_name_org_unique IF EXISTS;
			`,
		},
		{
			Version:     2,
			Name:        "migration_tracking",
			Description: "Create migration tracking node",
			UpQuery: `
				// Create migration tracking node
				MERGE (m:Migration {id: 'system'})
				SET m.current_version = 2,
					m.last_updated = datetime(),
					m.created_at = CASE WHEN m.created_at IS NULL THEN datetime() ELSE m.created_at END;
			`,
			DownQuery: `
				// Remove migration tracking
				MATCH (m:Migration {id: 'system'})
				DELETE m;
			`,
		},
		{
			Version:     3,
			Name:        "performance_indexes",
			Description: "Add performance indexes for common queries",
			UpQuery: `
				// Add performance indexes
				CREATE INDEX idx_repository_has_codeowners IF NOT EXISTS FOR (n:Repository) ON (n.has_codeowners_file);
				CREATE INDEX idx_repository_created_at IF NOT EXISTS FOR (n:Repository) ON (n.created_at);
				CREATE INDEX idx_user_created_at IF NOT EXISTS FOR (n:User) ON (n.created_at);
				CREATE INDEX idx_team_created_at IF NOT EXISTS FOR (n:Team) ON (n.created_at);
				
				// Add relationship indexes
				CREATE INDEX idx_owns_relationship IF NOT EXISTS FOR ()-[r:OWNS]-() ON (r.created_at);
				CREATE INDEX idx_member_of_relationship IF NOT EXISTS FOR ()-[r:MEMBER_OF]-() ON (r.created_at);
				CREATE INDEX idx_has_codeowner_relationship IF NOT EXISTS FOR ()-[r:HAS_CODEOWNER]-() ON (r.pattern);
			`,
			DownQuery: `
				// Drop performance indexes
				DROP INDEX idx_repository_has_codeowners IF EXISTS;
				DROP INDEX idx_repository_created_at IF EXISTS;
				DROP INDEX idx_user_created_at IF EXISTS;
				DROP INDEX idx_team_created_at IF EXISTS;
				DROP INDEX idx_owns_relationship IF EXISTS;
				DROP INDEX idx_member_of_relationship IF EXISTS;
				DROP INDEX idx_has_codeowner_relationship IF EXISTS;
			`,
		},
	}
}

// getCurrentMigrationVersion gets the current migration version (Impure Shell)
func getCurrentMigrationVersion(ctx context.Context, conn GraphConnection) (int, error) {
	if conn.Driver == nil {
		return 0, fmt.Errorf("driver not initialized")
	}

	query := "MATCH (m:Migration {id: 'system'}) RETURN m.current_version as version"
	results, err := executeQuery(ctx, conn, query, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get migration version: %w", err)
	}

	if len(results) == 0 {
		return 0, nil // No migrations applied yet
	}

	if version, ok := results[0]["version"].(int64); ok {
		return int(version), nil
	}

	if versionStr, ok := results[0]["version"].(string); ok {
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			return 0, fmt.Errorf("invalid version format: %s", versionStr)
		}
		return version, nil
	}

	return 0, fmt.Errorf("unexpected version type: %T", results[0]["version"])
}

// setMigrationVersion sets the current migration version (Impure Shell)
func setMigrationVersion(ctx context.Context, conn GraphConnection, version int) error {
	if conn.Driver == nil {
		return fmt.Errorf("driver not initialized")
	}

	query := `
		MERGE (m:Migration {id: 'system'})
		SET m.current_version = $version,
			m.last_updated = datetime(),
			m.created_at = CASE WHEN m.created_at IS NULL THEN datetime() ELSE m.created_at END
	`

	params := map[string]interface{}{
		"version": version,
	}

	_, err := executeQuery(ctx, conn, query, params)
	if err != nil {
		return fmt.Errorf("failed to set migration version: %w", err)
	}

	return nil
}

// applyMigration applies a single migration (Impure Shell)
func applyMigration(ctx context.Context, conn GraphConnection, migration Migration) error {
	if conn.Driver == nil {
		return fmt.Errorf("driver not initialized")
	}

	// Split the migration into individual statements
	statements := splitMigrationStatements(migration.UpQuery)

	for _, statement := range statements {
		statement = trimStatement(statement)
		if statement == "" || isComment(statement) {
			continue
		}

		_, err := executeQuery(ctx, conn, statement, nil)
		if err != nil {
			return fmt.Errorf("failed to execute migration statement: %s, error: %w", statement, err)
		}
	}

	// Update migration version
	err := setMigrationVersion(ctx, conn, migration.Version)
	if err != nil {
		return fmt.Errorf("failed to update migration version: %w", err)
	}

	return nil
}

// rollbackMigration rolls back a single migration (Impure Shell)
func rollbackMigration(ctx context.Context, conn GraphConnection, migration Migration) error {
	if conn.Driver == nil {
		return fmt.Errorf("driver not initialized")
	}

	// Split the rollback into individual statements
	statements := splitMigrationStatements(migration.DownQuery)

	for _, statement := range statements {
		statement = trimStatement(statement)
		if statement == "" || isComment(statement) {
			continue
		}

		_, err := executeQuery(ctx, conn, statement, nil)
		if err != nil {
			return fmt.Errorf("failed to execute rollback statement: %s, error: %w", statement, err)
		}
	}

	// Update migration version to previous version
	err := setMigrationVersion(ctx, conn, migration.Version-1)
	if err != nil {
		return fmt.Errorf("failed to update migration version: %w", err)
	}

	return nil
}

// runMigrationsUp runs all pending migrations (Impure Shell)
func runMigrationsUp(ctx context.Context, conn GraphConnection) error {
	if conn.Driver == nil {
		return fmt.Errorf("driver not initialized")
	}

	currentVersion, err := getCurrentMigrationVersion(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	migrations := getMigrations()

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Apply pending migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			fmt.Printf("Applying migration %d: %s\n", migration.Version, migration.Name)
			err := applyMigration(ctx, conn, migration)
			if err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
			}
		}
	}

	return nil
}

// runMigrationsDown rolls back migrations to target version (Impure Shell)
func runMigrationsDown(ctx context.Context, conn GraphConnection, targetVersion int) error {
	if conn.Driver == nil {
		return fmt.Errorf("driver not initialized")
	}

	if targetVersion < 0 {
		return fmt.Errorf("target version cannot be negative")
	}

	currentVersion, err := getCurrentMigrationVersion(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if targetVersion >= currentVersion {
		return fmt.Errorf("target version %d is not less than current version %d", targetVersion, currentVersion)
	}

	migrations := getMigrations()

	// Sort migrations by version in reverse order for rollback
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version > migrations[j].Version
	})

	// Rollback migrations
	for _, migration := range migrations {
		if migration.Version > targetVersion && migration.Version <= currentVersion {
			fmt.Printf("Rolling back migration %d: %s\n", migration.Version, migration.Name)
			err := rollbackMigration(ctx, conn, migration)
			if err != nil {
				return fmt.Errorf("failed to rollback migration %d: %w", migration.Version, err)
			}
		}
	}

	return nil
}

// validateMigrations validates all migration definitions (Pure Core)
func validateMigrations() error {
	migrations := getMigrations()

	if len(migrations) == 0 {
		return fmt.Errorf("no migrations defined")
	}

	versionsSeen := make(map[int]bool)

	for _, migration := range migrations {
		// Check for duplicate versions
		if versionsSeen[migration.Version] {
			return fmt.Errorf("duplicate migration version: %d", migration.Version)
		}
		versionsSeen[migration.Version] = true

		// Check required fields
		if migration.Version <= 0 {
			return fmt.Errorf("migration version must be positive: %d", migration.Version)
		}

		if migration.Name == "" {
			return fmt.Errorf("migration name cannot be empty for version %d", migration.Version)
		}

		if migration.UpQuery == "" {
			return fmt.Errorf("migration up query cannot be empty for version %d", migration.Version)
		}

		if migration.DownQuery == "" {
			return fmt.Errorf("migration down query cannot be empty for version %d", migration.Version)
		}
	}

	return nil
}

// Helper functions for migration statement processing (Pure Core)

// splitMigrationStatements splits a migration string into individual statements
func splitMigrationStatements(migration string) []string {
	// Simple split by semicolon for now
	// In production, might need more sophisticated parsing
	statements := []string{}
	current := ""

	for _, line := range splitLines(migration) {
		line = trimStatement(line)
		if line == "" || isComment(line) {
			continue
		}

		current += line + " "

		if endsWithSemicolon(line) {
			statements = append(statements, current)
			current = ""
		}
	}

	if current != "" {
		statements = append(statements, current)
	}

	return statements
}

// splitLines splits a string into lines
func splitLines(s string) []string {
	lines := []string{}
	current := ""

	for _, char := range s {
		if char == '\n' || char == '\r' {
			if current != "" {
				lines = append(lines, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		lines = append(lines, current)
	}

	return lines
}

// trimStatement trims whitespace and removes comments
func trimStatement(statement string) string {
	// Remove leading/trailing whitespace
	result := ""
	start := 0
	end := len(statement) - 1

	// Find first non-whitespace character
	for start < len(statement) && isWhitespace(rune(statement[start])) {
		start++
	}

	// Find last non-whitespace character
	for end >= 0 && isWhitespace(rune(statement[end])) {
		end--
	}

	if start <= end {
		result = statement[start : end+1]
	}

	return result
}

// isComment checks if a line is a comment
func isComment(line string) bool {
	trimmed := trimStatement(line)
	return len(trimmed) >= 2 && trimmed[:2] == "//"
}

// endsWithSemicolon checks if a line ends with a semicolon
func endsWithSemicolon(line string) bool {
	trimmed := trimStatement(line)
	return len(trimmed) > 0 && trimmed[len(trimmed)-1] == ';'
}

// isWhitespace checks if a character is whitespace
func isWhitespace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

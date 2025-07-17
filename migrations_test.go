package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMigrations(t *testing.T) {
	migrations := getMigrations()

	assert.NotEmpty(t, migrations, "Should have migrations defined")

	// Test that first migration is version 1
	assert.Equal(t, 1, migrations[0].Version, "First migration should be version 1")

	// Test that all migrations have required fields
	for _, migration := range migrations {
		assert.Positive(t, migration.Version, "Migration version should be positive")
		assert.NotEmpty(t, migration.Name, "Migration name should not be empty")
		assert.NotEmpty(t, migration.Description, "Migration description should not be empty")
		assert.NotEmpty(t, migration.UpQuery, "Migration up query should not be empty")
		assert.NotEmpty(t, migration.DownQuery, "Migration down query should not be empty")
	}
}

func TestValidateMigrations(t *testing.T) {
	t.Run("valid migrations", func(t *testing.T) {
		err := validateMigrations()
		assert.NoError(t, err, "Valid migrations should pass validation")
	})
}

func TestSplitMigrationStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "single statement",
			input: "CREATE INDEX test;",
			expected: []string{
				"CREATE INDEX test; ",
			},
		},
		{
			name: "multiple statements",
			input: `CREATE INDEX test1;
					CREATE INDEX test2;`,
			expected: []string{
				"CREATE INDEX test1; ",
				"CREATE INDEX test2; ",
			},
		},
		{
			name: "statements with comments",
			input: `// This is a comment
					CREATE INDEX test1;
					// Another comment
					CREATE INDEX test2;`,
			expected: []string{
				"CREATE INDEX test1; ",
				"CREATE INDEX test2; ",
			},
		},
		{
			name:  "statement without semicolon",
			input: "CREATE INDEX test",
			expected: []string{
				"CREATE INDEX test ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitMigrationStatements(tt.input)
			assert.Len(t, result, len(tt.expected), "Should have correct number of statements")

			for i, expected := range tt.expected {
				assert.Equal(t, expected, result[i], "Statement %d should match", i)
			}
		})
	}
}

func TestTrimStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no whitespace",
			input:    "CREATE INDEX",
			expected: "CREATE INDEX",
		},
		{
			name:     "leading whitespace",
			input:    "   CREATE INDEX",
			expected: "CREATE INDEX",
		},
		{
			name:     "trailing whitespace",
			input:    "CREATE INDEX   ",
			expected: "CREATE INDEX",
		},
		{
			name:     "both leading and trailing",
			input:    "   CREATE INDEX   ",
			expected: "CREATE INDEX",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimStatement(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsComment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "comment line",
			input:    "// This is a comment",
			expected: true,
		},
		{
			name:     "comment with whitespace",
			input:    "   // This is a comment",
			expected: true,
		},
		{
			name:     "not a comment",
			input:    "CREATE INDEX",
			expected: false,
		},
		{
			name:     "empty line",
			input:    "",
			expected: false,
		},
		{
			name:     "single slash",
			input:    "/",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isComment(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEndsWithSemicolon(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "ends with semicolon",
			input:    "CREATE INDEX test;",
			expected: true,
		},
		{
			name:     "ends with semicolon and whitespace",
			input:    "CREATE INDEX test;   ",
			expected: true,
		},
		{
			name:     "no semicolon",
			input:    "CREATE INDEX test",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "semicolon in middle",
			input:    "CREATE; INDEX test",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := endsWithSemicolon(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"space", ' ', true},
		{"tab", '\t', true},
		{"newline", '\n', true},
		{"carriage return", '\r', true},
		{"letter", 'a', false},
		{"number", '1', false},
		{"symbol", '!', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWhitespace(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single line",
			input:    "single line",
			expected: []string{"single line"},
		},
		{
			name:     "multiple lines with newline",
			input:    "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "multiple lines with carriage return",
			input:    "line1\rline2\rline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only newlines",
			input:    "\n\n\n",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitLines(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMigrationValidation(t *testing.T) {
	t.Run("migration with all required fields", func(t *testing.T) {
		migration := Migration{
			Version:     1,
			Name:        "test_migration",
			Description: "Test migration description",
			UpQuery:     "CREATE INDEX test;",
			DownQuery:   "DROP INDEX test;",
		}

		assert.Positive(t, migration.Version)
		assert.NotEmpty(t, migration.Name)
		assert.NotEmpty(t, migration.Description)
		assert.NotEmpty(t, migration.UpQuery)
		assert.NotEmpty(t, migration.DownQuery)
	})
}

func TestMigrationStateValidation(t *testing.T) {
	t.Run("valid migration state", func(t *testing.T) {
		state := MigrationState{
			CurrentVersion: 3,
			Migrations:     getMigrations(),
		}

		assert.GreaterOrEqual(t, state.CurrentVersion, 0)
		assert.NotEmpty(t, state.Migrations)
	})
}

package mysql

import (
	"database/sql"

	"github.com/uuzay/snippetbox/pkg/models"
)

// Wrapper struct for connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert a snippet
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// Get snippet by ID
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

// Get last <count> snippets
func (m *SnippetModel) Latest(count int) ([]*models.Snippet, error) {
	return nil, nil
}

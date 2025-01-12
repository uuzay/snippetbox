package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/uuzay/snippetbox/pkg/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// Wrapper struct for connection pool
type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(query, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	query := `SELECT id, hashed_password FROM users WHERE email = ? AND active = TRUE`

	var id int
	var hashedPassword []byte
	row := m.DB.QueryRow(query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		//TODO: This does not distinguish user found but inactive case
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

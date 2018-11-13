package uniquevalidator

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// UniqueRule to check data is unique or not in db
type UniqueRule struct {
	db       *sql.DB
	ruleName string
}

// Rule is func to register as custom rule
func (r *UniqueRule) Rule(field string, rule string, message string, value interface{}) error {
	var queryRow *sql.Row
	var total int

	query := `SELECT COUNT(*) as total FROM %s WHERE %s = ?`
	params := strings.Split(strings.TrimPrefix(rule, fmt.Sprintf("%s:", r.ruleName)), ",")

	if len(params) <= 2 {
		query = fmt.Sprintf(query, params[0], params[1])
		queryRow = r.db.QueryRow(query, value)
	} else {
		query += ` AND %s != ?`
		query = fmt.Sprintf(query, params[0], params[1], params[2])
		queryRow = r.db.QueryRow(query, value, params[3])
	}

	err := queryRow.Scan(&total)
	if err != nil {
		return err
	}

	if total > 0 {
		if message != "" {
			return errors.New(message)
		}

		return fmt.Errorf("The %s has already been taken", field)
	}

	return nil
}

// NewUniqueRule to create instance
func NewUniqueRule(db *sql.DB, ruleName string) *UniqueRule {
	return &UniqueRule{db, ruleName}
}

package postgres

import (
	"github.com/go-pg/pg/v10"
)

// SetIsolationLevel set isolation level for transaction.
// See https://www.postgresql.org/docs/current/sql-set-transaction.html for syntax.
// See https://www.postgresql.org/docs/current/transaction-iso.html fot meaning.
func SetIsolationLevel(tx *pg.Tx, level string) error {
	_, err := tx.Exec("SET TRANSACTION ISOLATION LEVEL " + level)
	return err
}

// SetIsolationLevelReadCommitted set to default PostgreSQL isolation level.
// In "select .. for update" levels "serializable" and "repeatable read" causes errors.
// See https://www.postgresql.org/docs/current/explicit-locking.html for details.
func SetIsolationLevelReadCommitted(tx *pg.Tx) error {
	return SetIsolationLevel(tx, "READ COMMITTED")
}

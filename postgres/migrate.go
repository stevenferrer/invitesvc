package postgres

import (
	"database/sql"

	// postgres driver
	_ "github.com/lib/pq"
	"github.com/lopezator/migrator"
)

// Migrate migrates the database.
func Migrate(db *sql.DB, opts ...migrator.Option) error {
	opts = append(opts, migrations)
	m, err := migrator.New(opts...)
	if err != nil {
		return err
	}

	return m.Migrate(db)
}

// MustMigrate migrates the database and panics if an error occurs.
func MustMigrate(db *sql.DB, opts ...migrator.Option) {
	err := Migrate(db, opts...)
	if err != nil {
		panic(err)
	}
}

// migrations is the list of migrations
var migrations = migrator.Migrations(
	&migrator.Migration{
		Name: "Create tokens table",
		Func: func(tx *sql.Tx) error {
			stmnt := `CREATE TABLE IF NOT EXISTS "tokens" (
				token varchar(12) PRIMARY KEY,
				disabled boolean NOT NULL DEFAULT FALSE,
				redeemed_at timestamp,
				updated_at timestamp,
				created_at timestamp NOT NULL DEFAULT NOW()
			)`
			if _, err := tx.Exec(stmnt); err != nil {
				return err
			}

			return nil
		},
	},
	&migrator.Migration{
		Name: "Create auth_keys table",
		Func: func(tx *sql.Tx) error {
			stmnt := `CREATE TABLE IF NOT EXISTS "auth_keys" (
				auth_key varchar(32) PRIMARY KEY,
				created_at timestamp NOT NULL DEFAULT NOW()
			)`
			if _, err := tx.Exec(stmnt); err != nil {
				return err
			}

			return nil
		},
	},
	// Add new migration
)

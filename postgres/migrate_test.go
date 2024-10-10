package postgres_test

import (
	"testing"

	"github.com/stevenferrer/invitesvc/postgres"
	"github.com/stevenferrer/invitesvc/postgres/txdb"
)

func TestMigrate(t *testing.T) {
	db := txdb.MustOpen()
	defer db.Close()
	postgres.MustMigrate(db)
}

package inmem

import (
	"fmt"

	"github.com/stevenferrer/invitesvc/authn"
	"github.com/stevenferrer/invitesvc/token"

	"github.com/hashicorp/go-memdb"
)

const (
	tokensTable = "tokens"
	authsTable  = "authns"
)

// Schema returns the memdb schema
func Schema() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tokensTable: {
				Name: tokensTable,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: tokenIDIndexer{},
					},
				},
			},
			authsTable: {
				Name: authsTable,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: authKeyIndexer{},
					},
				},
			},
		},
	}
}

// tokenIDIndexer implements memdb.Indexer and memdb.SingleIndexer
type tokenIDIndexer struct{}

func (tokenIDIndexer) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of args %d, expected 1", len(args))
	}

	id, ok := args[0].(token.ID)
	if !ok {
		return nil, fmt.Errorf("wrong type for arg %T, expected string", args[0])
	}

	return append([]byte(id), 0), nil
}

func (tokenIDIndexer) FromObject(raw interface{}) (bool, []byte, error) {
	t, ok := raw.(*token.Token)
	if !ok {
		return false, nil, fmt.Errorf("wrong type for arg %T, expected MyStruct", raw)
	}

	if t.ID == token.NilID {
		return false, nil, nil
	}

	return true, append([]byte(t.ID), 0), nil
}

// authKeyIndexer implements memdb.Indexer and memdb.SingleIndexer
type authKeyIndexer struct{}

func (authKeyIndexer) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of args %d, expected 1", len(args))
	}

	authKey, ok := args[0].(authn.AuthKey)
	if !ok {
		return nil, fmt.Errorf("wrong type for arg %T, expected string", args[0])
	}

	return append([]byte(authKey), 0), nil
}

func (authKeyIndexer) FromObject(raw interface{}) (bool, []byte, error) {
	a, ok := raw.(*authn.Auth)
	if !ok {
		return false, nil, fmt.Errorf("wrong type for arg %T, expected MyStruct", raw)
	}

	if a.Auth == authn.AuthKey("") {
		return false, nil, nil
	}

	return true, append([]byte(a.Auth), 0), nil
}

package blindindexstore

import (
	"errors"

	"github.com/gouniverse/sb"
)

// NewStore creates a new entity store
func NewStore(opts NewStoreOptions) (*Store, error) {
	store := &Store{
		tableName:          opts.TableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
	}

	if store.tableName == "" {
		return nil, errors.New("vault store: vaultTableName is required")
	}

	if store.db == nil {
		return nil, errors.New("vault store: DB is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = sb.DatabaseDriverName(store.db)
	}

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}

package blindindexstore

import (
	"errors"

	"github.com/gouniverse/sb"
)

// NewStore creates a new entity store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	store := &storeImplementation{
		tableName:          opts.TableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
		transformer:        opts.Transformer,
	}

	if store.tableName == "" {
		return nil, errors.New("blind index store: TableName is required")
	}

	if store.db == nil {
		return nil, errors.New("blind index store: DB is required")
	}

	if store.transformer == nil {
		return nil, errors.New("blind index store: Transformer is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = sb.DatabaseDriverName(store.db)
	}

	if store.automigrateEnabled {
		err := store.AutoMigrate()

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}

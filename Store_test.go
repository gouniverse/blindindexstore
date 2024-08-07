package blindindexstore

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/gouniverse/sb"
	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func TestWithAutoMigrate(t *testing.T) {
	db := initDB("test_store_with_automigrate.db")

	storeAutomigrateFalse, errAutomigrateFalse := NewStore(NewStoreOptions{
		TableName:          "test_blindindex_with_automigrate_false",
		DB:                 db,
		AutomigrateEnabled: false,
	})

	if errAutomigrateFalse != nil {
		t.Fatalf("automigrateEnabled: Expected [err] to be nill received [%v]", errAutomigrateFalse.Error())
	}

	if storeAutomigrateFalse.automigrateEnabled != false {
		t.Fatalf("automigrateEnabled: Expected [false] received [%v]", storeAutomigrateFalse.automigrateEnabled)
	}

	storeAutomigrateTrue, errAutomigrateTrue := NewStore(NewStoreOptions{
		TableName:          "test_blockindex_with_automigrate_true",
		DB:                 db,
		AutomigrateEnabled: true,
	})

	if errAutomigrateTrue != nil {
		t.Fatalf("automigrateEnabled: Expected [err] to be nill received [%v]", errAutomigrateTrue.Error())
	}

	if storeAutomigrateTrue.automigrateEnabled != true {
		t.Fatalf("automigrateEnabled: Expected [true] received [%v]", storeAutomigrateTrue.automigrateEnabled)
	}
}

func TestStoreSearchValueCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreSearchValueFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_find_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	value := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(value)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	valueFound, errFind := store.SearchValueFindByID(value.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if valueFound == nil {
		t.Fatal("Search value MUST NOT be nil")
		return
	}

	if valueFound.SourceReferenceID() != "RefId01" {
		t.Fatal("Search value reference ID MUST BE 'RefId01', found: ", valueFound.SourceReferenceID())
		return
	}

	if valueFound.SearchValue() != "SearchValue01" {
		t.Fatal("Search value MUST BE 'SearchValue01', found: ", valueFound.SearchValue())
	}

	if !strings.Contains(valueFound.DeletedAt(), sb.MAX_DATE) {
		t.Fatal("Search value MUST NOT be soft deleted", valueFound.DeletedAt())
		return
	}
}

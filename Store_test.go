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

func Test_Store_WithAutoMigrate(t *testing.T) {
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

func Test_Store_SearchValueCreate(t *testing.T) {
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

func Test_Store_SearchValueFindByID(t *testing.T) {
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

func Test_Store_SearchValueDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	if err != nil {
		t.Fatalf("Test_Store_ValueDelete: Expected [err] to be nil received [%v]", err.Error())
	}

	value := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(value)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	errDelete := store.SearchValueDelete(value)
	if errDelete != nil {
		t.Fatalf("ValueDelete Failed: " + errDelete.Error())
	}

	valueFound, errFind := store.SearchValueFindByID(value.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if valueFound != nil {
		t.Fatal("Search value MUST be nil")
		return
	}

	valuesFound2, errFind := store.SearchValueList(SearchValueQueryOptions{
		ID:          value.ID(),
		WithDeleted: true,
	})

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(valuesFound2) > 0 {
		t.Fatal("Search values MUST be 0")
		return
	}

}

func Test_Store_SearchValueSoftDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_soft_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	if err != nil {
		t.Fatalf("Test_Store_ValueDelete: Expected [err] to be nil received [%v]", err.Error())
	}

	value := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(value)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	errDelete := store.SearchValueSoftDelete(value)
	if errDelete != nil {
		t.Fatalf("ValueDelete Failed: " + errDelete.Error())
	}

	valueFound, errFind := store.SearchValueFindByID(value.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if valueFound != nil {
		t.Fatal("Search value MUST be nil")
		return
	}

	valuesFound2, errFind := store.SearchValueList(SearchValueQueryOptions{
		ID:          value.ID(),
		WithDeleted: true,
	})

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(valuesFound2) > 1 {
		t.Fatal("Search values MUST NOT be 0")
		return
	}

}

func Test_Store_SearchEqual(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_equals",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	data := []struct {
		RefID       string
		SearchValue string
	}{
		{
			RefID:       "USER01",
			SearchValue: "test01@test.com",
		},
		{
			RefID:       "USER02",
			SearchValue: "test02@test.com",
		},
		{
			RefID:       "USER03",
			SearchValue: "test03@test.com",
		},
	}

	for _, v := range data {
		value := NewSearchValue().
			SetSourceReferenceID(v.RefID).
			SetSearchValue(v.SearchValue)

		err = store.SearchValueCreate(value)

		if err != nil {
			t.Fatal("unexpected error:", err)
			return
		}

	}

	refsFound, errFind := store.Search("test02@test.com", SEARCH_TYPE_EQUALS)

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(refsFound) != 1 {
		t.Fatal("Search MUST return 1")
		return
	}

	if refsFound[0] != "USER02" {
		t.Fatal("Reference ID found MUST BE 'USER02', found: ", refsFound[0])
		return
	}
}

func Test_Store_SearchContains(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_contains",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	data := []struct {
		RefID       string
		SearchValue string
	}{
		{
			RefID:       "USER01",
			SearchValue: "test01@test.com",
		},
		{
			RefID:       "USER021",
			SearchValue: "test021@test.com",
		},
		{
			RefID:       "USER022",
			SearchValue: "test022@test.com",
		},
		{
			RefID:       "USER03",
			SearchValue: "test03@test.com",
		},
	}

	for _, v := range data {
		value := NewSearchValue().
			SetSourceReferenceID(v.RefID).
			SetSearchValue(v.SearchValue)

		err = store.SearchValueCreate(value)

		if err != nil {
			t.Fatal("unexpected error:", err)
			return
		}

	}

	refsFound, errFind := store.Search("st02", SEARCH_TYPE_CONTAINS)

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(refsFound) != 2 {
		t.Fatal("Search MUST return exactly 2 references. Returned: ", len(refsFound))
		return
	}

	if refsFound[0] != "USER021" {
		t.Fatal("Reference ID found MUST BE 'USER021', found: ", refsFound[0])
		return
	}

	if refsFound[1] != "USER022" {
		t.Fatal("Reference ID found MUST BE 'USER022', found: ", refsFound[1])
		return
	}
}

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
	_ = os.Remove(filepath) // remove database
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
		Transformer:        &NoChangeTransformer{},
	})

	if errAutomigrateFalse != nil {
		t.Fatalf("automigrateEnabled: Expected [err] to be nill received [%v]", errAutomigrateFalse.Error())
	}

	if storeAutomigrateFalse.IsAutomigrateEnabled() != false {
		t.Fatalf("automigrateEnabled: Expected [false] received [%v]", storeAutomigrateFalse.IsAutomigrateEnabled())
	}

	storeAutomigrateTrue, errAutomigrateTrue := NewStore(NewStoreOptions{
		TableName:          "test_blockindex_with_automigrate_true",
		DB:                 db,
		AutomigrateEnabled: true,
		Transformer:        &NoChangeTransformer{},
	})

	if errAutomigrateTrue != nil {
		t.Fatalf("automigrateEnabled: Expected [err] to be nill received [%v]", errAutomigrateTrue.Error())
	}

	if storeAutomigrateTrue.IsAutomigrateEnabled() != true {
		t.Fatalf("automigrateEnabled: Expected [true] received [%v]", storeAutomigrateTrue.IsAutomigrateEnabled())
	}
}

func Test_Store_SearchValueCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_create",
		AutomigrateEnabled: true,
		Transformer:        &Sha256Transformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	searchValue := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(searchValue)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if searchValue.SearchValue() != "ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9" {
		t.Fatal("Search value MUST BE 'ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9', found: ", searchValue.SearchValue())
	}

}

func Test_Store_SearchValueUpdate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_create",
		AutomigrateEnabled: true,
		Transformer:        &Sha256Transformer{},
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	searchValue := NewSearchValue().
		SetSourceReferenceID("RefId01").
		SetSearchValue("SearchValue01")

	err = store.SearchValueCreate(searchValue)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if searchValue.SearchValue() != "ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9" {
		t.Fatal("Search value MUST BE 'ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9', found: ", searchValue.SearchValue())
	}

	err = store.SearchValueUpdate(searchValue)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if searchValue.SearchValue() != "ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9" {
		t.Fatal("Search value MUST BE 'ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9', found: ", searchValue.SearchValue())
	}

	searchValue.SetSearchValue("SearchValue01Changed")
	err = store.SearchValueUpdate(searchValue)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if searchValue.SearchValue() != "cca12edb3e3e0ba645dd430f197ff27df630fe4ec7269d31d2fab51f13d1b01a" {
		t.Fatal("Search value MUST BE 'cca12edb3e3e0ba645dd430f197ff27df630fe4ec7269d31d2fab51f13d1b01a', found: ", searchValue.SearchValue())
	}

}

func Test_Store_SearchValueFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_find_by_id",
		AutomigrateEnabled: true,
		Transformer:        &Sha256Transformer{},
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

	if valueFound.SearchValue() != "ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9" {
		t.Fatal("Search value MUST BE 'ef46c0effb3e3a6d65fbbd46c039008205e67b8089339db1852ca0992804afb9', found: ", valueFound.SearchValue())
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
		Transformer:        &NoChangeTransformer{},
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
		Transformer:        &NoChangeTransformer{},
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
		Transformer:        &NoChangeTransformer{},
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

func Test_Store_SearchContains_NoChangeTransformer(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_contains",
		AutomigrateEnabled: true,
		Transformer:        &NoChangeTransformer{},
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

func Test_Store_SearchContains_Rot13Transformer(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_contains",
		AutomigrateEnabled: true,
		Transformer:        &Rot13Transformer{},
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

func Test_Store_SearchContains_Sha256Transformer(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_contains",
		AutomigrateEnabled: true,
		Transformer:        &Sha256Transformer{},
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

	if len(refsFound) != 0 {
		t.Fatal("Search MUST return exactly 0 references. Returned: ", len(refsFound))
		return
	}

	refsFound2, errFind := store.Search("test022@test.com", SEARCH_TYPE_EQUALS)

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if len(refsFound2) != 1 {
		t.Fatal("Search MUST return exactly 1 references. Returned: ", len(refsFound))
		return
	}

	if refsFound2[0] != "USER022" {
		t.Fatal("Reference ID found MUST BE 'USER022', found: ", refsFound2[0])
		return
	}
}

func Test_Store_SearchContains_UniTransformer(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		TableName:          "test_blindindex_value_search_contains",
		AutomigrateEnabled: true,
		Transformer:        &UniTransformer{},
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

package blindindexstore

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/doug-martin/goqu/v9"

	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

var _ StoreInterface = (*Store)(nil) // verify it extends the interface

// Store defines a session store
type Store struct {
	tableName          string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
}

// AutoMigrate auto migrate
func (st *Store) AutoMigrate() error {
	sql := st.sqlTableCreate()

	if st.debugEnabled {
		log.Println(sql)
	}

	_, err := st.db.Exec(sql)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

func (store *Store) SearchValueCreate(searchValue *SearchValue) error {
	searchValue.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	searchValue.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := searchValue.Data()

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Insert(store.tableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	if err != nil {
		return err
	}

	searchValue.MarkAsNotDirty()

	return nil
}

func (store *Store) SearchValueDelete(searchValue *SearchValue) error {
	if searchValue == nil {
		return errors.New("searchValue is nil")
	}

	return store.SearchValueDeleteByID(searchValue.ID())
}

func (store *Store) SearchValueDeleteByID(id string) error {
	if id == "" {
		return errors.New("searchValue id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.tableName).
		Prepared(true).
		Where(goqu.C("id").Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

func (store *Store) SearchValueFindByID(id string) (*SearchValue, error) {
	if id == "" {
		return nil, errors.New("searchValue id is empty")
	}

	list, err := store.SearchValueList(SearchValueQueryOptions{
		ID:    id,
		Limit: 1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

func (store *Store) SearchValueFindBySourceReferenceID(sourceReferenceID string) (*SearchValue, error) {
	if sourceReferenceID == "" {
		return nil, errors.New("searchValue objectID is empty")
	}

	list, err := store.SearchValueList(SearchValueQueryOptions{
		SourceReferenceID: sourceReferenceID,
		Limit:             1,
	})

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

func (store *Store) SearchValueList(options SearchValueQueryOptions) ([]SearchValue, error) {
	q := store.searchValueQuery(options)

	sqlStr, _, errSql := q.Select().ToSQL()

	if errSql != nil {
		return []SearchValue{}, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr)
	if err != nil {
		return []SearchValue{}, err
	}

	list := []SearchValue{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewSearchValueFromExistingData(modelMap)
		list = append(list, *model)
	})

	return list, nil
}

func (store *Store) SearchValueSoftDelete(searchValue *SearchValue) error {
	if searchValue == nil {
		return errors.New("searchValue is nil")
	}

	searchValue.SetDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.SearchValueUpdate(searchValue)
}

func (store *Store) SearchValueSoftDeleteByID(id string) error {
	searchValue, err := store.SearchValueFindByID(id)

	if err != nil {
		return err
	}

	return store.SearchValueSoftDelete(searchValue)
}

func (store *Store) SearchValueUpdate(searchValue *SearchValue) error {
	if searchValue == nil {
		return errors.New("order is nil")
	}

	// searchValue.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := searchValue.DataChanged()

	delete(dataChanged, "id")   // ID is not updateable
	delete(dataChanged, "hash") // Hash is not updateable
	delete(dataChanged, "data") // Data is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Update(store.tableName).
		Prepared(true).
		Set(dataChanged).
		Where(goqu.C("id").Eq(searchValue.ID())).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := store.db.Exec(sqlStr, params...)

	searchValue.MarkAsNotDirty()

	return err
}

func (store *Store) searchValueQuery(options SearchValueQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(store.dbDriverName).From(store.tableName)

	if options.ID != "" {
		q = q.Where(goqu.C("id").Eq(options.ID))
	}

	if options.SourceReferenceID != "" {
		q = q.Where(goqu.C(COLUMN_SOURCE_REFERENCE_ID).Eq(options.SourceReferenceID))
	}

	if options.SearchValue != "" {
		q = q.Where(goqu.C(COLUMN_SEARCH_VALUE).Eq(options.SearchValue))
	}

	if !options.CountOnly {
		if options.Limit > 0 {
			q = q.Limit(uint(options.Limit))
		}

		if options.Offset > 0 {
			q = q.Offset(uint(options.Offset))
		}
	}

	sortOrder := "desc"
	if options.SortOrder != "" {
		sortOrder = options.SortOrder
	}

	if options.OrderBy != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy).Desc())
		}
	}

	if !options.WithDeleted {
		q = q.Where(goqu.C(COLUMN_DELETED_AT).Gte(carbon.Now(carbon.UTC).ToDateTimeString()))
	}

	return q
}
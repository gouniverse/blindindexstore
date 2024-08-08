# Blind Index Store

Allows to create blind index database tables for searching.

The implementation is generic enough to allow custom transformers (recommended)

## Usage

1. Instantiate the blind index store, and add your transformer.

```golang
store, err := NewStore(NewStoreOptions{
    DB:                 db,
    TableName:          "blindindex_emails",
    AutomigrateEnabled: true,
    Transformer:        &Sha256Transformer{},
})

if err != nil {
    t.Fatal("unexpected error:", err)
}

if store == nil {
    t.Fatal("unexpected nil store")
}
```

2. Populate the index with your search values:

```golang
searchValue := NewSearchValue().
    SetSourceReferenceID("USER01").
    SetSearchValue("user01@test.com")

err = store.SearchValueCreate(searchValue)

if err != nil {
    t.Fatal("unexpected error:", err)
}
```

3. Search the index:

```golang
refsFound, errFind := store.Search("user01@test.com", SEARCH_TYPE_EQUALS)

if errFind != nil {
    t.Fatal("unexpected error:", errFind)
    return
}

// Refs found: [ "USER01" ]
```

## Frequently Asked Questions:


1. What is a transformer?

The transformer "blinds" the search value, before storing the data in the blind index.

Recommended is to create your own transformer, implementing the TransformerInterface.

2. What types of search are supported?

The supported search types depend on your transformer.

Provided search options are: SEARCH_TYPE_EQUALS, SEARCH_TYPE_CONTAINS, SEARCH_TYPE_STARTS_WITH, SEARCH_TYPE_ENDS_WITH.

A transformer may support one (i.e. SEARCH_TYPE_EQUALS) or more of them.
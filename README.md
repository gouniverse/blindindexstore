# Blind Index Store <a href="https://gitpod.io/#https://github.com/gouniverse/blindindexstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/gouniverse/blindindexstore/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/gouniverse/blindindexstore/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/blindindexstore)](https://goreportcard.com/report/github.com/gouniverse/blindindexstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/blindindexstore)](https://pkg.go.dev/github.com/gouniverse/blindindexstore)

The Blind Index Store allows for creating blind index database tables for searching encrypted data. 

It supports a wide range of custom transformers (recommended) for increased flexibility.

## License

This project is licensed under the GNU General Public License version 3 (AGPL-3.0). You can find a copy of the license at https://www.gnu.org/licenses/agpl-3.0.en.html

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

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

### 1. What is a blind index?

A blind index is a secure search mechanism that allows for searching encrypted data without revealing the underlying plaintext.
It involves hashing encrypted data, storing the hash (blind index) in a separate index, and comparing search queries to these hashes.
This technique ensures privacy and security while enabling efficient searches.

### 2. How does a blind index work?

The process involves encrypting data, hashing (or obfuscating) the encrypted data to create a blind index, storing the blind index in a separate index, and comparing search queries to these indexes.

### 3. What are the benefits of using a blind index?

- Privacy: Protects sensitive data by keeping it encrypted.
- Security: The hash function ensures that the original data cannot be easily recovered.
- Efficiency: Allows for efficient searching through encrypted data.

### 4. What is a transformer?

A transformer is a crucial component in a blind index system. It's responsible for "blinding" or obfuscating the search value before it's stored in the index.

This process ensures that the search value remains confidential, protecting sensitive data.

Recommendation: It's generally advisable to create your own custom transformer by implementing the TransformerInterface.
This allows for greater flexibility and control over the blinding process, tailoring it to your specific security requirements.

### 5. How do I create a custom transformer?

To create a custom transformer, implement the TransformerInterface. This gives you flexibility in customizing the blinding process.

### 6. What types of search are supported?

The available search types vary depending on your transformer.

Here are the potential options:

- SEARCH_TYPE_EQUALS: Exact match.
- SEARCH_TYPE_CONTAINS: Partial match within the string.
- SEARCH_TYPE_STARTS_WITH: Partial match at the beginning of the string.
- SEARCH_TYPE_ENDS_WITH: Partial match at the end of the string.

A transformer might support one or more of these search types.

Note: A fully "blinded" index typically only allows for SEARCH_TYPE_EQUALS searches, ensuring the highest level of privacy.

### 7. How do I instantiate the blind index store?
Instantiate the store by providing the database connection, table name, automigration settings, and your custom transformer.

### 8. How do I populate the index with search values?
Create a new SearchValue object, set the source reference ID and search value, and then call the SearchValueCreate method on the store.

### 9. How do I search the index?
Call the Search method on the store, providing the search term and the desired search type.

### 10. Can I use a different hash function?
Yes, you can customize the hash function used for creating the blind index. However, ensure that the chosen hash function is secure and appropriate for your use case.

### 11. Is the blind index store secure?
The security of the blind index store depends on the underlying encryption algorithm, hash function, and the implementation of your custom transformer. It's essential to choose strong cryptographic primitives and implement them securely.

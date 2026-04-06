# Search Guide

This guide covers the generic search system in gFly — a driver-based, fluent search API inspired by **Laravel Scout** that supports both PostgreSQL (no external service required) and Elasticsearch.

---

## Table of Contents

1. [Overview](#overview)
2. [Package Structure](#package-structure)
3. [Core Concepts](#core-concepts)
   - [Driver](#driver)
   - [Searchable Interface](#searchable-interface)
   - [Engine](#engine)
   - [Builder](#builder)
   - [Result & Hit](#result--hit)
4. [Drivers](#drivers)
   - [DatabaseDriver](#databasedriver)
   - [ElasticsearchDriver](#elasticsearchdriver)
5. [Making a Model Searchable](#making-a-model-searchable)
6. [Searching](#searching)
   - [Basic Search](#basic-search)
   - [Filtering](#filtering)
   - [Sorting](#sorting)
   - [Pagination](#pagination)
   - [Count Only](#count-only)
   - [Ad-hoc Search Without Searchable](#ad-hoc-search-without-searchable)
7. [Indexing (Elasticsearch)](#indexing-elasticsearch)
   - [Index a Single Model](#index-a-single-model)
   - [Update Index](#update-index)
   - [Remove a Model](#remove-a-model)
   - [Bulk Index](#bulk-index)
8. [Hydrating Results](#hydrating-results)
9. [Console Commands](#console-commands)
10. [Complete Example](#complete-example)
11. [Best Practices](#best-practices)

---

## Overview

```
Search Query
    ↓
Engine.For(model) / Engine.Index(name, fields...)
    ↓
Builder  →  .Query()  .Where()  .OrderBy()  .Page()  .PerPage()
    ↓
Driver.Search(Request)
    ↓
Result { Total, Hits[]{ID, Score, Data} }
    ↓
Hydrate IDs → []Model  (fetch from DB via mb.GetModelByID)
```

Two drivers are available out of the box:

| Driver | Backend | Indexing needed? | Fuzzy search |
|--------|---------|-----------------|--------------|
| `DatabaseDriver` | PostgreSQL `ILIKE` | No | No |
| `ElasticsearchDriver` | Elasticsearch 7.x/8.x | Yes | Yes |

---

## Package Structure

The search engine is provided by the external package `github.com/gflydev/search`. Your application code lives in:

```
internal/domain/models/
└── user_searchable.go          → Searchable interface implementation for models.User

internal/services/
└── user_search_services.go     → Engine init + SearchUsers, AddIndexUser,
                                  UpdateIndexUser, RemoveIndexUser, BulkIndexUsers

internal/console/commands/
└── user_search_command.go      → CLI commands: user-search, user-index-bulk
```

---

## Core Concepts

### Driver

`Driver` is the interface every search backend must implement:

```go
type Driver interface {
    Search(req Request) (*Result, error)
    Index(indexName string, id any, data core.Data) error
    Remove(indexName string, id any) error
    BulkIndex(indexName string, docs []Document) error
}
```

`Index`, `Remove`, and `BulkIndex` are no-ops for the `DatabaseDriver` — the database is already the source of truth.

---

### Searchable Interface

Any model that should be searchable must implement the four-method `Searchable` interface:

```go
type Searchable interface {
    SearchIndex() string            // table name (DB) or index name (ES)
    SearchKey() any                 // primary key
    SearchableFields() []string     // columns/fields to search against
    ToSearchDocument() core.Data    // document written to ES index
}
```

> **Field name convention:** Use **unqualified** field names (e.g. `"fullname"`, `"email"`) so they
> work both as plain SQL column names (single-table queries) and as Elasticsearch document field
> names. Only use qualified names (`"users.fullname"`) when the `DatabaseDriver` query involves
> explicit JOINs to avoid ambiguous column errors.

---

### Engine

`Engine` is the main entry point. Create one per driver and reuse it:

```go
engine := search.New(search.NewDatabaseDriver())
```

| Method | Description |
|--------|-------------|
| `engine.For(model Searchable)` | Returns a `Builder` pre-configured from the model |
| `engine.Index(name, fields...)` | Returns a `Builder` with explicit index and field list |
| `engine.IndexModel(model)` | Indexes a single model (ES) |
| `engine.RemoveModel(model)` | Removes a model from the index (ES) |
| `engine.BulkIndex(models)` | Bulk-indexes a slice of models (ES) |

---

### Builder

`Builder` provides the fluent query API:

```go
engine.For(models.User{}).
    Query("alice").
    Where("status", "active").
    OrderBy("id", search.Desc).
    Page(1).PerPage(20).
    Search()
```

| Method | Description |
|--------|-------------|
| `.Query(q string)` | Full-text search string |
| `.Where(field, value)` | Equality filter (AND-ed, multiple allowed) |
| `.OrderBy(field, direction)` | Sort by field (`search.Asc` / `search.Desc`) |
| `.Page(n int)` | 1-based page number (default: 1) |
| `.PerPage(n int)` | Results per page (default: 15) |
| `.Search()` | Execute and return `*Result` |
| `.Paginate()` | Alias for `.Search()` |
| `.Count()` | Return total matching count only |

---

### Result & Hit

```go
type Result struct {
    Total   int64   // total matching documents (before pagination)
    Page    int
    PerPage int
    Hits    []Hit
}

type Hit struct {
    ID    any        // primary key
    Score float64    // 1.0 for DB driver, ES _score for ES driver
    Data  core.Data  // ES _source document (nil for DB driver)
}
```

Convenience helpers on `*Result`:

```go
result.IDs()     // []any
result.IntIDs()  // []int  (skips non-int IDs)
```

---

## Drivers

### DatabaseDriver

Uses PostgreSQL `ILIKE` (case-insensitive) queries. No additional infrastructure is required.

```go
driver := search.NewDatabaseDriver()
engine := search.New(driver)
```

**Options:**

```go
driver := &search.DatabaseDriver{
    CaseSensitive: true,  // use LIKE instead of ILIKE (default: false)
    IDField:       "id",  // primary key column name (default: "id")
}
```

**How it works:**

For a `Query("alice")` across `["fullname", "email"]` the driver executes:

```sql
-- Count
SELECT COUNT(*) FROM users
WHERE (fullname ILIKE '%alice%' OR email ILIKE '%alice%')

-- Paginated IDs
SELECT id FROM users
WHERE (fullname ILIKE '%alice%' OR email ILIKE '%alice%')
ORDER BY id DESC
LIMIT 20 OFFSET 0
```

---

### ElasticsearchDriver

Uses the Elasticsearch REST API via Go's standard `net/http` — no extra dependency required.

```go
driver := search.NewElasticsearchDriver(search.ElasticsearchConfig{
    Host:     coreUtils.Getenv("ES_HOST", "http://localhost:9200"),  // from env
    Username: "",                       // optional (HTTP Basic Auth)
    Password: "",
    Timeout:  10 * time.Second,        // optional (default: 10 s)
})
engine := search.New(driver)
```

**How it works:**

A `Query("alice")` across `["fullname", "email"]` with `Where("status", "active")` generates:

```json
{
  "query": {
    "bool": {
      "must": [{
        "multi_match": {
          "query": "alice",
          "fields": ["fullname", "email"],
          "type": "best_fields",
          "fuzziness": "AUTO"
        }
      }],
      "filter": [{ "term": { "status": "active" } }]
    }
  },
  "from": 0,
  "size": 20
}
```

The driver returns ES `_id` and `_score` in each `Hit`. The `_source` document is available via `hit.Data`.

---

## Making a Model Searchable

Add a dedicated file inside the model's package (e.g. `internal/domain/models/user_searchable.go`):

```go
package models

import "github.com/gflydev/core"

// SearchIndex returns the table / Elasticsearch index name for User.
func (u User) SearchIndex() string { return TableUser }

// SearchKey returns the primary key used to address the document.
func (u User) SearchKey() any { return u.ID }

// SearchableFields lists the columns the DatabaseDriver will ILIKE against
// when a keyword is supplied in the search request.
// Field names are unqualified so they work in both SQL (single-table query)
// and as Elasticsearch document field names.
func (u User) SearchableFields() []string {
    return []string{
        "fullname",
        "email",
        "phone",
    }
}

// ToSearchDocument returns the flat map written to the Elasticsearch index
// when this user is indexed.
func (u User) ToSearchDocument() core.Data {
    return core.Data{
        "id":       u.ID,
        "fullname": u.Fullname,
        "email":    u.Email,
        "phone":    u.Phone,
        "status":   string(u.Status),
    }
}
```

> **Note:** `SearchableFields` uses plain column/field names (no table prefix). This works for
> both the `DatabaseDriver` (single-table query) and the `ElasticsearchDriver` (document field
> names). If your `DatabaseDriver` query requires JOINs, qualify the names (`"users.fullname"`)
> to prevent ambiguous column errors.

---

## Searching

### Basic Search

```go
engine := search.New(search.NewElasticsearchDriver(search.ElasticsearchConfig{
    Host: coreUtils.Getenv("ES_HOST", "http://localhost:9200"),
}))

result, err := engine.For(models.User{}).
    Query("alice").
    Search()

// result.Total → total matching records
// result.Hits  → []Hit with ID and Score
```

### Filtering

Additional `Where` calls are AND-ed together. Use the same field names as in `ToSearchDocument`:

```go
result, err := engine.For(models.User{}).
    Query("alice").
    Where("status", "active").
    Search()
```

> For the `DatabaseDriver` with JOINs, filter field names must be qualified (`"users.status"`).
> For the `ElasticsearchDriver`, filter field names match the document field names in the index.

### Sorting

```go
result, err := engine.For(models.User{}).
    Query("alice").
    OrderBy("id", search.Desc).
    Search()
```

Constants: `search.Asc`, `search.Desc`

### Pagination

```go
result, err := engine.For(models.User{}).
    Query("alice").
    Page(2).PerPage(15).
    Paginate()

// result.Total   → total count across all pages
// result.Page    → 2
// result.PerPage → 15
// result.Hits    → up to 15 hits for page 2
```

### Count Only

```go
total, err := engine.For(models.User{}).
    Query("alice").
    Where("status", "active").
    Count()
```

### Ad-hoc Search Without Searchable

When you do not want to implement the full `Searchable` interface, use `engine.Index()` directly:

```go
result, err := engine.
    Index("users", "fullname", "email").
    Query("alice").
    Where("status", "active").
    PerPage(10).
    Search()
```

---

## Indexing (Elasticsearch)

The `DatabaseDriver` does not require indexing — skip this section if you only use PostgreSQL search.

### Index a Single Model

Call after **creating** a model:

```go
func AddIndexUser(user models.User) error {
    if err := ESSearchEngine.IndexModel(user); err != nil {
        log.Errorf("AddIndexUser: failed to index user %d: %v", user.ID, err)
        return errors.New("error occurs while indexing user")
    }

    log.Infof("AddIndexUser: indexed user %d (%s)", user.ID, user.Email)

    return nil
}
```

### Update Index

Call after **updating** a model. `IndexModel` performs an upsert so the same function works for both create and update:

```go
func UpdateIndexUser(user models.User) error {
    if err := ESSearchEngine.IndexModel(user); err != nil {
        log.Errorf("UpdateIndexUser: failed to re-index user %d: %v", user.ID, err)
        return errors.New("error occurs while updating user index")
    }

    log.Infof("UpdateIndexUser: re-indexed user %d (%s)", user.ID, user.Email)

    return nil
}
```

### Remove a Model

Call before or after **deleting** a model:

```go
func RemoveIndexUser(user models.User) error {
    if err := ESSearchEngine.RemoveModel(user); err != nil {
        log.Errorf("RemoveIndexUser: failed to remove user %d from index: %v", user.ID, err)
        return errors.New("error occurs while removing user from index")
    }

    log.Infof("RemoveIndexUser: removed user %d from index", user.ID)

    return nil
}
```

### Bulk Index

Use for initial import or scheduled full re-sync. Convert the typed slice to `[]search.Searchable` first:

```go
func BulkIndexUsers(users []models.User) error {
    searchableData := make([]search.Searchable, len(users))
    for idx := range users {
        searchableData[idx] = users[idx]
    }

    if err := ESSearchEngine.BulkIndex(searchableData); err != nil {
        log.Errorf("BulkIndexUsers: bulk index failed: %v", err)
        return errors.New("error occurs while bulk indexing users")
    }

    log.Infof("BulkIndexUsers: indexed %d users", len(users))

    return nil
}
```

---

## Hydrating Results

Both drivers return `Hit.ID` values (primary keys). Use a private helper to fetch full model structs from the database, silently skipping missing records:

```go
func hydrateUsersByIDs(ids []int) ([]models.User, error) {
    users := make([]models.User, 0, len(ids))

    for _, id := range ids {
        user, err := mb.GetModelByID[models.User](id)
        if err != nil {
            log.Warnf("hydrateUsersByIDs: user %d not found, skipping", id)
            continue
        }
        users = append(users, *user)
    }

    return users, nil
}
```

> **Elasticsearch only:** `hit.Data` contains the raw `_source` document already available in the
> `Result`. If the indexed document contains all required fields you can use `hit.Data` directly
> and skip the database round-trip.

---

## Console Commands

Register console commands in `internal/console/commands/` for CLI-driven search operations:

```go
func init() {
    console.RegisterCommand(&userSearchCommand{}, "user-search")
    console.RegisterCommand(&userIndexBulkCommand{}, "user-index-bulk")
}
```

**`user-index-bulk`** — fetches all non-deleted users from the database and bulk-indexes them. Run this once to populate the index or after a schema change:

```go
type userIndexBulkCommand struct {
    console.Command
}

func (c *userIndexBulkCommand) Handle() {
    log.Info("=== user-index-bulk: starting ===")

    var users []models.User

    total, err := mb.Instance().
        Model(&models.User{}).
        Where(models.TableUser+".deleted_at", mb.Null, nil).
        Find(&users)

    if err != nil {
        log.Errorf("user-index-bulk: failed to load users: %v", err)
        return
    }

    log.Infof("user-index-bulk: loaded %d users from database", total)

    if err = services.BulkIndexUsers(users); err != nil {
        log.Errorf("user-index-bulk: %v", err)
        return
    }

    log.Infof("=== user-index-bulk: done at %s ===", time.Now().Format("2006-01-02 15:04:05"))
}
```

Run via artisan:

```bash
./build/artisan cmd:run user-index-bulk
./build/artisan cmd:run user-search
```

---

## Complete Example

### 1. Implement Searchable (`internal/domain/models/user_searchable.go`)

```go
package models

import "github.com/gflydev/core"

// SearchIndex returns the table / Elasticsearch index name for User.
func (u User) SearchIndex() string { return TableUser }

// SearchKey returns the primary key used to address the document.
func (u User) SearchKey() any { return u.ID }

// SearchableFields lists the columns the DatabaseDriver will ILIKE against.
// Use unqualified names — they work for both single-table SQL and ES document fields.
func (u User) SearchableFields() []string {
    return []string{
        "fullname",
        "email",
        "phone",
    }
}

// ToSearchDocument returns the flat map written to the Elasticsearch index.
func (u User) ToSearchDocument() core.Data {
    return core.Data{
        "id":       u.ID,
        "fullname": u.Fullname,
        "email":    u.Email,
        "phone":    u.Phone,
        "status":   string(u.Status),
    }
}
```

### 2. Create service file (`internal/services/user_search_services.go`)

```go
package services

import (
    "gfly/internal/domain/models"
    "github.com/gflydev/search"

    "github.com/gflydev/core/errors"
    "github.com/gflydev/core/log"
    coreUtils "github.com/gflydev/core/utils"
    mb "github.com/gflydev/db"
)

// ESSearchEngine is the application-wide Elasticsearch search engine.
// The host is read from the ES_HOST environment variable.
var ESSearchEngine = search.New(search.NewElasticsearchDriver(search.ElasticsearchConfig{
    Host: coreUtils.Getenv("ES_HOST", "http://localhost:9200"),
}))

// SearchUsers searches users in Elasticsearch by keyword with optional status
// filter and returns hydrated models.User slice alongside the total count.
func SearchUsers(keyword, status string, page, perPage int) ([]models.User, int64, error) {
    builder := ESSearchEngine.For(models.User{}).
        Query(keyword).
        Page(page).
        OrderBy("id", "desc").
        PerPage(perPage)

    if status != "" {
        builder = builder.Where("status", status)
    }

    result, err := builder.Search()
    if err != nil {
        log.Errorf("SearchUsers: elasticsearch query failed: %v", err)
        return nil, 0, errors.New("error occurs while searching users")
    }

    users, err := hydrateUsersByIDs(result.IntIDs())
    if err != nil {
        return nil, 0, err
    }

    return users, result.Total, nil
}

// AddIndexUser indexes a newly created user into Elasticsearch.
func AddIndexUser(user models.User) error {
    if err := ESSearchEngine.IndexModel(user); err != nil {
        log.Errorf("AddIndexUser: failed to index user %d: %v", user.ID, err)
        return errors.New("error occurs while indexing user")
    }
    log.Infof("AddIndexUser: indexed user %d (%s)", user.ID, user.Email)
    return nil
}

// UpdateIndexUser re-indexes an updated user in Elasticsearch.
func UpdateIndexUser(user models.User) error {
    if err := ESSearchEngine.IndexModel(user); err != nil {
        log.Errorf("UpdateIndexUser: failed to re-index user %d: %v", user.ID, err)
        return errors.New("error occurs while updating user index")
    }
    log.Infof("UpdateIndexUser: re-indexed user %d (%s)", user.ID, user.Email)
    return nil
}

// RemoveIndexUser removes a user from the Elasticsearch index.
func RemoveIndexUser(user models.User) error {
    if err := ESSearchEngine.RemoveModel(user); err != nil {
        log.Errorf("RemoveIndexUser: failed to remove user %d from index: %v", user.ID, err)
        return errors.New("error occurs while removing user from index")
    }
    log.Infof("RemoveIndexUser: removed user %d from index", user.ID)
    return nil
}

// BulkIndexUsers re-indexes all provided users in a single Elasticsearch bulk request.
func BulkIndexUsers(users []models.User) error {
    searchableData := make([]search.Searchable, len(users))
    for idx := range users {
        searchableData[idx] = users[idx]
    }

    if err := ESSearchEngine.BulkIndex(searchableData); err != nil {
        log.Errorf("BulkIndexUsers: bulk index failed: %v", err)
        return errors.New("error occurs while bulk indexing users")
    }

    log.Infof("BulkIndexUsers: indexed %d users", len(users))
    return nil
}

// hydrateUsersByIDs fetches full User models from the database for the given IDs.
// Missing or deleted records are silently skipped.
func hydrateUsersByIDs(ids []int) ([]models.User, error) {
    users := make([]models.User, 0, len(ids))

    for _, id := range ids {
        user, err := mb.GetModelByID[models.User](id)
        if err != nil {
            log.Warnf("hydrateUsersByIDs: user %d not found, skipping", id)
            continue
        }
        users = append(users, *user)
    }

    return users, nil
}
```

### 3. Use in a controller (`internal/http/controllers/api/user/`)

```go
func (h *SearchUsersApi) Handle(c *core.Ctx) error {
    keyword := c.QueryParam("q")
    status  := c.QueryParam("status")
    page, _ := strconv.Atoi(c.QueryParam("page"))
    if page < 1 {
        page = 1
    }

    users, total, err := services.SearchUsers(keyword, status, page, 20)
    if err != nil {
        return c.Error(http.Error{Code: "SEARCH_ERROR", Message: err.Error()})
    }

    data := http.ToListResponse(users, transformers.ToUserResponse)
    return c.Success(http.List[response.User]{
        Meta: http.Meta{Page: page, PerPage: 20, Total: int(total)},
        Data: data,
    })
}
```

### 4. Wire index calls into the main user service

Call the index functions immediately after successful create/update/delete operations:

```go
// After create
user, err := CreateUser(dto)
if err == nil {
    _ = AddIndexUser(*user)
}

// After update
user, err := UpdateUser(dto)
if err == nil {
    _ = UpdateIndexUser(*user)
}

// Before/after delete
_ = RemoveIndexUser(user)
DeleteUserByID(user.ID)
```

---

## Best Practices

| Topic | Recommendation |
|-------|----------------|
| **Driver selection** | Use `DatabaseDriver` for simple keyword search on small datasets with no extra infra. Use `ElasticsearchDriver` for fuzzy matching, relevance scoring, or high query throughput. |
| **Field names** | Use unqualified field names in `SearchableFields` and `Where`/`OrderBy` calls — they match both SQL columns (single-table) and ES document fields. Only qualify (`table.column`) when the `DatabaseDriver` query involves JOINs. |
| **Engine reuse** | Define `ESSearchEngine` as a package-level var in the service file and reuse it across all functions — the underlying HTTP client should be shared. |
| **Named service functions** | Create distinct named functions per operation (`AddIndexUser`, `UpdateIndexUser`, `RemoveIndexUser`, `BulkIndexUsers`) rather than calling `ESSearchEngine` directly in controllers or other services. |
| **ES indexing** | Call `AddIndexUser` / `UpdateIndexUser` inside the service layer, immediately after a successful create or update, to keep the index in sync. |
| **Bulk re-index** | Implement `BulkIndexUsers` as a console command (`internal/console/commands/`) and run it via `./build/artisan cmd:run user-index-bulk`. Never call it inline in a request handler. |
| **Hydration** | Always fetch full models from PostgreSQL using `mb.GetModelByID` in a private `hydrateByIDs` helper. Skip missing records with a warn log rather than returning an error. |
| **Empty query** | An empty `Query("")` skips the full-text match and returns all documents matching only the `Where` filters, behaving like a normal filtered list. |
| **Sorting with ES** | Only sort on `keyword` or `numeric` field types in ES. Sorting on `text` fields requires a `.keyword` sub-field mapping (`field.keyword`). |
| **Error wrapping** | Use `errors.New("context message")` from `github.com/gflydev/core/errors` and log the original error with `log.Errorf` before returning a generic message to the caller. |

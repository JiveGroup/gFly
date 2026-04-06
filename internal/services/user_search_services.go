package services

import (
	"gfly/internal/domain/models"
	"github.com/gflydev/search"

	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	coreUtils "github.com/gflydev/core/utils"
	mb "github.com/gflydev/db"
)

// ====================================================================
// ========================= Search Engine ============================
// ====================================================================

// searchEngine is the application-wide Elasticsearch search engine for users.
// The host is read from the ES_HOST environment variable (default: http://localhost:9200).
var searchEngine = search.New(search.NewElasticsearchDriver(search.ElasticsearchConfig{
	Host: coreUtils.Getenv("ES_HOST", "http://localhost:9200"),
}))

// ====================================================================
// ========================= Main functions ===========================
// ====================================================================

// SearchUsers searches users in Elasticsearch by keyword with optional status
// filter and returns hydrated models.User slice alongside the total count.
//
// Parameters:
//   - keyword (string): Full-text search term matched against fullname, email and phone.
//   - status (string): Optional status filter (active / pending / blocked). Pass "" to skip.
//   - page (int): 1-based page number.
//   - perPage (int): Number of results per page.
//
// Returns:
//   - ([]models.User, int64, error): Matched users, total count, and any error.
func SearchUsers(keyword, status string, page, perPage int) ([]models.User, int64, error) {
	builder := searchEngine.For(models.User{}).
		Query(keyword).
		Page(page).
		OrderBy("id", "asc").
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
// Call this immediately after a successful CreateUser.
//
// Parameters:
//   - user (models.User): The user model to index.
//
// Returns:
//   - error: An error if indexing fails.
func AddIndexUser(user models.User) error {
	if err := searchEngine.IndexModel(user); err != nil {
		log.Errorf("AddIndexUser: failed to index user %d: %v", user.ID, err)
		return errors.New("error occurs while indexing user")
	}

	log.Infof("AddIndexUser: indexed user %d (%s)", user.ID, user.Email)

	return nil
}

// UpdateIndexUser re-indexes an updated user in Elasticsearch.
// Call this immediately after a successful UpdateUser.
//
// Parameters:
//   - user (models.User): The updated user model to re-index.
//
// Returns:
//   - error: An error if re-indexing fails.
func UpdateIndexUser(user models.User) error {
	if err := searchEngine.IndexModel(user); err != nil {
		log.Errorf("UpdateIndexUser: failed to re-index user %d: %v", user.ID, err)
		return errors.New("error occurs while updating user index")
	}

	log.Infof("UpdateIndexUser: re-indexed user %d (%s)", user.ID, user.Email)

	return nil
}

// RemoveIndexUser removes a user from the Elasticsearch index.
// Call this before or after a successful DeleteUserByID.
//
// Parameters:
//   - user (models.User): The user model to remove from the index.
//
// Returns:
//   - error: An error if removal fails.
func RemoveIndexUser(user models.User) error {
	if err := searchEngine.RemoveModel(user); err != nil {
		log.Errorf("RemoveIndexUser: failed to remove user %d from index: %v", user.ID, err)
		return errors.New("error occurs while removing user from index")
	}

	log.Infof("RemoveIndexUser: removed user %d from index", user.ID)

	return nil
}

// BulkIndexUsers re-indexes all provided users in a single Elasticsearch bulk
// request.  Useful for initial import or full re-sync jobs.
//
// Parameters:
//   - users ([]models.User): Slice of user models to index.
//
// Returns:
//   - error: An error if the bulk operation fails.
func BulkIndexUsers(users []models.User) error {
	searchableData := make([]search.Searchable, len(users))
	for idx := range users {
		searchableData[idx] = users[idx]
	}

	if err := searchEngine.BulkIndex(searchableData); err != nil {
		log.Errorf("BulkIndexUsers: bulk index failed: %v", err)
		return errors.New("error occurs while bulk indexing users")
	}

	log.Infof("BulkIndexUsers: indexed %d users", len(users))

	return nil
}

// ====================================================================
// ======================== Helper Functions ==========================
// ====================================================================

// hydrateUsersByIDs fetches full User models from the database for the
// given primary key list.  Missing or deleted records are silently skipped.
//
// Parameters:
//   - ids ([]int): Slice of user primary keys returned by the search engine.
//
// Returns:
//   - ([]models.User, error): Hydrated user models and any error encountered.
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

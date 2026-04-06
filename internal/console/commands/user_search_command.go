package commands

import (
	"gfly/internal/services"

	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
	mb "github.com/gflydev/db"

	"gfly/internal/domain/models"
	"time"
)

// ---------------------------------------------------------------
//
//	Register commands.
//	./artisan cmd:run user-search          — full-text search demo
//	./artisan cmd:run user-index-bulk      — bulk re-index all users
//
// ---------------------------------------------------------------

func init() {
	console.RegisterCommand(&userSearchCommand{}, "user-search")
	console.RegisterCommand(&userIndexBulkCommand{}, "user-index-bulk")
}

// ---------------------------------------------------------------
//                     userSearchCommand
// ---------------------------------------------------------------

// userSearchCommand tests SearchUsers against Elasticsearch.
type userSearchCommand struct {
	console.Command
}

// Handle runs a set of search scenarios and prints results to the log.
func (c *userSearchCommand) Handle() {
	log.Info("=== user-search: starting ===")

	scenarios := []struct {
		keyword string
		status  string
	}{
		{"admin", ""},
		{"admin", "active"},
		{"", "active"},
		{"", "pending"},
	}

	for _, s := range scenarios {
		label := s.keyword
		if label == "" {
			label = "(empty)"
		}

		users, total, err := services.SearchUsers(s.keyword, s.status, 1, 10)
		if err != nil {
			log.Errorf("  [keyword=%q status=%q] error: %v", s.keyword, s.status, err)
			continue
		}

		log.Infof("  [keyword=%q status=%q] total=%d returned=%d",
			label, s.status, total, len(users))

		for idx := range users {
			u := users[idx]
			log.Infof("    → id=%-4d email=%-30s status=%s", u.ID, u.Email, u.Status)
		}
	}

	log.Infof("=== user-search: done at %s ===", time.Now().Format("2006-01-02 15:04:05"))
}

// ---------------------------------------------------------------
//                   userIndexBulkCommand
// ---------------------------------------------------------------

// userIndexBulkCommand fetches every user from the database and
// bulk-indexes them into Elasticsearch.  Run this once to populate
// the index or after a schema change.
type userIndexBulkCommand struct {
	console.Command
}

// Handle loads all users and calls BulkIndexUsers.
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

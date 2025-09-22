package commands

import (
	"github.com/gflydev/cache"
	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
	"time"
)

// ---------------------------------------------------------------
//                      Register command.
// ./artisan cmd:run redis-test
// ---------------------------------------------------------------

// Auto-register command.
func init() {
	console.RegisterCommand(&redisCommand{}, "redis-test")
}

// ---------------------------------------------------------------
//                      RedisCommand struct.
// ---------------------------------------------------------------

// RedisCommand struct for hello command.
type redisCommand struct {
	console.Command
}

// Handle Process command.
func (c *redisCommand) Handle() {
	// Add new key
	if err := cache.Set("foo", "Hello world", time.Duration(15*24*3600)*time.Second); err != nil {
		log.Error(err)
	}

	// Get data key
	bar, err := cache.Get("foo")
	if err != nil {
		log.Error(err)
	}
	log.Infof("foo `%v`\n", bar)

	log.Infof("RedisCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}

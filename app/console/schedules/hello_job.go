package schedules

import (
	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
	"time"
)

// ---------------------------------------------------------------
// 					Register job.
// ---------------------------------------------------------------

// Auto-register job into scheduler.
func init() {
	console.RegisterJob(&helloJob{})
}

// ---------------------------------------------------------------
// 					HelloJob struct.
// ---------------------------------------------------------------

// helloJob struct for hello job.
type helloJob struct{}

// GetTime Get time format.
func (c *helloJob) GetTime() string {
	return "0/2 * * * * *"
}

// Handle Process the job.
func (c *helloJob) Handle() {
	log.Infof("HelloJob :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}

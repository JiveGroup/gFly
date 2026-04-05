// Package events registers all application event listeners.
// Import this package with a blank identifier to auto-load listeners:
//
//	_ "gfly/internal/events"

package events

import (
	"gfly/internal/events/user"
	"github.com/gflydev/event"
)

// init auto-registers all event subscribers when this package is imported.
func init() {
	event.Subscribe(&user.UserSubscriber{})
}

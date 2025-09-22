# Page Controllers

This directory contains controllers that handle web page requests and return HTML responses.

## Purpose

Page controllers are responsible for:
- Processing web page requests
- Retrieving data needed for page rendering
- Passing data to view templates
- Handling form submissions
- Managing redirects and flash messages
- Implementing server-side rendering logic

## Structure

- **home_page.go**: Controller for the home page

## Usage

Page controllers typically retrieve data and pass it to view templates:

```
// Example pseudocode for a page controller
// Note: This is not actual implementation code

// In a page controller file (dashboard_page.go):
package page

import (
    "gfly/internal/http/request"
    "gfly/internal/services"
)

// DashboardController handles dashboard page requests
type DashboardController struct {
    userService *services.UserService
    statsService *services.StatsService
}

// NewDashboardController creates a new dashboard controller
func NewDashboardController(userService *services.UserService, statsService *services.StatsService) *DashboardController {
    return &DashboardController{
        userService: userService,
        statsService: statsService,
    }
}

// Show handles GET /dashboard requests
func (c *DashboardController) Show(ctx *fasthttp.RequestCtx) {
    // Get current user
    user, err := c.userService.GetCurrentUser(ctx)
    if err != nil {
        // Redirect to login page
        ctx.Redirect("/login", 302)
        return
    }

    // Get dashboard stats
    stats, err := c.statsService.GetUserStats(user.ID)
    if err != nil {
        // Handle error
        ctx.Error("Failed to load dashboard", 500)
        return
    }

    // Render dashboard template with data
    ctx.SetContentType("text/html")
    ctx.Render("dashboard", map[string]interface{}{
        "user": user,
        "stats": stats,
    })
}
```

## Best Practices

- Keep controllers focused on handling requests and passing data to views
- Use services for business logic
- Implement proper authentication and authorization checks
- Validate form inputs
- Use flash messages for user feedback
- Implement proper error handling and user-friendly error pages
- Use consistent layout templates
- Consider using view models to prepare data for templates
- Implement proper redirects for form submissions
- Use PRG (Post/Redirect/Get) pattern for form submissions
- Document controller methods with clear descriptions
- Test controllers with both valid and invalid inputs

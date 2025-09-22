# Schedules

This directory contains scheduled tasks that run automatically at specified intervals.

## Purpose

The schedules directory is responsible for:
- Defining tasks that run on a schedule
- Specifying the frequency and timing of scheduled tasks
- Implementing background processing for routine operations
- Automating maintenance and cleanup tasks
- Running periodic data processing jobs

## Structure

- **hello_job.go**: Example scheduled job

## Usage

Scheduled tasks are defined as Go structs that implement the scheduler interface:

```
// Example pseudocode for a scheduled task
// Note: This is not actual implementation code

// In a schedule file (cleanup_task.go):
package schedules

import (
    "github.com/gflydev/console"
    "github.com/gflydev/core/log"
    "time"
)

// CleanupTask performs routine cleanup operations
type CleanupTask struct {
    console.Schedule
}

// Register the task with the scheduler
func init() {
    console.RegisterSchedule(&CleanupTask{})
}

// GetSchedule defines when the task should run
func (t *CleanupTask) GetSchedule() string {
    // Run daily at 2:00 AM
    return "0 2 * * *"
}

// Handle executes the scheduled task
func (t *CleanupTask) Handle() {
    log.Info("Starting cleanup task...")

    // Perform cleanup operations
    cleanupTempFiles()
    purgeOldLogs()
    removeExpiredSessions()

    log.Info("Cleanup task completed")
}
```

To run the scheduler:

```bash
./build/artisan schedule:run
```

## Cron Expression Format

Schedules use cron expressions to define when tasks should run:

```
┌─────────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌─────────── day of month (1 - 31)
│ │ │ ┌───────── month (1 - 12)
│ │ │ │ ┌─────── day of week (0 - 6) (Sunday to Saturday)
│ │ │ │ │
* * * * *
```

Common patterns:
- `* * * * *`: Every minute
- `0 * * * *`: Every hour at minute 0
- `0 0 * * *`: Every day at midnight
- `0 0 * * 0`: Every Sunday at midnight
- `0 0 1 * *`: First day of every month at midnight
- `*/15 * * * *`: Every 15 minutes

## Best Practices

- Keep scheduled tasks focused on a single responsibility
- Use descriptive names for task classes
- Log the start and completion of tasks
- Implement error handling and retry logic
- Avoid long-running tasks that could overlap with the next scheduled run
- Consider using queues for resource-intensive operations
- Monitor task execution and set up alerts for failures
- Document the purpose and schedule of each task
- Test tasks in isolation before scheduling them
- Use environment variables for configurable aspects of tasks

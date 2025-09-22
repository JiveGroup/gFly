# Queues

This directory contains background tasks that are processed asynchronously by queue workers.

## Purpose

The queues directory is responsible for:
- Defining background tasks that can be queued for later processing
- Implementing asynchronous processing for time-consuming operations
- Offloading resource-intensive tasks from the request-response cycle
- Providing retry mechanisms for failed tasks
- Enabling distributed task processing

## Structure

- **hello_task.go**: Example queue task

## Usage

Queue tasks are defined as Go structs that implement the task interface:

```
// Example pseudocode for a queue task
// Note: This is not actual implementation code

// In a queue task file (email_task.go):
package queues

import (
    "github.com/gflydev/console"
    "github.com/gflydev/core/log"
)

// EmailTask sends emails asynchronously
type EmailTask struct {
    console.Task

    // Task parameters
    To      string
    Subject string
    Body    string
}

// Register the task with the queue system
func init() {
    console.RegisterTask(&EmailTask{})
}

// Handle processes the queued task
func (t *EmailTask) Handle() error {
    log.Infof("Sending email to %s with subject: %s", t.To, t.Subject)

    // Send the email
    err := sendEmail(t.To, t.Subject, t.Body)
    if err != nil {
        log.Errorf("Failed to send email: %v", err)
        return err // Task will be retried
    }

    log.Info("Email sent successfully")
    return nil
}
```

To dispatch a task to the queue:

```
// In application code:
emailTask := &queues.EmailTask{
    To:      "user@example.com",
    Subject: "Welcome to our service",
    Body:    "Thank you for signing up!",
}

// Dispatch the task to the queue
if err := console.Dispatch(emailTask); err != nil {
    log.Error(err)
}
```

To run the queue worker:

```bash
./build/artisan queue:run
```

## Best Practices

- Keep tasks focused on a single responsibility
- Make tasks idempotent (safe to run multiple times)
- Include all necessary data in the task struct
- Implement proper error handling and logging
- Use appropriate retry strategies for different types of failures
- Consider task priorities for critical operations
- Monitor queue depth and processing time
- Set appropriate timeouts for long-running tasks
- Document task parameters and expected behavior
- Test tasks in isolation before deploying
- Consider using separate queues for different types of tasks
- Implement graceful shutdown for queue workers

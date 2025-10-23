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

Queue tasks follow a specific pattern with separate payload and task structs:

### Task Structure

Each queue task consists of:
1. **Task Registration** - Automatic registration in `init()` with a unique identifier
2. **Constructor Function** - Returns payload and task identifier
3. **Payload Struct** - Contains the data to be processed
4. **Task Struct** - Implements the `Dequeue` method

```go
// Example based on hello_task.go implementation

package queues

import (
    "github.com/gflydev/console"
    "github.com/gflydev/core/errors"
    "github.com/gflydev/core/log"
)

// ---------------------------------------------------------------
// Register task
// ---------------------------------------------------------------

// Auto-register task into queue with a unique identifier
func init() {
    console.RegisterTask(&EmailTask{}, "send-email")
}

// ---------------------------------------------------------------
// Task info
// ---------------------------------------------------------------

// NewEmailTask Constructor for EmailTask
// Returns: (payload, taskIdentifier)
func NewEmailTask(to, subject, body string) (EmailTaskPayload, string) {
    return EmailTaskPayload{
        To:      to,
        Subject: subject,
        Body:    body,
    }, "send-email"
}

// EmailTaskPayload Task payload structure
type EmailTaskPayload struct {
    To      string
    Subject string
    Body    string
}

// EmailTask sends emails asynchronously
type EmailTask struct {
    console.Task
}

// Dequeue processes the queued task
func (t EmailTask) Dequeue(task *console.TaskPayload) error {
    // Decode task payload
    var payload EmailTaskPayload
    if err := task.BindPayload(&payload); err != nil {
        return errors.New("json.Unmarshal failed: %v: %s", err, task.GetType())
    }

    // Process payload
    log.Infof("Sending email to %s with subject: %s", payload.To, payload.Subject)

    // Send the email
    err := sendEmail(payload.To, payload.Subject, payload.Body)
    if err != nil {
        return errors.New("failed to send email: %v", err)
    }

    log.Info("Email sent successfully")
    return nil
}
```

### Dispatching Tasks

To dispatch a task to the queue, use the constructor function to create the payload:

```go
// In application code (e.g., in a service or controller):
import "your-project/internal/console/queues"

// Create task payload using constructor
payload, taskName := queues.NewEmailTask(
    "user@example.com",
    "Welcome to our service",
    "Thank you for signing up!",
)

// Dispatch the task to the queue
if err := console.DispatchTask(payload, taskName); err != nil {
    log.Error(err)
}
```

To run the queue worker:

```bash
./build/artisan queue:run
```

## Best Practices

### Task Design
- **Separate Concerns**: Use separate payload and task structs (e.g., `HelloTaskPayload` and `HelloTask`)
- **Constructor Pattern**: Create constructor functions that return `(payload, taskIdentifier)` tuple
- **Unique Identifiers**: Register tasks with unique string identifiers (e.g., `"send-email"`, `"hello-world"`)
- **Single Responsibility**: Keep tasks focused on one specific operation
- **Idempotent**: Make tasks safe to run multiple times with the same input

### Implementation
- **Payload Decoding**: Always use `task.BindPayload(&payload)` to decode the payload
- **Error Handling**: Use `github.com/gflydev/core/errors` package for creating errors with context
- **Error Messages**: Include task type in error messages using `task.GetType()`
- **Logging**: Use structured logging with `log.Infof()`, `log.Error()`, etc.
- **Method Signature**: Implement `Dequeue(task *console.TaskPayload) error` method

### Example Error Handling
```go
// Decode payload
if err := task.BindPayload(&payload); err != nil {
    return errors.New("json.Unmarshal failed: %v: %s", err, task.GetType())
}

// Business logic error
if err := someOperation(); err != nil {
    return errors.New("operation failed: %v", err)
}
```

### Operations
- **Retry Strategies**: Implement appropriate retry logic for different failure types
- **Task Priorities**: Consider priorities for critical operations
- **Monitoring**: Track queue depth and processing time
- **Timeouts**: Set appropriate timeouts for long-running tasks
- **Graceful Shutdown**: Implement proper shutdown for queue workers
- **Testing**: Test tasks in isolation before deploying
- **Separate Queues**: Consider using different queues for different task types

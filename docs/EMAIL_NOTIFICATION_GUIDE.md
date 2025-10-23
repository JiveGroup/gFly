# Email Notification Guide

This guide explains how to create and use email notifications in gFly using the official notification system.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Quick Start](#quick-start)
4. [Creating Notifications](#creating-notifications)
5. [Email Templates](#email-templates)
6. [Configuration](#configuration)
7. [Sending Notifications](#sending-notifications)
8. [Complete Examples](#complete-examples)
9. [Best Practices](#best-practices)
10. [Troubleshooting](#troubleshooting)

## Overview

gFly provides a unified notification system through the `github.com/gflydev/notification` package. The system separates notification logic from delivery mechanisms, allowing flexible implementation of email-based alerts.

### Key Features

- **Separation of Concerns**: Notification data is separate from delivery logic
- **Template-Based**: Uses Pongo2 templates for consistent email design
- **Type-Safe**: Strongly-typed notification structs
- **Multi-Channel**: Extensible to SMS, push notifications, etc.
- **Error Handling**: Built-in error management

## Architecture

### Components

```
┌─────────────────────────────────────────────────────────────┐
│                     Application Layer                       │
│              (Services, Controllers, Jobs)                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ notification.Send(notificationStruct)
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                 Notification Struct                         │
│          (Implements ToEmail() method)                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ Returns notifyMail.Data
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                  gFly Mail Notifier                         │
│            (github.com/gflydev/notification/mail)           │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ Sends via SMTP
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    Email Provider                           │
│              (SMTP Server / MailHog)                        │
└─────────────────────────────────────────────────────────────┘
```

### Directory Structure

```
project/
├── cmd/web/main.go                          # Registration: notificationMail.AutoRegister()
├── pkg/modules/auth/
│   ├── notifications/                       # Module-specific notifications
│   │   ├── verify_email_notification.go
│   │   ├── welcome_notification.go
│   │   ├── reset_password_notification.go
│   │   └── change_password_notification.go
│   └── services/
│       └── auth_services.go                 # Usage: notification.Send()
├── internal/notifications/                  # Global notifications
│   └── send_mail_notification.go
└── resources/views/mails/                   # Email templates
    ├── master.tpl                           # Base template
    ├── verify_email.tpl
    ├── welcome.tpl
    ├── forgot_password.tpl
    └── change_password.tpl
```

## Quick Start

### Step 1: Register Notification System

In `cmd/web/main.go` (already done in this project):

```go
package main

import (
    notificationMail "github.com/gflydev/notification/mail"
    // ... other imports
)

func main() {
    // Register mail notification
    notificationMail.AutoRegister()

    // ... rest of application setup
}
```

### Step 2: Create a Notification

Create `pkg/modules/auth/notifications/verify_email_notification.go`:

```go
package notifications

import (
    "github.com/gflydev/core"
    notifyMail "github.com/gflydev/notification/mail"
    view "github.com/gflydev/view/pongo"
)

type VerifyEmail struct {
    Email string
    Name  string
    Token string
}

func (n VerifyEmail) ToEmail() notifyMail.Data {
    body := view.New().Parse("mails/verify_email", core.Data{
        // For master template
        "title":    "Verify Your Email",
        "base_url": core.AppURL,
        "email":    n.Email,
        // For verify_email template
        "user_name": n.Name,
        "token":     n.Token,
    })

    return notifyMail.Data{
        To:      n.Email,
        Subject: "Verify Your Email Address",
        Body:    body,
    }
}
```

### Step 3: Create Email Template

Create `resources/views/mails/verify_email.tpl`:

```html
{% extends "master.tpl" %}
{% block body %}
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    Hi {{ user_name }}
</p>
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    Thank you for registering! Please verify your email address by clicking the button below.
</p>
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    <a href="{{ base_url }}/verify-email?token={{ token }}" target="_blank" style="border: solid 2px #0867ec; border-radius: 4px; box-sizing: border-box; cursor: pointer; display: inline-block; font-size: 16px; font-weight: bold; margin: 0; padding: 12px 24px; text-decoration: none; text-transform: capitalize; background-color: #0867ec; border-color: #0867ec; color: #ffffff;">
        Verify Email
    </a>
</p>
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    This link will expire in 24 hours.
</p>
{% endblock %}
```

### Step 4: Send Notification

In your service:

```go
package services

import (
    "gfly/pkg/modules/auth/notifications"
    "github.com/gflydev/notification"
    "github.com/gflydev/core/log"
)

func SendVerificationEmail(email, name, token string) error {
    if err := notification.Send(notifications.VerifyEmail{
        Email: email,
        Name:  name,
        Token: token,
    }); err != nil {
        log.Errorf("Failed to send verification email: %v", err)
        return err
    }
    return nil
}
```

## Creating Notifications

### Notification Struct

A notification struct contains all the data needed to compose the email:

```go
type MyNotification struct {
    // Required fields
    Email string  // Recipient email address

    // Optional fields based on your needs
    Name     string
    Token    string
    URL      string
    Amount   float64
    // ... any other data for the template
}
```

### Implementing ToEmail() Method

The `ToEmail()` method is required and must return `notifyMail.Data`:

```go
func (n MyNotification) ToEmail() notifyMail.Data {
    // 1. Parse template with data
    body := view.New().Parse("mails/my_template", core.Data{
        // Data for master.tpl
        "title":    "Email Title",
        "base_url": core.AppURL,
        "email":    n.Email,

        // Data for my_template.tpl
        "user_name": n.Name,
        "token":     n.Token,
        // ... other variables
    })

    // 2. Return mail data
    return notifyMail.Data{
        To:      n.Email,        // Required: recipient email
        Subject: "Email Subject", // Required: email subject
        Body:    body,           // Required: HTML body from template
    }
}
```

### Common Patterns

#### Pattern 1: Simple Notification

```go
type WelcomeEmail struct {
    Email string
    Name  string
}

func (n WelcomeEmail) ToEmail() notifyMail.Data {
    body := view.New().Parse("mails/welcome", core.Data{
        "title":     "Welcome!",
        "base_url":  core.AppURL,
        "email":     n.Email,
        "user_name": n.Name,
    })

    return notifyMail.Data{
        To:      n.Email,
        Subject: "Welcome to Our Platform",
        Body:    body,
    }
}
```

#### Pattern 2: Notification with Token/Link

```go
type ResetPassword struct {
    Email string
    Name  string
    Token string
}

func (n ResetPassword) ToEmail() notifyMail.Data {
    resetPasswordURI := utils.Getenv(constants.AuthResetPasswordUri, "/reset-password")

    body := view.New().Parse("mails/forgot_password", core.Data{
        "title":              "Reset Password",
        "base_url":           core.AppURL,
        "email":              n.Email,
        "user_name":          n.Name,
        "token":              n.Token,
        "reset_password_uri": resetPasswordURI,
    })

    return notifyMail.Data{
        To:      n.Email,
        Subject: "Reset Your Password",
        Body:    body,
    }
}
```

#### Pattern 3: Notification with Dynamic Data

```go
type OrderConfirmation struct {
    Email       string
    Name        string
    OrderID     int
    TotalAmount float64
    Items       []OrderItem
}

func (n OrderConfirmation) ToEmail() notifyMail.Data {
    body := view.New().Parse("mails/order_confirmation", core.Data{
        "title":        "Order Confirmation",
        "base_url":     core.AppURL,
        "email":        n.Email,
        "user_name":    n.Name,
        "order_id":     n.OrderID,
        "total_amount": n.TotalAmount,
        "items":        n.Items,
    })

    return notifyMail.Data{
        To:      n.Email,
        Subject: fmt.Sprintf("Order #%d Confirmed", n.OrderID),
        Body:    body,
    }
}
```

## Email Templates

### Template Structure

All email templates use Pongo2 syntax and extend the master template:

```html
{% extends "master.tpl" %}
{% block body %}
    <!-- Your email content here -->
{% endblock %}
```

### Master Template

The master template (`resources/views/mails/master.tpl`) provides:
- Responsive HTML email structure
- Professional styling
- Header/footer
- Mobile-friendly design
- Unsubscribe link

Available variables in master template:
- `{{ title }}` - Email title (for HTML `<title>` tag)
- `{{ base_url }}` - Application URL (from `core.AppURL`)
- `{{ email }}` - Recipient email (for unsubscribe link)

### Pongo2 Syntax

#### Variables

```html
<p>Hi {{ user_name }}</p>
<p>Your token is: {{ token }}</p>
```

#### Conditionals

```html
{% if user_name %}
    <p>Hi {{ user_name }}</p>
{% else %}
    <p>Hi there!</p>
{% endif %}
```

#### Loops

```html
<ul>
{% for item in items %}
    <li>{{ item.name }} - ${{ item.price }}</li>
{% endfor %}
</ul>
```

#### Filters

```html
<p>{{ user_name|upper }}</p>
<p>{{ description|truncate:100 }}</p>
<p>{{ created_at|date:"2006-01-02" }}</p>
```

### Template Best Practices

1. **Inline Styles**: Always use inline styles for email compatibility
2. **Simple Layout**: Keep layouts simple and table-based
3. **Alt Text**: Include alt text for images
4. **Test Rendering**: Test across email clients (Gmail, Outlook, etc.)
5. **Mobile First**: Design for mobile screens first
6. **Clear CTA**: Make call-to-action buttons prominent

### Example Template with Styling

```html
{% extends "master.tpl" %}
{% block body %}
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    Hi {{ user_name }}
</p>

<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    Thank you for your order! Your order details:
</p>

<table style="width: 100%; margin-bottom: 16px; border-collapse: collapse;">
    <tr style="background-color: #f4f5f6;">
        <th style="padding: 8px; text-align: left; border: 1px solid #eaebed;">Item</th>
        <th style="padding: 8px; text-align: right; border: 1px solid #eaebed;">Price</th>
    </tr>
    {% for item in items %}
    <tr>
        <td style="padding: 8px; border: 1px solid #eaebed;">{{ item.name }}</td>
        <td style="padding: 8px; text-align: right; border: 1px solid #eaebed;">${{ item.price }}</td>
    </tr>
    {% endfor %}
</table>

<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    <a href="{{ base_url }}/orders/{{ order_id }}" target="_blank" style="border: solid 2px #0867ec; border-radius: 4px; box-sizing: border-box; cursor: pointer; display: inline-block; font-size: 16px; font-weight: bold; margin: 0; padding: 12px 24px; text-decoration: none; text-transform: capitalize; background-color: #0867ec; border-color: #0867ec; color: #ffffff;">
        View Order
    </a>
</p>
{% endblock %}
```

## Configuration

### Environment Variables

Email configuration is handled through environment variables:

```env
# SMTP Configuration
MAIL_HOST=localhost
MAIL_PORT=1025
MAIL_USERNAME=
MAIL_PASSWORD=
MAIL_FROM_ADDRESS=noreply@example.com
MAIL_FROM_NAME="Your App Name"

# Application URL (used in email links)
APP_URL=http://localhost:7789
```

### Development with MailHog

For local development, use MailHog to catch emails:

```yaml
# docker-compose.yml
services:
  mailhog:
    image: mailhog/mailhog
    ports:
      - "1025:1025"  # SMTP
      - "8025:8025"  # Web UI
```

Configuration for MailHog:

```env
MAIL_HOST=localhost
MAIL_PORT=1025
MAIL_USERNAME=
MAIL_PASSWORD=
```

Access MailHog UI at `http://localhost:8025` to view sent emails.

### Production Configuration

For production, use a real SMTP provider (Gmail, SendGrid, AWS SES, etc.):

```env
# Example: Gmail SMTP
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USERNAME=your-email@gmail.com
MAIL_PASSWORD=your-app-password
MAIL_FROM_ADDRESS=noreply@yourdomain.com
MAIL_FROM_NAME="Your Company"

# Example: SendGrid SMTP
MAIL_HOST=smtp.sendgrid.net
MAIL_PORT=587
MAIL_USERNAME=apikey
MAIL_PASSWORD=your-sendgrid-api-key
MAIL_FROM_ADDRESS=noreply@yourdomain.com
MAIL_FROM_NAME="Your Company"
```

## Sending Notifications

### Basic Send

```go
import (
    "gfly/pkg/modules/auth/notifications"
    "github.com/gflydev/notification"
)

// Send notification
err := notification.Send(notifications.VerifyEmail{
    Email: "user@example.com",
    Name:  "John Doe",
    Token: "verification-token-123",
})
if err != nil {
    // Handle error
}
```

### With Error Handling

```go
if err := notification.Send(notifications.VerifyEmail{
    Email: user.Email,
    Name:  user.Fullname,
    Token: token,
}); err != nil {
    log.Errorf("Failed to send verification email to %s: %v", user.Email, err)
    return errors.New("Failed to send verification email")
}
```

### Non-Blocking Send (Best Practice)

For operations where email failure shouldn't block the main flow:

```go
// Send email but don't fail if it doesn't work
if err := notification.Send(notifications.Welcome{
    Email: user.Email,
    Name:  user.Fullname,
}); err != nil {
    // Log error but continue execution
    log.Errorf("Failed to send welcome email to %s: %v", user.Email, err)
    // Don't return error - user is still created successfully
}
```

### Queued Send (For Better Performance)

For high-volume applications, send emails asynchronously using queues:

```go
// Add to queue instead of sending immediately
queue.Push(job.SendEmail{
    NotificationType: "verify_email",
    Email:           user.Email,
    Name:            user.Fullname,
    Token:           token,
})
```

## Complete Examples

### Example 1: Email Verification Flow

**Notification Struct** (`pkg/modules/auth/notifications/verify_email_notification.go`):

```go
package notifications

import (
    "github.com/gflydev/core"
    notifyMail "github.com/gflydev/notification/mail"
    view "github.com/gflydev/view/pongo"
)

type VerifyEmail struct {
    Email string
    Name  string
    Token string
}

func (n VerifyEmail) ToEmail() notifyMail.Data {
    body := view.New().Parse("mails/verify_email", core.Data{
        "title":    "Verify Your Email",
        "base_url": core.AppURL,
        "email":    n.Email,
        "user_name": n.Name,
        "token":     n.Token,
    })

    return notifyMail.Data{
        To:      n.Email,
        Subject: "Verify Your Email Address",
        Body:    body,
    }
}
```

**Template** (`resources/views/mails/verify_email.tpl`):

```html
{% extends "master.tpl" %}
{% block body %}
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    Hi {{ user_name }}
</p>
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    Thank you for registering! Please verify your email address.
</p>
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    <a href="{{ base_url }}/verify-email?token={{ token }}" target="_blank" style="border: solid 2px #0867ec; border-radius: 4px; box-sizing: border-box; cursor: pointer; display: inline-block; font-size: 16px; font-weight: bold; margin: 0; padding: 12px 24px; text-decoration: none; text-transform: capitalize; background-color: #0867ec; border-color: #0867ec; color: #ffffff;">
        Verify Email
    </a>
</p>
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    Or copy and paste this link into your browser:<br/>
    <a href="{{ base_url }}/verify-email?token={{ token }}">{{ base_url }}/verify-email?token={{ token }}</a>
</p>
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    This link will expire in 24 hours. If you didn't create an account, you can safely ignore this email.
</p>
{% endblock %}
```

**Usage in Service** (`pkg/modules/auth/services/auth_services.go`):

```go
func SignUp(signUp dto.SignUp) (*models.User, error) {
    // ... user creation logic ...

    // Send verification email (non-blocking)
    if err := notification.Send(notifications.VerifyEmail{
        Email: user.Email,
        Name:  user.Fullname,
        Token: token,
    }); err != nil {
        log.Errorf("Failed to send verification email to %s: %v", user.Email, err)
        // Don't fail signup if email fails
    }

    return user, nil
}
```

### Example 2: Password Reset Flow

**Notification** (`pkg/modules/auth/notifications/reset_password_notification.go`):

```go
package notifications

import (
    "gfly/pkg/constants"
    "github.com/gflydev/core"
    "github.com/gflydev/core/utils"
    notifyMail "github.com/gflydev/notification/mail"
    view "github.com/gflydev/view/pongo"
)

type ResetPassword struct {
    ID    int
    Email string
    Name  string
    Token string
}

func (n ResetPassword) ToEmail() notifyMail.Data {
    resetPasswordURI := utils.Getenv(constants.AuthResetPasswordUri, "/reset-password")

    body := view.New().Parse("mails/forgot_password", core.Data{
        "title":              "Reset Password",
        "base_url":           core.AppURL,
        "email":              n.Email,
        "user_name":          n.Name,
        "token":              n.Token,
        "reset_password_uri": resetPasswordURI,
    })

    return notifyMail.Data{
        To:      n.Email,
        Subject: "Reset Your Password",
        Body:    body,
    }
}
```

**Usage** (`pkg/modules/auth/services/password_services.go`):

```go
func ForgotPassword(email string) error {
    user := repository.Pool.GetUserByEmail(email)
    if user == nil {
        return errors.New("User not found")
    }

    // Generate reset token
    token := generateResetToken()
    user.Token = dbNull.String(token)
    user.UpdatedAt = time.Now()

    if err := mb.UpdateModel(user); err != nil {
        return errors.New("Failed to generate reset token")
    }

    // Send reset email
    if err := notification.Send(notifications.ResetPassword{
        Email: user.Email,
        Name:  user.Fullname,
        Token: token,
    }); err != nil {
        log.Errorf("Failed to send reset password email: %v", err)
        return errors.New("Failed to send reset email")
    }

    return nil
}
```

### Example 3: Welcome Email

**Notification** (`pkg/modules/auth/notifications/welcome_notification.go`):

```go
package notifications

import (
    "github.com/gflydev/core"
    notifyMail "github.com/gflydev/notification/mail"
    view "github.com/gflydev/view/pongo"
)

type Welcome struct {
    Email string
    Name  string
}

func (n Welcome) ToEmail() notifyMail.Data {
    body := view.New().Parse("mails/welcome", core.Data{
        "title":     "Welcome!",
        "base_url":  core.AppURL,
        "email":     n.Email,
        "user_name": n.Name,
    })

    return notifyMail.Data{
        To:      n.Email,
        Subject: "Welcome to Our Platform",
        Body:    body,
    }
}
```

**Template** (`resources/views/mails/welcome.tpl`):

```html
{% extends "master.tpl" %}
{% block body %}
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    Hi {{ user_name }}
</p>
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    Welcome to our platform! Your email has been verified and your account is now active.
</p>
<p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
    <a href="{{ base_url }}/dashboard" target="_blank" style="border: solid 2px #0867ec; border-radius: 4px; box-sizing: border-box; cursor: pointer; display: inline-block; font-size: 16px; font-weight: bold; margin: 0; padding: 12px 24px; text-decoration: none; text-transform: capitalize; background-color: #0867ec; border-color: #0867ec; color: #ffffff;">
        Go to Dashboard
    </a>
</p>
{% endblock %}
```

**Usage**:

```go
// Send welcome email after successful verification (non-blocking)
_ = notification.Send(notifications.Welcome{
    Email: user.Email,
    Name:  user.Fullname,
})
```

## Best Practices

### 1. Error Handling

**Do**:
```go
if err := notification.Send(myNotification); err != nil {
    log.Errorf("Failed to send notification: %v", err)
    // Decide: return error or continue based on criticality
}
```

**Don't**:
```go
// Silent failure
notification.Send(myNotification)
```

### 2. Non-Critical Emails

For non-critical emails (welcome, confirmation), don't fail the operation:

```go
// User registration should succeed even if welcome email fails
if err := createUser(user); err != nil {
    return err
}

// Welcome email failure doesn't prevent registration
_ = notification.Send(notifications.Welcome{
    Email: user.Email,
    Name:  user.Fullname,
})
```

### 3. Template Organization

- Keep templates in `resources/views/mails/`
- Use descriptive names: `verify_email.tpl`, `order_confirmation.tpl`
- Extend master template for consistency
- Test templates across email clients

### 4. Notification Organization

- Module-specific: `pkg/modules/{module}/notifications/`
- Global: `internal/notifications/`
- One file per notification type
- Clear, descriptive struct names

### 5. Data Validation

Validate required fields before sending:

```go
if email == "" || token == "" {
    return errors.New("Email and token are required")
}

err := notification.Send(notifications.VerifyEmail{
    Email: email,
    Name:  name,
    Token: token,
})
```

### 6. Security

- Never expose sensitive data in email links (use tokens)
- Always use HTTPS for production links
- Set token expiry times
- Sanitize user input before including in templates

### 7. Testing

Test emails in development:
- Use MailHog for local testing
- Test with different email clients (Gmail, Outlook)
- Verify mobile responsiveness
- Check spam score

### 8. Performance

For high-volume emails:
- Use queue system for async sending
- Batch notifications when possible
- Monitor SMTP rate limits
- Implement retry logic

## Troubleshooting

### Emails Not Sending

**Check configuration**:
```go
// Verify environment variables
MAIL_HOST=localhost
MAIL_PORT=1025
MAIL_FROM_ADDRESS=noreply@example.com
```

**Check MailHog** (development):
- Ensure MailHog is running: `docker ps | grep mailhog`
- Check web UI: http://localhost:8025
- Verify SMTP port: 1025

**Check logs**:
```go
log.Errorf("Failed to send email: %v", err)
```

### Template Not Found

**Error**: `template not found: mails/my_template`

**Solution**: Verify template exists at correct path:
```
resources/views/mails/my_template.tpl
```

**Check template name** in code:
```go
view.New().Parse("mails/my_template", core.Data{...})
//                 ^^^^^^^^^^^^^^^^ must match filename without .tpl
```

### Template Rendering Issues

**Variable not showing**:
- Check variable is passed in `core.Data`
- Verify variable name matches in template

**Syntax error**:
- Check Pongo2 syntax
- Ensure `{% %}` and `{{ }}` are closed properly
- Validate `{% extends %}` and `{% block %}` structure

### SMTP Authentication Failed

**Gmail**:
- Use App Password, not regular password
- Enable "Less secure app access" (or use OAuth2)

**SendGrid**:
- Use "apikey" as username
- Use API key as password

### Emails Going to Spam

- Set up SPF, DKIM, DMARC records
- Use verified sender domain
- Avoid spam trigger words
- Include unsubscribe link
- Test with mail-tester.com

### Performance Issues

**Symptoms**: Slow API responses, timeouts

**Solutions**:
- Implement async email sending with queues
- Don't block HTTP responses waiting for email
- Use connection pooling
- Monitor SMTP server performance

## Additional Resources

- **gFly Documentation**: https://www.gfly.dev
- **Pongo2 Templates**: https://github.com/flosch/pongo2
- **Email Testing**: https://www.mail-tester.com
- **MailHog**: https://github.com/mailhog/MailHog

## Summary

The gFly notification system provides a clean, type-safe way to send emails:

1. **Create** notification struct with `ToEmail()` method
2. **Design** Pongo2 template extending master.tpl
3. **Send** using `notification.Send(notificationStruct)`
4. **Handle** errors appropriately
5. **Test** with MailHog in development

This pattern separates concerns, makes emails testable, and provides a consistent interface for all notification types.

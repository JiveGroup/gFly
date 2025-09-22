# Resources

This directory contains resources used by the application for the presentation layer and other non-public assets.

## Purpose

The resources directory is used for:
- View templates
- Email templates
- Localization files
- Frontend source files (before compilation)
- Configuration templates
- TLS certificates and keys
- Other non-public assets

## Structure

- **internal/**: Application-specific resources
- **views/**: View templates for rendering HTML
- **tls/**: TLS certificates and keys for HTTPS

## Usage

Resources in this directory are typically:
- Processed or compiled before being used
- Used by the application to generate dynamic content
- Not directly accessible to end users
- Used as source files for the build process

## Best Practices

- Organize templates in a logical directory structure
- Use a consistent naming convention for template files
- Keep sensitive information (like TLS keys) secure
- Use template inheritance or partials to avoid duplication
- Consider internationalization needs when designing templates
- Separate business logic from presentation concerns

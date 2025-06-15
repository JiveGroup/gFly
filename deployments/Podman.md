# Podman Configuration

This guide explains how to use [Podman](https://podman.io/) as an alternative to [Docker](https://www.docker.com/) for running the required services.

## Purpose

The docker directory is used for:
- Defining container services required by the application
- Providing configuration files for each service
- Enabling consistent development environments
- Simplifying the setup process for new developers

## Structure

- **docker-compose.yml**: Main Docker Compose configuration file
- **mailpit/**: Configuration files for the Mailpit service
  - **README.md**: Instructions for generating TLS certificates
  - **authfile**: Authentication configuration for Mailpit
  - **cert.pem**: TLS certificate for secure connections
  - **key.pem**: TLS private key for secure connections
- **redis/**: Configuration files for the Redis service
  - **redis.env**: Environment variables for Redis configuration

## Services

The Docker Compose configuration includes the following services:

### PostgreSQL Database (db)
- Image: postgres:16.4-alpine
- Port: 5432
- Credentials: user/secret
- Database: gfly

### Mail Testing (mail)
- Image: axllent/mailpit
- Ports: 8025 (web interface), 1025 (SMTP)
- Web interface: http://localhost:8025
- SMTP server: localhost:1025

### Redis (redis)
- Image: redis:7.4.0-alpine
- Port: 6379
- Password: secret (defined in redis.env)


## Config Podman

To use `Podman` instead of `Docker`, you need to modify the `Makefile` to replace command `docker compose` with `podman compose` in `Makefile`:

1. Open the `Makefile` in your preferred text editor
2. Find all occurrences of `docker compose` in the `container.run`, `container.logs`, `container.stop`, and `container.delete` targets
3. Replace each occurrence of `docker compose` with `podman compose`

For example, change:

```bash
docker compose --env-file deployments/docker/docker.env -f deployments/docker/docker-compose.yml -p gfly up -d db
```

To:

```bash
podman compose --env-file deployments/docker/docker.env -f deployments/docker/docker-compose.yml -p gfly up -d db
```

After making these changes, you can use the same Makefile commands as before:

```bash
make container.run   # Start all services with Podman
make container.stop  # Stop all services
```

## Usage

To start all services:

```bash
make container.run
```

To stop all services:

```bash
make container.stop
```

To view logs:

```bash
podman logs gfly-db    # Database logs
podman logs gfly-mail  # Mail logs
podman logs gfly-redis # Redis logs

# via Makefile command
make container.logs
```

## Best Practices

- Do not commit sensitive information in container configuration files
- Use environment variables for secrets and configuration
- Keep container images updated to the latest stable versions
- Use health checks to ensure services are properly initialized
- Document any special setup requirements for each service
- Use volumes for persistent data storage

# gFly v1.18.0

**Laravel inspired web framework written in Go**

Built on top of [FastHttp - the fastest HTTP engine](https://github.com/valyala/fasthttp), [FluentSQL - flexible and powerful SQL builder](https://github.com/jivegroup/fluentsql). Quick development with minimum memory allocation and high performance. Very simple and easy to use.

Note: You can find [documentation](https://doc.gfly.dev/) about [gFly](https://gfly.dev/)

# Tour of gFly

## I. Install environment

### 1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop/), [Podman Desktop](https://podman-desktop.io/) or [OrbStack](https://orbstack.dev/)

### 2. Install Golang

### 2.1 On Linux
```bash
# Install go at folder /home/$USER/Apps
mkdir -p /home/$USER/Apps
wget https://go.dev/dl/go1.26.0.linux-amd64.tar.gz
tar -xvzf go1.26.0.linux-amd64.tar.gz
```
Add `~/.profile` or `~/.zshrc` to the bottom of the file 
```bash
# Edit .zshrc
❯ vi ~/.zshrc

# ----------- Golang -----------
export GOROOT=/home/$USER/Apps/go
export GOPATH=/home/$USER/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
Check
```bash
❯ source ~/.profile
# Or
❯ source ~/.zshrc

# Check Go
go version
```

### 2.2 On Mac
```bash
# Install go at folder /Users/$USER/Apps
mkdir -p /Users/$USER/Apps
wget https://go.dev/dl/go1.26.0.darwin-arm64.tar.gz
tar -xvzf go1.26.0.darwin-arm64.tar.gz
```
Add `~/.profile` or `~/.zshrc` to the bottom of the file 
```bash
# Edit .profile
❯ vi ~/.profile

# ----------- Golang -----------
export GOROOT=/Users/$USER/Apps/go
export GOPATH=/Users/$USER/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
Check
```bash
❯ source ~/.profile
# Or
❯ source ~/.zshrc

# Check Go
go version
```

### 3. Install `Swag`, `Air`, `Migrate`, `Lint`, `GoSec`, `GoVulncheck`, and `GoCritic`

In Go programming, In addition to having an IDE that suits you, Other supporting tools will help control code quality better. Therefore, we recommend that you install and use the following tools as an indispensable part of the application development process with Go.

```bash
# ----- Install swag -----
go install github.com/swaggo/swag/cmd/swag@latest
swag -v

# ----- Install air -----
go install github.com/air-verse/air@latest
air -v

# ----- Install migrate -----
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate --version

# ----- Install Lint -----
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint --version
  
# ----- Option 2 -----
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.9.0
golangci-lint --version

# ----- Option 3 (Ubuntu) -----
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v2.9.0 && sudo chown $USER:$USER $(go env GOPATH)/bin/golangci-lint

# ----- Install GoSecure -----
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec version

# ----- Option 2 -----
curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.23.0
gosec version
  
# ----- Option 3 (Ubuntu) -----
curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v2.23.0 && sudo chown $USER:$USER $(go env GOPATH)/bin/gosec

# ----- Install Go vulncheck -----
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck -version

# ----- Install Critic -----
go install github.com/go-critic/go-critic/cmd/gocritic@latest
gocritic version
```

### 4. Create project skeleton from `gFly` repository
```bash
git clone https://github.com/jivegroup/gfly.git myweb
cd myweb 
rm -rf .git* && cp .env.example .env
```

## II. Start `redis`, `mail`, and `db` services and `application`

Make sure that you don't have any services running on localhost ports `6379`, `1025`, `8025`, and `5432` . 

### 1. Start docker services
```bash
# Docker run (Create DB, Redis, Mail services)
make container.run
```

Note: If you get any problems with `Docker`, review this tutorial [README.md](deployments/docker/README.md). 
Or you can check this document [Podman.md](deployments/Podman.md) if you are using `Podman` in your machine
### 2. Check services
```bash
# Docker
❯ docker ps

# Podman
❯ podman ps

>>> CONTAINER ID   IMAGE                  COMMAND                  CREATED         STATUS                   PORTS                                                                                            NAMES
>>> 38fb5bd004df   redis:7.4.0-alpine     "docker-entrypoint.s…"   9 minutes ago   Up 9 minutes             0.0.0.0:6379->6379/tcp, :::6379->6379/tcp                                                        gfly-redis
>>> 9e52bdb5a4ae   axllent/mailpit        "/mailpit"               9 minutes ago   Up 9 minutes (healthy)   0.0.0.0:1025->1025/tcp, :::1025->1025/tcp, 0.0.0.0:8025->8025/tcp, :::8025->8025/tcp, 1110/tcp   gfly-mail
>>> d62e30b0d548   postgres:16.4-alpine   "docker-entrypoint.s…"   9 minutes ago   Up 9 minutes (healthy)   0.0.0.0:5432->5432/tcp, :::5432->5432/tcp                                                        gfly-db
```

### 3. Start app
```bash
# Migrate DB
make migrate.up

# Generate Doc
make doc

# Run App
make dev
```

### 4. Check app

Browse URL [http://localhost:7789/](http://localhost:7789/)

Check API  via CLI
```
curl -v -X GET http://localhost:7789/api/v1/info | jq
```

Note: Install [jq](https://jqlang.github.io/jq/) tool to view JSON format

API doc [http://localhost:7789/docs/](http://localhost:7789/docs/)

### 5. CLI Actions

Run below commands in 3 different terminals

#### 5.1 Schedule 
```bash
# Run schedule (Terminal 1)
./build/artisan schedule:run
```

Note: Will get the message of `Schedule Job` file `hello_job.go` every 2 seconds

#### 5.2 Queue 
```bash
# Run queue (Terminal 2)
./build/artisan queue:run
```
Note: Nothing happens because we don't have any job that was queued!

#### 5.3 Command
```bash
# Run command `hello-world` (Terminal 3)
./build/artisan cmd:run hello-world
```

Note: Check the output of `Terminal 2` and `Terminal 3`. The `Terminal 2` has some output because the `Task` file `hello_task.go` was queued from `Command` file `hello_command.go` in `Terminal 3`

You can check more details about [command](https://doc.gfly.dev/docs/03-digging-deeper/03-01-02.command/), [schedule](https://doc.gfly.dev/docs/03-digging-deeper/03-01-03.schedule/), and [queue](https://doc.gfly.dev/docs/03-digging-deeper/03-01-04.queue/) at link [https://doc.gfly.dev/](https://doc.gfly.dev/)

**Important! Should run 2 commands `make schedule` and `make queue` to get full deployment environment.**

### 6. Build CLI and Web
```bash
make build
```

Check some binary files in folder build/

## III. Service connection

Add some code to check `application` connect to services `redis`, `mail`, and `db`.

### 1. Connect `Database` service

Make sure run successfully `make migrate.up` to initial tables

Note: Check DB connection and see 4 tables: `users`, `roles`, `user_roles`, and `address`.

#### Create command

Create a new command line `db-test`. Add file `internal/console/commands/db_command.go` 

```go
package commands

import (
    "gfly/internal/domain/models"
    "github.com/gflydev/console"
    "github.com/gflydev/core/log"
    mb "github.com/gflydev/db"
    "time"
)

// ---------------------------------------------------------------
//                      Register command.
// ./artisan cmd:run db-test
// ---------------------------------------------------------------

// Auto-register command.
func init() {
    console.RegisterCommand(&dbCommand{}, "db-test")
}

// ---------------------------------------------------------------
//                      DBCommand struct.
// ---------------------------------------------------------------

// DBCommand struct for hello command.
type dbCommand struct {
    console.Command
}

// Handle Process command.
func (c *dbCommand) Handle() {
    user, err := mb.GetModelBy[models.User]("email", "admin@gfly.dev")
    if err != nil || user == nil {
        log.Panic(err)
    }
    log.Infof("User %v\n", user)

    log.Infof("DBCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
```

#### Build and run command

```bash
# Build
make build

# Run command `db-test`
./build/artisan cmd:run db-test
```

### 2. Connect `Redis` service

Create a new command line `redis-test`. Add file `internal/console/commands/redis_command.go` 

```go
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

```

#### Build and run command

```bash
# Build
make build

# Run command `redis-test`
./build/artisan cmd:run redis-test
```

### 3. Connect `Mail` service

Create a new command line `mail-test`. Add file `internal/console/commands/mail_command.go` 

```go
package commands

import (
    "github.com/gflydev/cache"
    "github.com/gflydev/console"
    "github.com/gflydev/core/log"
    "time"
)

// ---------------------------------------------------------------
//                      Register command.
// ./artisan cmd:run mail-test
// ---------------------------------------------------------------

// Auto-register command.
func init() {
    console.RegisterCommand(&mailCommand{}, "mail-test")
}

// ---------------------------------------------------------------
//                      MailCommand struct.
// ---------------------------------------------------------------

// MailCommand struct for hello command.
type mailCommand struct {
    console.Command
}

// Handle Process command.
func (c *mailCommand) Handle() {
    // ============== Send mail ==============
    sendMail := notifications.SendMail{
        Email: "admin@gfly.dev",
    }

    if err := notification.Send(sendMail); err != nil {
        log.Error(err)
    }

    log.Infof("MailCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
```

#### Build and run command

```bash
# Build
make build

# Run command `mail-test`
./build/artisan cmd:run mail-test
```

Check mail at http://localhost:8025/

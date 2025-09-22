## Console `internal/console/`

Use CLI tool `./build/artisan` to run `Queue`, `Schedule`, `Command`.

gFly's console configuration options are stored in your application's `.env` configuration file. Make sure Redis server ready

```bash
REDIS_HOST="localhost"
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB_NUMBER=0
REDIS_QUEUE_NUMBER=1
```

**Note: Redis is required by gFly to run Queue, Session and Caching**


### Queue

While building your web application, you may have some `tasks`, such as parsing and storing an uploaded CSV file, that take too long to perform during a typical web request. Thankfully, gFly allows you to easily create queued `tasks` that may be processed in the background. By moving time intensive `tasks` to a queue, your application can respond to web requests with blazing speed and provide a better user experience to your customers.

gFly queues provide a unified queueing API across a single queue backend via Redis. Your `tasks` are defined in the `internal/console/queues` directory. To help you get started, a simple example `hello-world` task is defined within folder. Note: Run `./build/artisan cmd:run hello-world` for testing by pushing a task into queue.

The `task` will be activated and run asynchronously. Run below command to start queue handler:

    ./build/artisan queue:run

### Scheduler

gFly's scheduler offers a fresh approach to managing scheduled `jobs` on your server. The scheduler allows you to fluently and expressively define your command schedule within your gFly application itself. When using the scheduler, only a single cron entry is needed on your server. Your `job` schedule is defined in the `internal/console/schedules` directory. To help you get started, a simple example `hello-world run every 2 seconds` job is defined within folder.

Run below command to start scheduler:

    ./build/artisan schedule:run

### Command

Not only HTTP request to push data into your app. Sometimes you need more action from CLI. Artisan is the command-line interface included with gFly. It provides a number of helpful commands that can assist you while you build your application.

In addition to the commands provided with Artisan, you may also build your own custom commands. Commands are typically stored in the `internal/console/commands` directory. Command format

    ./build/artisan cmd:run <CMD> --<PARAM_NAME>=<PARAM_VALUE>

To help you get started, a simple example `hello-world` command is defined within folder. You can try it below command:

    ./build/artisan cmd:run hello-world

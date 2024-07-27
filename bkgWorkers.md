<h1>Adding Background Asynchronous workers to increase app performance</h1>
<h2>Inside the ./workers folder will be located to logic to integrate the workers</h2>

<h3>Distributor</h3>
<p>The distributor will be in charge of distribute tasks and enqueue them to be proccesed after</p>

<h3>Processor</h3>
<p>The processor will be in charge of process the tasks in the correct order</p>

<h2>Running Redis on our local machine</h2>
To test the Redis background asyncronous workers we need run Redis on our local machine.<br/>
We can do so using Docker to run Redix on [this](https://hub.docker.com/_/redis) DockerHub image

On the `Makefile` file, add the command to run the redis image on our machine with `Docker`

```Makefile
    redis:
        docker run --name redis -p 6379:6379 -d redis:7-alpine
```

We can check the connection with `Redis` with the next command:

```docker
    docker exec -it redis redis-cli ping
```

Add the redis address port to the environment variables:

```env
REDIS_ADDRESS=0.0.0.0:6379
```

<h2>Adding the workers to our app flow</h2>
<h3>Adding the distributor to the flow</h3>
On the `main` function we can create the distributor task and the redis Opt
```go
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	workers.NewReditTaskDistributor(redisOpt)
```

Let's also add some configuration to distribute the task to send the email

```go
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue("critical"),
	}
```

Add the logic on the `gapi server` to receive and user the `distributor` and the logic to distribute the `SendVerifyEmail` task after create a user

<h3>Adding the processor to the flow</h3>
For now, we are just adding the tasks to the queue, now we need a function that is responsible to run the async server and process those tasks

```go
    func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
        taskProcessor := workers.NewRedisTaskProcessor(redisOpt, store)
        log.Info().Msg("start task processor")
        err := taskProcessor.Start()
        if err != nil {
            log.Fatal().Err(err).Msg("failed to start task processor")
        }
    }
```

Now we can add the function on another go routine since his execution bloks the main thread of the application. So in order to run the processor to pull the tasks, the gateway server and the grpc server we need has something like this:

```go
	go runTaskProcessor(redisOpt, store)
	go runGatewayServer(config, store, distributor)
	rungRPCServer(config, store, distributor)
```

Add the configurations to take the critical tasks also

```go
	server := asynq.NewServer(
		redisOpt, asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
		},
	)
```

<h3>Adding error handlers</h3>
<h4>In order to handle the errors when the distributor or processor failed we need add te next code in to our configuration</h4>
```go
	server := asynq.NewServer(
		redisOpt, asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("process task failed")
			}),
		},
	)
```

We can also settup the default logs of the asynq server used for example to show when the server is running. The logic will be inside `./workers/logger.go`. Basically we are overriding the Loggers functions that are inside the asynq library.
And now we need add this logger configuration to the asynq server configuration.

```go
    Logger: NewLogger(),
```

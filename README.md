# Simple Queue Service

A simple queue that uses sqlite to store queued task.

## Status

I started this project to learn go and currently it comes with a TCP server which can be used to communicate with the queue service. I am still working on writing a client for it.

## Modules

- handlers: All the "commands" sent to the TCP server are forwarded to handlers to process the request.
- queue: Contains function which allows to save/consume a task
- server: Simple TCP server
- logs: contains server logs

## Usage

[Download and install Go](https://go.dev/dl/)

If make is installed in your system then run `make run` else run `go run .`

## Customization

You can use following env vars to customize server and queue's behaviour:

| Env Var         | Description                                                                                | Default Value |
|-----------------|--------------------------------------------------------------------------------------------|---------------|
| SQS_PORT        | Port number to which server should be binded                                               | 4500          |
| SQS_BUFFER_SIZE | Max amount of data per request, used by `Server.readData` to read specified amount of data | 4096          |
| SQS_TIMEOUT     | Read timeout in seconds                                                                    | 5             |
| SQS_VERBOSE     | Whether to print logs to stdout, 0 = false, any other value means true                     | 0             |

**Note: Support for config file is under progress**

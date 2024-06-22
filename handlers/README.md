# Handlers

All requests sent to the TCP server is forwarded here, all handlers are provided with `Request` and `Queue` pointers and are expected to return `Response` and error(optional).

## Actions

Each handler is mapped to an action, following are the supported actions:

### job

Fetches detail for the first queued job, the job returned is then dequeued from the queue and it's status is set to complete.

Example:

```json
{
    "Action": "job",
}
```

### jobs

Returns list of all the **queued** jobs

Example:
```json
{
	"Action": "jobs"
}
```

**Note: This handler is still in progress, in future it might have ability to filter jobs**

### add_job

Adds a new job in the queue, returns job id

Example:
```json
{
    "Action": "add_job",
    "Data": {
        "email": "john@doe.com"
    }
}
```

### get_job

Fetches detail for a single job based on it's id

Example:
```json
{
    "Action": "get_job",
    "Data": {
        "Id": "2589fa4f-ca98-4e86-9427-e9706ef0fa6a"
    }
}
```

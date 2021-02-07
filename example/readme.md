This is an example to demonstrate how the exporter may work.

We will eventually visualize this with jaeger, so lets start up the jaeger otlp example

```shell
docker run --rm -it \
  -p 16686:16686 \
  -p 55680:55680 \
  jaegertracing/opentelemetry-all-in-one 
```

Run the sample application to write to a file the demo metrics and traces

```go
$ go run ./app
$ ls
    metric.json
    trace.json
```

Then run 
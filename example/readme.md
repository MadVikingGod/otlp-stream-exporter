This is an example to demonstrate how the exporter may work.

We will eventually visualize this with jaeger, so lets start up the jaeger otlp example

```shell
docker run --rm -it \
  -p 16686:16686 \
  -p 55680:55680 \
  jaegertracing/opentelemetry-all-in-one 
```

Run the sample application to write to a file the demo metrics and traces

```shell
$ go run ./app
$ ls
...
trace.json
```

Star the jaeger otlp container
```shell
$ docker run -d --name jaeger -p 16686:16686 -p 55680:55680 jaegertracing/opentelemetry-all-in-one
```


Export the tracers to jaeger 
```shell
$ go run ./exporter
```

View the traces on jaeger [localhost:16686](http://localhost:16686/)
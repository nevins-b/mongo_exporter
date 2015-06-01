# Mongo Exporter

Export Mongo metrics to Prometheus.

To run it:

```bash
make
./mongo_exporter [flags]
```

### Flags

```bash
./mongo_exporter --help
```

* __`mongo.server`:__ Address (host and port) of the mongo instance we should
    connect to.
* __`web.listen-address`:__ Address to listen on for web interface and telemetry.
* __`web.telemetry-path`:__ Path under which to expose metrics.

---
title: "Sinks"
date: 2021-06-07T15:01:29+02:00
draft: false
---

Actually bruco supports the following event sinks:

1. Apache Kafka
2. NATS

## The kafka sink

The kafka sink is activated using the `kafka` kind

Example:
```yaml
# OPTIONAL: you can have a source and a processor without a sink
sink:
  kind: kafka
  brokers:
    - localhost:9092
  topic: test-out
  # OPTIONAL: target partition
  partition: 42
  # OPTIONAL: default hash if partition is not defined. values: manual, hash, random
  partitioner: hash

```

## The nats sink

The nats sink is activated using the `nats` kind

Example:
```yaml
# OPTIONAL: you can have a source and a processor without a sink
sink:
  kind: nats
  serverUrl: localhost:4222
  subject: out.sub
```


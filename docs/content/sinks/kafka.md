---
title: "Kafka"
weight: 1
draft: false
---

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
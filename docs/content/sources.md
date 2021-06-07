---
title: "Sources"
date: 2021-06-07T15:01:21+02:00
draft: false
---

Actually bruco supports the following event sources:

1. HTTP
2. Apache Kafka
3. NATS

## The HTTP event source

The http source is activated using the `http` kind

| yaml path | description |
| ----------- | ----------- |
| ignoreProcessorResponse | If true doen't return the prcessed value to the http caller |
| port | the port the http server listens too (do not change this if you are using bruco with k8s) |

Example:
```yaml
source:
  kind: http
  ignoreProcessorResponse: false
  # port: 8090
```

## The Kafka event source
The kafka source is activated using the `kafka` kind

| yaml path | description |
| ----------- | ----------- |
| brokers | a list of kafka brokers to connect to |
| topics | a list of topics to subscribe too |
| offset | the initial offset where the kafka consumer will start to consume |
| consumerGroup | the kafka consumer group |
| balanceStrategy | the kafka consumer balance startegy |
| ... | ... |

Example: 
```yaml
source:
  kind: kafka
  brokers:
    - localhost:9092
  topics:
    - test1
    - test2
  # OPTIONAL: default latest. values: latest, earliest
  offset: latest
  consumerGroup: my-consumer-group
  # OPTIONAL: default range. values: range, sticky, roundrobin
  balanceStrategy: range
  # NOTE: the async version (fireAndForget=true) will not guarantee
  # messages handling order between same partition
  fireAndForget: false
  # OPTIONAL: default 60
  # Increase this one if you have slow consumers to prevent
  # rebalance loop
  rebalanceTimeout: 120
  # OPTIONAL: default to 256. Put a low vaclue here in case of slow consumers
  # to prevent rebalancing loop. It can be as low as 0
  channelBufferSize: 256
  # OPTIONAL: default 1024 * 1024. Again for slow consumers you could keep this low
  fetchDefaultBytes: 8
```

## The NATS event source
The NATS source is activated using the `nats` kind

Example:
```yaml
source:
  kind: nats
  serverUrl: localhost:4222
  queueGroup: test
  subject: in.sub
```
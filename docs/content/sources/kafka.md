---
title: "Kafka"
weight: 1
draft: false
---

Bruco uses Sarama (https://github.com/Shopify/sarama) as kafka library.

The kafka source is activated using the `kafka` kind

| yaml path | description |
| ----------- | ----------- |
| brokers | a list of kafka brokers to connect to |
| topics | a list of topics to subscribe too |
| offset | the initial offset where the kafka consumer will start to consume |
| consumerGroup | the kafka consumer group |
| balanceStrategy | the kafka consumer balance startegy |
| fireAndForget | if true, the first available worker will process a message, regardless of partition message order |
| rebalanceTimeout | kafka  rebalanceTimeout |
| channelBufferSize | the size of consumer buffer |
| fetchDefaultBytes | the size in bytes to fetch from the topic each time |



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
#### Understand the balanceStrategy param
The balanceStrategy param, let's you config the partition assignment strategy per consumer.
There are 3 available values:
1. Range (default)
2. Stick
3. Round Robin

The **range** partition strategy, distributes partitions to consumers as ranges.
Example:
suppose you have 6 partitions and two consumers. The assignement strategy will be like:
```
  c1: [p0, p1, p2]
  c2: [p3, p4, p5]
```

The **sticky** partition strategy, will assign partitions to consumer trying to keep previous assignment while at the sime time obtaining a balanced partition distribution.
Example:
you have 6 partition and two consumer with an assignment like:
```
  c1: [p0, p2, p4]
  c2: [p1, p3, p5]
```
then a new consumer joins the consumer group. You could obtain a reassigment like:
```
  c1: [p0, p2]
  c2: [p1, p3]
  c3: [p4, p5]
```

The **roundrobin** partition strategy, uses a round robin parition distribution between consumers.
Example:
with 6 paritions and two consumers, you will get:
```
  c1: [p0, p2, p4]
  c2: [p1, p3, p5]
```
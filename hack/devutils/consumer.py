import time
from kafka import KafkaConsumer
from kafka.coordinator.assignors.roundrobin import RoundRobinPartitionAssignor

topic = "test"
kafka_broker = "localhost:9092"

consumer = KafkaConsumer(
    topic,
    bootstrap_servers=[kafka_broker],
    auto_offset_reset='latest',
    enable_auto_commit=True,
    group_id='python-consumer',
    max_poll_records=1,
    partition_assignment_strategy=[RoundRobinPartitionAssignor],
    value_deserializer=lambda x: x.decode('utf-8'))

for msg in consumer:
    partitions = consumer.assignment()
    l = []
    for p in partitions:
        l.append(p.partition)
    print(f"value: {msg.value}, partitions: {l}")
    time.sleep(12)
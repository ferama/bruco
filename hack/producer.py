from kafka import KafkaProducer
from kafka.admin import KafkaAdminClient, NewTopic, NewPartitions
import time

kafka_broker = "localhost:9092"
topic = "test"
num_partitions = 4

admin_client = KafkaAdminClient(
        bootstrap_servers = kafka_broker, 
        client_id = 'admin-client'
    )
try:
    topic_list = []
    topic_list.append(NewTopic(name=topic, num_partitions=num_partitions, replication_factor=1))
    admin_client.create_topics(new_topics=topic_list, validate_only=False)
except Exception as e:
    # if topic already exists, set partitions
    try:
        admin_client.create_partitions({topic: NewPartitions(num_partitions)})
    except:
        pass

producer = KafkaProducer(bootstrap_servers=[kafka_broker])

i = 0
while True:
    i += 1
    partition = i % num_partitions
    producer.send(topic, partition=partition, value=str(i).encode('utf-8'))
    print(f"Sent {i} on partition {partition}")
    # time.sleep(0.01)
from kafka import KafkaProducer
import time

kafka_broker = "localhost:9092"
topic = "test"

producer = KafkaProducer(bootstrap_servers=[kafka_broker])

i = 0
while True:
    i += 1
    producer.send(topic, value=str(i).encode('utf-8'))
    print(f"Sent {i}")
    time.sleep(1)
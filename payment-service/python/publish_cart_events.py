import stomp
import json
import random
from faker import Faker
import time

fake = Faker()

ACTIVEMQ_HOST = "localhost"
ACTIVEMQ_PORT = 61613
QUEUE_NAME = "/queue/cart_events"

def generate_cart_event():
    #event_type = random.choice(["add.cart", "remove.cart"])
    event_type = "add.cart"
    return {
        "event": event_type,
        #"studentId": str(random.randint(1000, 9999)),
        "studentId": "12345",
        "courseId": f"COURSE{random.randint(1, 100):03}",
        "courseName": fake.sentence(nb_words=3),
        "price": round(random.uniform(10, 500), 2) if event_type == "add.cart" else 0.0,
    }

# Connect to ActiveMQ
class ActiveMQPublisher:
    def __init__(self, host, port):
        self.conn = stomp.Connection([(host, port)])
        self.conn.connect(wait=True)

    def publish(self, queue_name, message):
        self.conn.send(body=message, destination=queue_name, content_type="application/json")

    def disconnect(self):
        self.conn.disconnect()

# Publish 100 random cart events
publisher = ActiveMQPublisher(ACTIVEMQ_HOST, ACTIVEMQ_PORT)

for _ in range(100):
    event = generate_cart_event()
    event_json = json.dumps(event)
    publisher.publish(QUEUE_NAME, event_json)
    print(f"Published event: {event}")
    time.sleep(10)

publisher.disconnect()
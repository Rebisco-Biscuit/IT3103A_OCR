import psycopg2
from faker import Faker
import random

conn = psycopg2.connect(
    dbname="payment-service",
    user="postgres",
    password="admin",
    host="localhost",
    port="5432"
)

cursor = conn.cursor()
fake = Faker()

def generate_card():
    return (
        fake.name(),
        fake.credit_card_number(),
        '1' + ''.join([str(random.randint(0, 9)) for _ in range(9)]),
        fake.credit_card_expire(),
        fake.credit_card_security_code(),
        round(random.uniform(0, 1000), 2),
        str(random.randint(1000, 9999)),
        random.choice([True, False])
    )

for _ in range(100):
    cursor.execute("""
        INSERT INTO mock_payments (card_holder, card_number, phone_number, expiry_date, cvv, balance, pin, is_valid)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
    """, generate_card())

conn.commit()
cursor.close()
conn.close()
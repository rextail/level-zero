import json
import random
import uuid
from datetime import datetime, timedelta

# Оригинальный JSON-шаблон
template = {
    "order_uid": "b563feb7b2b84b6test",
    "track_number": "WBILMTESTTRACK",
    "entry": "WBIL",
    "delivery": {
        "name": "Test Testov",
        "phone": "+9720000000",
        "zip": "2639809",
        "city": "Kiryat Mozkin",
        "address": "Ploshad Mira 15",
        "region": "Kraiot",
        "email": "test@gmail.com"
    },
    "payment": {
        "transaction": "b563feb7b2b84b6test",
        "request_id": "",
        "currency": "USD",
        "provider": "wbpay",
        "amount": 1817,
        "payment_dt": 1637907727,
        "bank": "alpha",
        "delivery_cost": 1500,
        "goods_total": 317,
        "custom_fee": 0
    },
    "items": [
        {
            "chrt_id": 9934930,
            "track_number": "WBILMTESTTRACK",
            "price": 453,
            "rid": "ab4219087a764ae0btest",
            "name": "Mascaras",
            "sale": 30,
            "size": "0",
            "total_price": 317,
            "nm_id": 2389212,
            "brand": "Vivienne Sabo",
            "status": 202
        }
    ],
    "locale": "en",
    "internal_signature": "",
    "customer_id": "test",
    "delivery_service": "meest",
    "shardkey": "9",
    "sm_id": 99,
    "date_created": "2021-11-26T06:22:19Z",
    "oof_shard": "1"
}

# Функция для генерации случайной даты и времени
def random_date(start, end):
    return start + timedelta(seconds=random.randint(0, int((end - start).total_seconds())))

# Генерация 100 JSON-объектов
data = []
for i in range(100):
    new_entry = template.copy()

    # Генерация уникальных идентификаторов
    new_entry["order_uid"] = str(uuid.uuid4())
    new_entry["payment"]["transaction"] = str(uuid.uuid4())
    new_entry["items"][0]["rid"] = str(uuid.uuid4())

    # Изменение времени создания заказа
    new_entry["date_created"] = random_date(datetime(2021, 1, 1), datetime(2023, 12, 31)).isoformat() + "Z"

    # Иногда намеренно создаем невалидную дату
    if random.random() < 0.1:
        new_entry["date_created"] = "invalid-date"

    data.append(new_entry)

# Сохранение сгенерированных данных в файл
with open('generated_data.json', 'w') as f:
    json.dump(data, f, indent=4)

print("Сгенерировано 100 JSON-объектов с уникальными ID и сохранено в файл 'generated_data.json'.")

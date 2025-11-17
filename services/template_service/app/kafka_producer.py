import json
import uuid
import logging

logger = logging.getLogger(__name__)

try:
    from kafka import KafkaProducer
    from django.conf import settings

    producer = KafkaProducer(
        bootstrap_servers=settings.KAFKA_BOOTSTRAP_SERVERS,
        value_serializer=lambda v: json.dumps(v).encode("utf-8"),
        key_serializer=lambda k: str(k).encode("utf-8") if k else None,
    )
    KAFKA_AVAILABLE = True
    logger.info("Kafka producer initialized successfully.")
except Exception as e:
    logger.warning(f"Kafka unavailable, using MockProducer: {e}")

    class MockProducer:
        def send(self, topic, value, key=None):
            logger.info(f"[MOCK-KAFKA] topic={topic}, key={key}")
            logger.info(json.dumps(value, indent=2))

        def flush(self):
            pass

    producer = MockProducer()
    KAFKA_AVAILABLE = False


def send_render_request(message: dict):
    """Send a message to Kafka, safely fallback if Kafka is unavailable."""
    message["correlation_id"] = str(uuid.uuid4())
    topic = settings.KAFKA_TEMPLATE_TOPIC if KAFKA_AVAILABLE else "mock-topic"

    try:
        producer.send(topic, value=message, key=message.get("template_id"))
        producer.flush()
        return True
    except Exception as e:
        logger.error(f"Failed to send Kafka message: {e}")
        return False

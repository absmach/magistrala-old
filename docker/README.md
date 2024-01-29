# Docker Composition

Configure environment variables and run Magistrala Docker Composition.

\*Note\*\*: `docker-compose` uses `.env` file to set all environment variables. Ensure that you run the command from the same location as .env file.

## Installation

Follow the [official documentation](https://docs.docker.com/compose/install/).

## Usage

Run the following commands from the project root directory.

```bash
docker-compose -f docker/docker-compose.yml up
```

```bash
docker-compose -f docker/addons/<path>/docker-compose.yml  up
```

To pull docker images from a specific release you need to change the value of `MG_RELEASE_TAG` in `.env` before running these commands.

## Broker Configuration

Magistrala supports configurable MQTT broker and Message broker, which also acts as an events store. Magistrala uses two types of brokers:

1. MQTT_BROKER: Handles MQTT communication between MQTT adapters and message broker. This can either be 'VerneMQ' or 'NATS'.
2. MESSAGE_BROKER: Manages message exchange between Magistrala core, optional, and external services. This can either be 'NATS' or 'RabbitMQ'. This is used to store messages for distributed processing.

Events store: This is the same as MESSAGE_BROKER. This can either be 'NATS' or 'RabbitMQ' or 'Redis'. This is used by Magistrala services to store events for distributed processing. If Redis is used as an events store, then RabbitMQ or NATS is used as a message broker.

The current deployment strategy for Magistrala in `docker/docker-compose.yml` is to use VerneMQ as a MQTT_BROKER and NATS as a MESSAGE_BROKER and EVENTS_STORE.

Therefore, the following combinations are possible:

- MQTT_BROKER: VerneMQ, MESSAGE_BROKER: NATS, EVENTS_STORE: NATS
- MQTT_BROKER: VerneMQ, MESSAGE_BROKER: NATS, EVENTS_STORE: Redis
- MQTT_BROKER: VerneMQ, MESSAGE_BROKER: RabbitMQ, EVENTS_STORE: RabbitMQ
- MQTT_BROKER: VerneMQ, MESSAGE_BROKER: RabbitMQ, EVENTS_STORE: Redis
- MQTT_BROKER: NATS, MESSAGE_BROKER: RabbitMQ, EVENTS_STORE: RabbitMQ
- MQTT_BROKER: NATS, MESSAGE_BROKER: RabbitMQ, EVENTS_STORE: Redis
- MQTT_BROKER: NATS, MESSAGE_BROKER: NATS, EVENTS_STORE: NATS
- MQTT_BROKER: NATS, MESSAGE_BROKER: NATS, EVENTS_STORE: Redis

For Message brokers other than NATS, you would need to build the docker images with RabbitMQ as the build tag and change the `docker/.env`. For example, to use RabbitMQ as a message broker:

```bash
MG_MESSAGE_BROKER_TYPE=rabbitmq make dockers
```

```env
MG_MESSAGE_BROKER_TYPE=rabbitmq
MG_MESSAGE_BROKER_URL=${MG_RABBITMQ_URL}
```

For Redis as an events store, you would need to run RabbitMQ or NATS as a message broker. For example, to use Redis as an events store with rabbitmq as a message broker:

```bash
MG_ES_TYPE=redis MG_MESSAGE_BROKER_TYPE=rabbitmq make dockers
```

```env
MG_MESSAGE_BROKER_TYPE=rabbitmq
MG_MESSAGE_BROKER_URL=${MG_RABBITMQ_URL}
MG_ES_TYPE=redis
MG_ES_URL=${MG_REDIS_URL}
```

For MQTT broker other than VerneMQ, you would need to change the `docker/.env`. For example, to use NATS as a MQTT broker:

```env
MG_MQTT_BROKER_TYPE=nats
MG_MQTT_BROKER_HEALTH_CHECK=${MG_NATS_HEALTH_CHECK}
MG_MQTT_ADAPTER_MQTT_QOS=${MG_NATS_MQTT_QOS}
MG_MQTT_ADAPTER_MQTT_TARGET_HOST=${MG_MQTT_BROKER_TYPE}
MG_MQTT_ADAPTER_MQTT_TARGET_PORT=1883
MG_MQTT_ADAPTER_MQTT_TARGET_HEALTH_CHECK=${MG_MQTT_BROKER_HEALTH_CHECK}
MG_MQTT_ADAPTER_WS_TARGET_HOST=${MG_MQTT_BROKER_TYPE}
MG_MQTT_ADAPTER_WS_TARGET_PORT=8080
MG_MQTT_ADAPTER_WS_TARGET_PATH=${MG_NATS_WS_TARGET_PATH}
```

### RabbitMQ configuration

```yaml
services:
  rabbitmq:
    image: rabbitmq:3.12.12-management-alpine
    container_name: magistrala-rabbitmq
    restart: on-failure
    environment:
      RABBITMQ_ERLANG_COOKIE: ${MG_RABBITMQ_COOKIE}
      RABBITMQ_DEFAULT_USER: ${MG_RABBITMQ_USER}
      RABBITMQ_DEFAULT_PASS: ${MG_RABBITMQ_PASS}
      RABBITMQ_DEFAULT_VHOST: ${MG_RABBITMQ_VHOST}
    ports:
      - ${MG_RABBITMQ_PORT}:${MG_RABBITMQ_PORT}
      - ${MG_RABBITMQ_HTTP_PORT}:${MG_RABBITMQ_HTTP_PORT}
    networks:
      - magistrala-base-net
```

### Redis configuration

```yaml
services:
  redis:
    image: redis:7.2.4-alpine
    container_name: magistrala-es-redis
    restart: on-failure
    networks:
      - magistrala-base-net
    volumes:
      - magistrala-broker-volume:/data
```

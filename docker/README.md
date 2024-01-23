# Docker Composition

Configure environment variables and run Magistrala Docker Composition.

\*Note\*\*: `docker-compose` uses `.env` file to set all environment variables. Ensure that you run the command from the same location as .env file.

## Installation

Follow the [official documentation](https://docs.docker.com/compose/install/).

## Usage

Run following commands from project root directory.

```bash
docker-compose -f docker/docker-compose.yml up
```

```bash
docker-compose -f docker/addons/<path>/docker-compose.yml  up
```

To pull docker images from a specific release you need to change the value of `MG_RELEASE_TAG` in `.env` before running these commands.

## Broker Configuration

Magistrala supports configurable MQTT broker and Message broker, which also acts as events store. Magistrala uses two types of brokers:

1. MQTT_BROKER: Handles MQTT communication between MQTT adapters and message broker. This can either be 'vernemq' or 'nats'.
2. MESSAGE_BROKER: Manages communication between adapters and Magistrala writer services. This can either be 'nats' or 'rabbitmq' or 'redis'. This is used to store messages for distributed processing.

Events store: This is the same as MESSAGE_BROKER. This can either be 'nats' or 'rabbitmq' or 'redis'. This is used by Magistrala services to store events for distributed processing. If redis is used as an events store, then rabbitmq or nats is used as a message broker.

Since nats is used as both MQTT_BROKER and MESSAGE_BROKER, it is not possible to run nats as an MQTT_BROKER and nats as a MESSAGE_BROKER at the same time, this is the current depolyment strategy for Magistrala in `docker/docker-compose.yml`.

. Therefore, the following combinations are possible:

- MQTT_BROKER: vernemq, MESSAGE_BROKER: nats, EVENTS_STORE: nats
- MQTT_BROKER: vernemq, MESSAGE_BROKER: nats, EVENTS_STORE: redis
- MQTT_BROKER: vernemq, MESSAGE_BROKER: rabbitmq, EVENTS_STORE: rabbitmq
- MQTT_BROKER: vernemq, MESSAGE_BROKER: rabbitmq, EVENTS_STORE: redis
- MQTT_BROKER: nats, MESSAGE_BROKER: rabbitmq, EVENTS_STORE: rabbitmq
- MQTT_BROKER: nats, MESSAGE_BROKER: rabbitmq, EVENTS_STORE: redis
- MQTT_BROKER: nats, MESSAGE_BROKER: nats, EVENTS_STORE: nats
- MQTT_BROKER: nats, MESSAGE_BROKER: nats, EVENTS_STORE: redis

For Message brokers other than nats, you would need to change the `docker/.env`. For example, to use rabbitmq as a message broker:

```env
MG_MESSAGE_BROKER_TYPE=rabbitmq
MG_MESSAGE_BROKER_URL=${MG_RABBITMQ_URL}
```

For redis as an events store, you would need to run rabbitmq or nats as a message broker. For example, to use redis as an events store with rabbitmq as a message broker:

```env
MG_MESSAGE_BROKER_TYPE=rabbitmq
MG_MESSAGE_BROKER_URL=${MG_RABBITMQ_URL}
MG_ES_TYPE=redis
MG_ES_URL=${MG_REDIS_URL}
```

For MQTT brokers other than nats, you would need to change the `docker/.env`. For example, to use vernemq as a MQTT broker:

```env
MG_MQTT_BROKER_TYPE=vernemq
MG_MQTT_BROKER_HEALTH_CHECK=${MG_VERNEMQ_HEALTH_CHECK}
MG_MQTT_ADAPTER_MQTT_QOS=${MG_VERNEMQ_MQTT_QOS}
MG_MQTT_ADAPTER_MQTT_TARGET_HOST=${MG_MQTT_BROKER_TYPE}
MG_MQTT_ADAPTER_MQTT_TARGET_PORT=1883
MG_MQTT_ADAPTER_MQTT_TARGET_HEALTH_CHECK=${MG_MQTT_BROKER_HEALTH_CHECK}
MG_MQTT_ADAPTER_WS_TARGET_HOST=${MG_MQTT_BROKER_TYPE}
MG_MQTT_ADAPTER_WS_TARGET_PORT=8080
MG_MQTT_ADAPTER_WS_TARGET_PATH=${MG_VERNEMQ_WS_TARGET_PATH}
```

### VerneMQ configuration

```yaml
services:
  vernemq:
    image: magistrala/vernemq:${MG_RELEASE_TAG}
    container_name: magistrala-vernemq
    restart: on-failure
    environment:
      DOCKER_VERNEMQ_ALLOW_ANONYMOUS: ${MG_DOCKER_VERNEMQ_ALLOW_ANONYMOUS}
      DOCKER_VERNEMQ_LOG__CONSOLE__LEVEL: ${MG_DOCKER_VERNEMQ_LOG__CONSOLE__LEVEL}
    networks:
      - magistrala-base-net
    volumes:
      - magistrala-broker-volume:/var/lib/vernemq
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

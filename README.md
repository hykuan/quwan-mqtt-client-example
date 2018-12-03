# quwan-mqtt-client-example

This example creates a mqtt client that subscribe to the given topic,

and will publish some example message to given topic.


## Getting Started

```bash
# you need to modify the environments in docker-compose.yaml
$ docker-compose up
```

## Run Locally

```bash
# run golang main.go locally. (go1.11 required)
export BROKER=ssl://[YOUR_NGINX_HOST]:8883
export PASSWORD=[YOUR_THING_TOKEN]
$ make run
```

## Dev site information
Please see the wiki page for the secret information

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable | Description                            | Default                      |
|----------|----------------------------------------|------------------------------|
| TOPIC    | channels/mainflux-channel-id/messages  | channels/1/messages          |
| BROKER   | mainflux-nginx-host:mainflux-mqtt-port | ssl://[YOUR_NGINX_HOST]:8883 |
| USER     | thing id                               | 1                            |
| PASSWORD | thing token                            | [YOUR_THING_TOKEN]           |

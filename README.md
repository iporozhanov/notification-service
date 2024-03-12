# Notification Service


## Table of Contents

- [Description](#description)
- [Usage](#usage)
  - [Running Locally](#running-locally)
  - [Running with Docker](#running-with-docker)
  - [Run tests](#run-tests)
  - [Authenticate](#authenticate)
- [ApiDoc](openapi.yaml)

## Description

This app receives notifications via HTTP endpoint and then sends the notification through the propper channel.

It uses PostgreSQL to guarantee delivery and retry mechanism.

Supported notification types are:
   
   * Email (type=1) - Supported providers: Mailgun
   * SMS (type=2) - Supported providers: Sinch
   * Slack (type=3)

To send a notification call the `POST /notifications` endpoint with the propper payload.

Example:

```curl
curl --location 'http://localhost:8080/notifications' \
--header 'AUTH_TOKEN: {KEY_FROM_CONFIG}' \
--header 'Content-Type: application/json' \
--data '{
    "identifier": "U06L2KFUMRD",
    "subject": "Test Subject",
    "message": "Hello, this is a test message",
    "type": 3
}'
```

[More information on the enpoint](./openapi.yaml)

If a notification send fails it will be retried until a configured number of attempts is reached.

### Extending functionality

#### Adding additional notification types

To add new notification types:
1.  add the new type to [/notification/notification.go](/notification/notification.go) and update the available types and `String()` method
2. update the config

#### Adding additional senders

You can add new senders to existing notification types (ex. Sendgid for Email)

1. Add the new sender type in [sender/sender.go](/sender/sender.go) and update the available types
2. Add the new sender implementing the `Sender interface` [sender package](/sender/)
3. Update the `NewSender()` method - [sender/sender.go](/sender/sender.go)

## Usage

### Running Locally
1. Setup  file

   ```shell
   cp config.yaml.example config.yaml
   ```
And fill it with the correct settings

2. Run it

   ```shell
   make run
   ```

### Running with Docker

1. Build the Docker image:

   ```shell
   make API_PORT=8080 build-docker
   ```

2. Run the Docker container with docker compose:

   ```shell
   make API_PORT=8080 run-docker
   ```
Notice: The api port you pass to the above commands should match the on in the config.yaml

### Run tests

   ```shell
   make test
   ```
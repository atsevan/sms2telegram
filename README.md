# SMS to Telegram

This project checks an [SMS Gammu Gateway](https://github.com/pajikos/sms-gammu-gateway) endpoint and sends a Telegram message for each new SMS it receives. It doesn’t store SMSs, doesn't rely on a relational database (RDBMS), and has a minimal list of dependencies.

## Prerequisites

- Go 1.20 or later
- Docker
- Docker Compose

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/sms2telegram.git
    cd sms2telegram
    ```

2. Install dependencies:

    ```sh
    go mod download || true
    ```

## Usage

### Running Locally

1. Set the required environment variables:

    ```sh
    export TELEGRAM_TOKEN=your_telegram_token      # Provided by BotFather
    export TELEGRAM_CHAT_ID=your_telegram_chat_id  # Forward the message to https://t.me/getidsbot
    export ENDPOINT=your_sms_endpoint              # http://sms-gammu-gateway:5000
    export USERNAME=your_username
    export PASSWORD=your_password
    export INTERVAL=10s
    ```

2. Run the application:

    ```sh
    go run main.go
    ```

### Running with Docker

1. Build the Docker image:

    ```sh
    docker build -t sms2telegram .
    ```

2. Run the Docker container:

    ```sh
    docker run -e TELEGRAM_TOKEN=your_telegram_token \
               -e TELEGRAM_CHAT_ID=your_telegram_chat_id \
               -e ENDPOINT=your_sms_endpoint \
               -e USERNAME=your_username \
               -e PASSWORD=your_password \
               -e INTERVAL=10s
    ```

### Running with Docker Compose

1. Create a `.env` file with the following content:

    ```env
    TELEGRAM_TOKEN=your_telegram_token
    TELEGRAM_CHAT_ID=your_telegram_chat_id
    ENDPOINT=your_sms_endpoint
    USERNAME=your_username
    PASSWORD=your_password
    INTERVAL=10s
    ```

2. Run the application using Docker Compose:

    ```sh
    docker-compose up --build
    ```

## Stopping the Application

To stop the application when running locally, press `Ctrl+C`.

To stop the Docker container, run:

```sh
docker-compose down
```


## docker-compose.yaml example with sms-gammu-gateway
```
version: '3'
services:
  sms-gammu-gateway:
    container_name: sms-gammu-gateway
    restart: on-failure
    image: pajikos/sms-gammu-gateway
    #    environment:
    #      - PIN="1234"
    ports:
      - "5000:5000"
    devices:
      # see [udev/rules.d/usb-modem.rules]
      - /dev/mobile:/dev/mobile
  sms2telegram:
    restart: unless-stopped
    image: atsevan/sms2telegram:latest
    depends_on:
      - sms-gammu-gateway
    # those should be set in the environment or in a .env file
    environment:
      - TELEGRAM_TOKEN=${TELEGRAM_TOKEN}
      - TELEGRAM_CHAT_ID=${TELEGRAM_CHAT_ID}
      - ENDPOINT=http://sms-gammu-gateway:5000
```
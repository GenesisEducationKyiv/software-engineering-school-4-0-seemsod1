# Exchange Rate Service

## Introduction

This project is a test case for Genesis Software School. The service provides APIs to:
1. Get the current USD to UAH exchange rate (`/rate`).
2. Subscribe to receive daily email updates on exchange rate changes (`/subscribe`).
3. emails with the current rate should be sent to all subscribed users once a day.

## Technologies

- Golang
- GORM (ORM for database)
- gocron (Scheduler)
- Docker
- Docker Compose
- other packages

## Installation

### Prerequisites

- Docker
- Docker Compose

### Steps

1. Clone the repository:
    ```sh
    git clone https://github.com/seemsod1/api-project.git
    cd api-project
    ```

2. Build and start the service using Docker:
    ```sh
    docker build -t api-project .
    docker-compose up --build
    ```

## Usage

### Endpoints

- **Get Exchange Rate**
    - URL: `http://localhost:8080/rate`
    - Method: GET
    - Response Codes:
        - 200: Success
        - 400: Bad Request

- **Subscribe for Updates**
    - URL: `http://localhost:8080/subscribe`
    - Method: POST
    - Response Codes:
        - 200: Success
        - 400: Bad Request (invalid email address)
        - 409: Conflict (email address already subscribed)
        - 500: Internal Server Error

## Configuration

### Scheduling and Email Sending

- Exchange rate emails are sent to all subscribed users once a day.
- The service adjusts sending times based on the user's timezone to reduce server load.
- A goroutine-based pool manages the email sending process efficiently.


## Past Version

### Initial Approach

The first version of the email sending mechanism retrieved all available subscribers and sent emails to them once a day at a specific server time. While this approach met the requirement of daily updates, it caused a high server load at the scheduled time each day.

### Improved Solution

To mitigate the server load issue, the sending mechanism was revised to distribute the email sending process throughout the day, relative to the users' time zones. For example, to ensure a subscriber in Kyiv (UTC+3) receives the email at 9 AM local time, the server (UTC) schedules the email to be sent at 6 AM.

### Implementation Details

- The revised solution uses goroutines to create a pool of email workers.
- Each worker pulls emails for sending at the appropriate time, based on the user's time zone.
- This distributed approach reduces the server load by spreading the email-sending tasks across different times.

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
- Taskfile

## Installation

### Prerequisites

- Docker
- Docker Compose

### Steps

1. Clone the repository:
    ```sh
    git clone https://github.com/seemsod1/api-project.git
   ```
    ```
    cd api-project
    ```
2. Build and start the service using Docker-Compose:
    ```sh
    task dcb
    ```
3. The service will be available at `http://localhost:8080`.
4. To stop the service, run:
    ```sh
    task dcd
    ```
   or 
    ```sh
    task docker-compose-down
    ```

### Executable file (Windows)

- To run the service without Docker, you can build the executable file:
    ```sh
    task build
    ```
- Build and run the executable:
    ```sh
    task run
    ```
Do not forget to set/change the environment variables before running the executable.


### Linters

- To run the linters, use the following command:
    ```sh
    task lint
    ```
  or 
    ```sh
    task l
    ```
  
### Tests

- To run the tests, use the following command:
    ```sh
    task test
    ```
  or 
    ```sh
    task t
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



# Metrics and Alerts

For monitoring and alerting, I use Victoria Metrics and Grafana to collect and visualize metrics. I define custom metrics for different parts of the service, such as API calls, subscriber activity, and email sending. By setting up alerts based on these metrics, in future we can proactively identify issues and take corrective actions.

P.S. All thresholds are arbitrary and should be adjusted based on the service's requirements and performance characteristics.
### 1. Coinbase,PrivatBank and NBU API Metrics

- `[provider]_fetched_success_total`: Counter for successful API calls to [provider].
- `[provider]_fetched_failed_total`: Counter for failed API calls to [provider].
- `[provider]_duration_seconds`: Summary of API call durations to [provider].

- We can also calculate:
  - Rate of successful API calls: `rate([provider]_fetched_success_total[5m])`
  - Rate of failed API calls: `rate([provider]_fetched_failed_total[5m])`
  - Average duration of API calls: `avg_over_time([provider]_duration_seconds_sum[5m]) / avg_over_time([provider]_duration_seconds_count[5m])`
  - Total number of API calls for [provider]: `[provider]_fetched_success_total + [provider]_fetched_failed_total`
  
**Alerts:**

- **High Failure Rate**: Trigger an alert if the rate of failed API calls exceeds a certain threshold.
  ```promql
  rate([provider]_fetched_failed_total[5m]) / rate([provider]_fetched_success_total[5m]) > 0.1
    ```
- **High Duration**: Trigger an alert if the average duration of API calls exceeds a threshold.
  ```promql
    avg_over_time([provider]_duration_seconds_sum[5m]) / avg_over_time([provider]_duration_seconds_count[5m]) > 2
    ```

### 2. Subscriber and Customer Metrics

- `new_subscribers_total`: Counter for the total number of new subscribers.
- `new_customers_total`: Counter for the total number of new customers.
- `unsubscribe_total`: Counter for the total number of unsubscribed users.
- `subscribe_status_conflict_total`: Counter for conflicts when subscribing.

- We can also calculate:
  - Rate of new subscribers: `rate(new_subscribers_total[5m])`
  - Rate of unsubscribed users: `rate(unsubscribe_success_total[5m])`
  - Rate of conflicts when subscribing: `rate(subscribe_status_conflict_total[5m])`
  
**Alerts:**

- **High Unsubscribe Rate**: Trigger an alert if the rate of unsubscribed users exceeds a certain threshold.
  ```promql
  rate(unsubscribe_success_total[5m]) / rate(new_subscribers_total[5m]) > 0.1
  ```
  
- **High Conflict Rate**: Trigger an alert if the rate of conflicts when subscribing exceeds a certain threshold in a given time window.
  ```promql
    rate(subscribe_status_conflict_total[1m]) > 100
    ```  
- **High New Customer Rate**: Trigger an alert if the rate of new customers exceeds a certain threshold.
  ```promql
    rate(new_subscribers_total[1m]) > 500
    ```
- **Small Number of New Subscribers**: Trigger an alert if the rate of new subscribers is below a certain threshold.
  ```promql
    rate(new_subscribers_total[1m]) < 10
    ```
### 3. Email Sending Metrics

- `emails_sent_successfully_total`: Counter for the total number of emails sent successfully.
- `emails_sent_failed_total`: Counter for the total number of failed email sending attempts.
- `email_sending_duration_seconds`: Summary of email sending durations.

- We can also calculate:
  - Rate of successful email sending attempts: `rate(emails_sent_successfully_total[5m])`
  - Rate of failed email sending attempts: `rate(emails_sent_failed_total[5m])`
  - Average duration of email sending: `avg_over_time(email_sending_duration_seconds_sum[5m]) / avg_over_time(email_sending_duration_seconds_count[5m])`
  - Total number of email sending attempts: `emails_sent_successfully_total + emails_sent_failed_total`

**Alerts:**

- **High Failure Rate**: Trigger an alert if the rate of failed email sending attempts exceeds a certain threshold.
  ```promql
  rate(emails_sent_failed_total[5m]) / rate(emails_sent_successfully_total[5m]) > 0.1
  ```
  
- **High Duration**: Trigger an alert if the average duration of email sending exceeds a threshold.
  ```promql
    avg_over_time(email_sending_duration_seconds_sum[5m]) / avg_over_time(email_sending_duration_seconds_count[5m]) > 2
    ```
- **High Number of Failed Email Sending Attempts**: Trigger an alert if the rate of failed email sending attempts exceeds a certain threshold.
    ```promql
        rate(emails_sent_failed_total[1m]) > 100
    ```
  
### 4. API Usage Metrics

- `api_calls_total`: The sum of requests_total{path="/api/v1/..."} metrics.

**Alerts:**

- **High API Usage**: Trigger an alert if the rate of API calls exceeds a certain threshold.
  ```promql
  rate(requests_total{path="/api/v1/rate"}[5m]) > 1000
  ```

### 5. Broker Metrics

- `broker_messages_produced_success_total`: The sum of messages produced by services. ex. `broker_messages_produced_total{service="emailstreamer"}` + ...
- `broker_messages_produced_failed_total`: The sum of messages failed by services. ex. `broker_messages_failed_total{service="emailstreamer"}` + ...
- `broker_messages_consumed_total`: The sum of messages consumed by services. ex. `broker_messages_consumed_total{service="sender"}` + ...


**Alerts:**

- **High data loss**: Trigger an alert if the rate of failed messages exceeds a certain threshold.
  ```promql
  rate(broker_messages_produced_failed_total[5m]) / rate(broker_messages_produced_success_total[5m]) > 0.1
  ```

- **High consumer lag**: Trigger an alert if difference between produced and consumed messages in time window high.

### etc
We can also add database metrics, system(like cpu and gpu usage). And add alerts on them.
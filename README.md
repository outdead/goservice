# goservice
Golang daemon skeleton.

# Config
<details>
<summary>config-example.yaml</summary>

```yaml
app:
  port: 8080
  profiler_port: 8099
  check_connections_interval: 10m
  error_buffer: 100
  log:
    level: info

connections:
  postgres:
    addr: "db_postgresql:5432"
    database: "test"
    username: "postgres"
    password: "postgres"
    debug: false
  clickhouse:
    addr: "db_clickhouse:9000"
    database: "test"
    debug: false
  rabbitmq:
    server:
      qos: 5000
      server: "amqp://test:tes@queue_rabbitmq:5672/"
      routing_key: "rk.test.income"
      queue:
        name: "test.income"
        durable: true
      exchange:
        name: "test"
        type: "direct"
        durable: true
    consumers:
      test_income:
        queue_name: "test.income"
        routing_key: "rk.test.income"
```
</details>

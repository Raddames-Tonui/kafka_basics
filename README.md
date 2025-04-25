kafka/
├── docker-compose.yml          # Kafka + PostgreSQL with KRaft mode
├── go.mod
├── go.sum
├── README.md
├── .env                        # Environment variables (e.g., Kafka port)
├── proto/                      # (Optional) for schema registry or protobuf later
├── config/
│   └── config.go               # Load config (e.g., DB creds, Kafka URL)
├── internal/
│   ├── producer/
│   │   ├── handler.go          # HTTP handler for /payments
│   │   └── producer.go         # Kafka producer logic
│   ├── consumer/
│   │   └── consumer.go         # Kafka consumer logic
│   ├── db/
│   │   ├── db.go               # PostgreSQL connection setup
│   │   └── insert.go           # Logic to insert events into DB
│   └── models/
│       └── payment_event.go    # PaymentEvent struct + JSON tags
├── cmd/
│   ├── producer/
│   │   └── main.go             # Starts HTTP server + publishes to Kafka
│   └── consumer/
│       └── main.go             # Starts Kafka consumer and writes to DB
└── scripts/
    └── init.sql                # SQL for creating payment_events table


docker exec -it kafka-postgres psql -U postgres -d payments_events
psql (15.12 (Debian 15.12-1.pgdg120+1))
Type "help" for help.

payments_events=# 


docker compose up --build -d

docker compose up

docker-compose logs -f kafka

go run cmd/producer/main.go

go run cmd/consumer/main.go

docker exec -it kafka /bin/bash

/opt/bitnami/kafka/bin/kafka-topics.sh --create --topic payment_events --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

/opt/bitnami/kafka/bin/kafka-topics.sh --list --bootstrap-server localhost:9092


 kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic payment_events --from-beginning

{"user_id":"12345","amount":100,"status":"success"}
{"user_id":"12345","amount":100,"status":"success"}

kafka-consumer-groups.sh --bootstrap-server localhost:9092 --group payment-consumer --reset-offsets --topic payment_events --to-earliest --execute
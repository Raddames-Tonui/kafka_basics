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


cmd/producer/main.go: Runs an HTTP server (e.g. POST /payments) that sends events to Kafka.

cmd/consumer/main.go: Kafka consumer that listens to payments topic and writes to PostgreSQL.

docker-compose.yml: Spins up Kafka (Bitnami or Confluent), PostgreSQL, and optionally Zookeeper.

scripts/init.sql: Initializes payment_events table.

internal/models/: Central place for event structs.

config/: Loads .env or hardcoded config for ports and credentials.


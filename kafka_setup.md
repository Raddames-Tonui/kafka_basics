docker exec -it kafka kafka-topics.sh --bootstrap-server localhost:9092 --create --topic payments --partitions 1 --replication-factor 1

docker exec -it kafka:
Run a command inside the running Docker container named kafka interactively (-it).

kafka-topics.sh:
Kafka’s built-in script for managing topics (create, delete, list, describe, etc.).

--bootstrap-server localhost:9092:
This tells the Kafka CLI which broker to connect to.
localhost:9092 refers to the Kafka server's listener on port 9092.

--create:
This flag tells Kafka to create a new topic.

--topic payments:
The name of the new topic is payments.

--partitions 1:
The topic will have 1 partition.
More partitions allow better scalability and parallelism.

--replication-factor 1:
Only 1 copy of the partition data will be maintained (i.e., no replication).
This is fine for local dev, but in production you'd want at least 2–3 for fault tolerance.

This command creates a Kafka topic called payments with 1 partition and no replication — perfect for testing locally.


See the topics
 docker exec -it kafka kafka-topics.sh --bootstrap-server localhost:9092 --list



 go mod init payment-tracker
go get github.com/gin-gonic/gin
go get github.com/segmentio/kafka-go
go get github.com/lib/pq

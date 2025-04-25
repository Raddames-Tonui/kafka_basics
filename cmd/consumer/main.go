package main

import (
	"kafka/internal/consumer"
	"kafka/internal/db"
)

func main (){
	dbConn := db.InitDB()
	defer dbConn.Close()

	consumer.StartConsumer("localhost:9092", "payment_events", "payment_group", dbConn)

}
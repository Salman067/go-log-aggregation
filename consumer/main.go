package main

import (
	consumer "log-aggregation/consumer/containers"
)

func main() {
	consumer.Serve()
}

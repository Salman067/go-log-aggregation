package main

import (
	producer "log-aggregation/producer/containers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	producer.Serve(e)
}

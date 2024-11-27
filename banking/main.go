package main

import (
	"banking/app"
	"github.com/wandz2810/banking-lib/logger"
)

func main() {
	logger.Info("Starting banking app....")
	app.Start()
}

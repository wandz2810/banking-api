#!/bin/bash
$env:SERVER_ADDRESS="localhost"
$env:SERVER_PORT="8080"
$env:MONGODB_URI="mongodb://localhost:27017"
go run main.go
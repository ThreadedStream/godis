#!/bin/sh
export POSTGRES_USER=$1
export POSTGRES_PASSWORD=$2
export POSTGRES_DB=$3
export POSTGRES_HOST=$4
export POSTGRES_PORT=$5

go run main.go client.go server.go server_utils.go queries.go json_models.go parser.go


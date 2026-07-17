#!/bin/sh

go test -coverprofile=coverage/cover.out ./...
go tool cover -html=coverage/cover.out

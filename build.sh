#!/bin/sh
set -x
rm -f pholcus_pkg/history/*
go build app/example.go && ./example

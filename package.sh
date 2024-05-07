#!/bin/bash
rm bootstrap*
go build -o bootstrap *.go
zip bootstrap.zip bootstrap
#!/bin/bash
export docker_host=192.168.31.128:12375
export GOOS=linux
export GOARCH=amd64
go build -o exam ../cmd

docker build . -t wuyuan_exam
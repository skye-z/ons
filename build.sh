#!/usr/bin/env bash

echo "Start build nas sync project..."

init(){
    go mod download
    go mod tidy
    rm -rf ./out
    mkdir ./out
}

generate(){
    CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -o $3_$1_$2 -ldflags '-s -w'
    mv $3_$1_$2 ./out/
}

cd cloud-server
init
echo "build cloud server from linux/x86_64"
generate linux amd64 cloud-server
echo "build cloud server from linux/arm64"
generate linux arm64 cloud-server

cd ../nas-server
init
echo "build nas server from linux/x86_64"
generate linux amd64 nas-server
echo "build nas server from linux/arm64"
generate linux arm64 nas-server
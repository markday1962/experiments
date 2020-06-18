#!/usr/bin/env bash
env GOOS=linux GOARCH=amd64 go build -o getArangoBackup -v main.go

echo "copying getArangoBackup to s3://aistemos-cloud-init/backup-checker/getArangoBackup"
aws s3 cp ./getArangoBackup s3://aistemos-cloud-init/backup-checker/getArangoBackup

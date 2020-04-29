#!/usr/bin/env bash
env GOOS=linux GOARCH=amd64 go build -o loadData -v main.go

echo "copying loadData to s3://aistemos-cloud-init/backup-checker/loadData"
aws s3 cp ./loadData s3://aistemos-cloud-init/backup-checker/loadData

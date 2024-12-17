#!/bin/bash


#GOOS=linux GOARCH=amd64 go build -o main .

docker build --platform=linux/amd64 -t ratelimiter:latest .
#
#aws-vault exec live -- aws ecr get-login-password --region eu-west-1 | docker login --username AWS --password-stdin  831753857064.dkr.ecr.eu-west-1.amazonaws.com
#
#
#docker tag ratelimiter:latest 831753857064.dkr.ecr.eu-west-1.amazonaws.com/ratelimiter:latest
#
#docker push 831753857064.dkr.ecr.eu-west-1.amazonaws.com/ratelimiter:latest

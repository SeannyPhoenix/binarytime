#!/bin/bash

ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
docker build -t "${ACCOUNT_ID}.dkr.ecr.us-west-2.amazonaws.com/binarytime-now" .
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin "${ACCOUNT_ID}.dkr.ecr.us-west-2.amazonaws.com"
docker push "${ACCOUNT_ID}.dkr.ecr.us-west-2.amazonaws.com/binarytime-now:latest"
aws lambda update-function-code --function-name binarytime-now --image-uri "${ACCOUNT_ID}.dkr.ecr.us-west-2.amazonaws.com/binarytime-now:latest" | cat

#!/bin/bash

export AWS_DEFAULT_REGION=ap-northeast-1
export DATABASE_URL=`aws ssm get-parameters --names fascia.$SERVICE_ENV.database_url --no-with-decryption --region $AWS_DEFAULT_REGION --query "Parameters[0].Value" --output text`
export DB_USER=`aws ssm get-parameters --names fascia.$SERVICE_ENV.db_user --no-with-decryption --region $AWS_DEFAULT_REGION --query "Parameters[0].Value" --output text`
export DB_NAME=`aws ssm get-parameters --names fascia.$SERVICE_ENV.db_name --no-with-decryption --region $AWS_DEFAULT_REGION --query "Parameters[0].Value" --output text`
export DB_PASSWORD=`aws ssm get-parameters --names fascia.$SERVICE_ENV.db_password --with-decryption --region $AWS_DEFAULT_REGION --query "Parameters[0].Value" --output text`
export CLIENT_ID=`aws ssm get-parameters --names fascia.$SERVICE_ENV.client_id --no-with-decryption --region $AWS_DEFAULT_REGION --query "Parameters[0].Value" --output text`
export CLIENT_SECRET=`aws ssm get-parameters --names fascia.$SERVICE_ENV.client_secret --with-decryption --region $AWS_DEFAULT_REGION --query "Parameters[0].Value" --output text`
export SLACK_URL=`aws ssm get-parameters --names fascia.$SERVICE_ENV.slack_url --no-with-decryption --region $AWS_DEFAULT_REGION --query "Parameters[0].Value" --output text`
export SECRET=`aws ssm get-parameters --names fascia.$SERVICE_ENV.secret --with-decryption --region $AWS_DEFAULT_REGION --query "Parameters[0].Value" --output text`

exec "$@"

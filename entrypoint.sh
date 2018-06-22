#!/bin/sh

export AWS_DEFAULT_REGION=ap-northeast-1
export DATABASE_URL=`myaws ssm parameter get fascia.$SERVICE_ENV.database_url --region $AWS_DEFAULT_REGION`
export DB_USER=`myaws ssm parameter get fascia.$SERVICE_ENV.db_user --region $AWS_DEFAULT_REGION`
export DB_NAME=`myaws ssm parameter get fascia.$SERVICE_ENV.db_name --region $AWS_DEFAULT_REGION`
export DB_PASSWORD=`myaws ssm parameter get fascia.$SERVICE_ENV.db_password --region $AWS_DEFAULT_REGION`
export CLIENT_ID=`myaws ssm parameter get fascia.$SERVICE_ENV.client_id --region $AWS_DEFAULT_REGION`
export CLIENT_SECRET=`myaws ssm parameter get fascia.$SERVICE_ENV.client_secret --region $AWS_DEFAULT_REGION`
export SLACK_URL=`myaws ssm parameter get fascia.$SERVICE_ENV.slack_url --region $AWS_DEFAULT_REGION`
export SECRET=`myaws ssm parameter get fascia.$SERVICE_ENV.secret --region $AWS_DEFAULT_REGION`

exec "$@"

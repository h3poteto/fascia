#!/bin/sh

export REGION=ap-northeast-1
export DATABASE_URL=`myaws ssm parameter get fascia.$SERVICE_ENV.database_url --region $REGION`
export DB_USER=`myaws ssm parameter get fascia.$SERVICE_ENV.db_user --region $REGION`
export DB_NAME=`myaws ssm parameter get fascia.$SERVICE_ENV.db_name --region $REGION`
export DB_PASSWORD=`myaws ssm parameter get fascia.$SERVICE_ENV.db_password --region $REGION`
export CLIENT_ID=`myaws ssm parameter get fascia.$SERVICE_ENV.client_id --region $REGION`
export CLIENT_SECRET=`myaws ssm parameter get fascia.$SERVICE_ENV.client_secret --region $REGION`
export SLACK_URL=`myaws ssm parameter get fascia.$SERVICE_ENV.slack_url --region $REGION`
export SECRET=`myaws ssm parameter get fascia.$SERVICE_ENV.secret --region $REGION`

exec "$@"

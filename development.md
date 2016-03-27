# Development
This document explain how to develop fascia.

## Environment Variables

I recommend to use `direnv` for environment variables, so please install.
For example,

```
export DATABASE_URL=$MYSQL_PORT_3306_TCP_ADDR  ## This will receive from mysql docker
export DB_USER="root"
export DB_PASSWORD="mysql"  ## This is specified by docker-compose.yml
export DB_NAME="fascia"
export DB_TEST_NAME="fascia_test"
export GOJIENV="development"
export CLIENT_ID="hogehoge"
export CLIENT_SECRET="fugafuga"
export TEST_TOKEN="testhoge"
export SLACK_URL="https://hooks.slack.com/services/hogehoge/fugafuga"
export AWS_ACCESS_KEY_ID=hogehoge   ## These will use AWS SES in mailer
export AWS_SECRET_ACCESS_KEY=fugafuga
export AWS_REGION=region
```

## Server Application

Development environment for fascia require Docker and Docker Compose, so you will need them.
Please install [Docker](https://docs.docker.com/mac/) and [Docker Compose](https://docs.docker.com/compose/).

Then, you can run docker container.

```
$ docker-compose run --rm --service-ports fascia /bin/bash
```

Please install dependent packages.

```
$ gom install
```


At first time, you need to create database, like this:

```
$ mysql -u root -p mysql -h $MYSQL_PORT_3306_TCP_ADDR
mysql > create database fascia;
```

And prepare database tables.

```
$ gom exec goose up
$ gom run db/seed/seed.go
```

After that, you can start server.

```
$ gom run server.go
```

Please open browser and access `localhost:9090`, you can access fascia on localhost.


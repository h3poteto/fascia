# Development
This document explain how to develop fascia.

## Environment Variables

Create file `.docker-env`, and write follwing environments:

```
DATABASE_URL=mysql                      ## MySQL docker host name
DB_USER=root
DB_PASSWORD=mysql                       ## This is specified by docker-compose.yml
DB_NAME=fascia
DB_TEST_NAME=fascia_test
CLIENT_ID=hogehoge                      ## GitHub application client id
CLIENT_SECRET=fugafuga                  ## GitHub application client secret key
TEST_TOKEN=testhoge                     ## GitHub access token for test environments
SLACK_URL=https://hooks.slack.com/services/hogehoge/fugafuga
AWS_ACCESS_KEY_ID=hogehoge              ## These will use AWS SES in mailer
AWS_SECRET_ACCESS_KEY=fugafuga
AWS_REGION=region
```

## Docker

Development environment for fascia require Docker and Docker Compose, so you will need them.
Please install [Docker](https://docs.docker.com/mac/) and [Docker Compose](https://docs.docker.com/compose/).


## JavaScript

It's necessary to prepare node packages, so please run npm install in docker container.

```
$ docker-compose run --rm node /bin/bash
node@b8446c2db58c:/var/opt/app$ npm install
```

## Server Application

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

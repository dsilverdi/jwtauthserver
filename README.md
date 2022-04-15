# JWT AUTH SERVER

Simple Quickstart/Boilerplate to Implement JWT Authentication Service

## Why I Develop This?
As a Boilerplate/Quickstart Server In case i need a fully customized JWT Authentication Service for My Next Project

## Depedencies
* Docker
* Docker Compose

## How to run
via terminal in this directory

* Build
```sh-session
$ docker-compose build
```

* run
```sh-session
$ docker-compose up
```

## Available API

For current version, jwtauthserver provides:

- Method `POST`: `/register` -  register user to the server
- Method `POST`: `/authorize` - returns jwt token (json)

# surl

Yet another URL shortener.

## Quickstart

You can bring up a development stack locally with Docker Compose:

```sh
$ sudo docker-compose up --build
```

This will build the Docker image for surl, and set up a node for you. It
will also provision a [dynalite] instance which mocks DynamoDB for local
development and testing.

To create a short URL for https://www.example.com:

```sh
$ curl -X POST -H "Content-Type: application/json" \
   -d '{"url": "https://www.example.com/"}' http://localhost:3000/submit
{"url":"https://www.example.com/","shorten_url":"M9Yv6VB2"}
```

From now on, http://localhost:3000/M9Yv6VB2 will redirect you there.

[dynalite]: https://github.com/mhart/dynalite

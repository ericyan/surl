# surl

surl is a DynamoDB-backed URL shortener written in Go. It is designed
to be high performance and operationally simple.

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

## Production deployment

surl is delivered as a single static-linked binary with no dependency
other than libc. A `Dockerfile` is also provided for building Docker
image -- the resulting image is only 26.2MB in size.

surl is 12factor-compliant: there are no configuration
files or command-line flags, everything is configured via environment
variables. See `docker-compose.yml` for details.

Amazon DynamoDB, which serves as the persistent storage, is the only
external dependency. Caches such as DynamoDB Accelerator, Memcached
and Redis are not needed, as surl has [bigcache], a very fast in-memory
cache embedded in.

Although a single instance of surl can handle thousands of requests per
second easily, it is highly recommended to deploy multiple instances
behind a load-balancer (Amazon ELB, HAProxy, Envoy, NGINX, etc) for the
sake of high availability.

surl itself is stateless so you can simply scale horizontally. Because
each surl instance has its own in-memory cache, you should configure the
load-balancer(s) to use consistent hashing or similar strategy to enable
sticky sessions, which will improve the cache hit rate significantly.

## Benchmarks

surl is designed to be fast. If the requested short URL is in its cache,
a single node can handle 100K req/s easily:

```sh
$ sudo docker run --rm --net=host williamyeh/wrk -t30 -c1000 -d30s --latency http://localhost:3000/M9Yv6VB2
Running 30s test @ http://localhost:3000/M9Yv6VB2
  30 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     9.56ms    3.57ms  70.05ms   77.48%
    Req/Sec     3.47k   571.76    35.01k    92.71%
  Latency Distribution
     50%    8.97ms
     75%   11.20ms
     90%   13.71ms
     99%   21.66ms
  3113729 requests in 30.10s, 671.10MB read
Requests/sec: 103446.90
Transfer/sec:     22.30MB
```

Even in the worst case senario that every request yields a cache miss,
in which case every request will hit the backend DynamoDB/dynalite, it
can still do nearly 2,000 req/s:

```sh
$ sudo docker run --rm --net=host williamyeh/wrk -t30 -c1000 -d30s --latency http://localhost:3000/notcached
Running 30s test @ http://localhost:3000/notcached
  30 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   298.21ms  314.81ms   2.00s    90.87%
    Req/Sec   113.62     95.98   350.00     57.71%
  Latency Distribution
     50%  179.19ms
     75%  346.08ms
     90%  571.73ms
     99%    1.70s
  54973 requests in 30.09s, 10.28MB read
  Socket errors: connect 0, read 0, write 0, timeout 2366
  Non-2xx or 3xx responses: 54973
Requests/sec:   1827.02
Transfer/sec:    349.70KB
```

The above benchmarks are from a single Debian Linux desktop of i7-7700
CPU and 16GB RAM. Just like any benchmarks, it should be interpreted
with care and YMMV.

[dynalite]: https://github.com/mhart/dynalite
[bigcache]: https://github.com/allegro/bigcache
[12factor]: https://12factor.net/

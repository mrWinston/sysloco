# sysloco

> People tend to overengineer stuff.  - Me

`sysloco` is a very very simple log accumulator. Think about Graylog or ELK.
With way less features. And currently limited to only syslog intake, a storage
solution that's as naïve as it gets and a simple webui that can't even draw nice
graphs.


## Motivation

So why even implement such a thing. The motivation behind sysloco is twofold:

* I wanted to improve on my golang and vuejs skillz
* I was frustrated with the available log accumulators

Let me explain (the second point):

Imagine you have a simple service running with docker-compose. No ominous
"WEB-SCALE" service that handles millions of requests per minute, but something
moderate, say 100 - 1000 req/min. You've deployed it with docker-compose, cause,
lets be honest, settings up all that k8s stuff for that little service is
overkill.

What if you just want to quickly check the logs of that service? Log into the
server with ssh, run `docker-compose logs`, maybe grep that output for errors.

What if you need s.o. else to check the logs? What if you want the logs to be
viewable without ssh access?

The standard answer is ELK. Just deploy the Whole elk-stack somewhere. Set up all
the log forwarders with maybe a sidecar container, configure the elasticserch
indicies, setup a dashboard on kibana and when you're done, you realise that
you're now paying double the infrastructure cost and spent two days just fiddling
around with logging.

This is where `sysloco` comes into play.

## Getting Started

Just add the following to each service in your docker-compose file:

```
logging:
  driver: syslog
  options:
    syslog-address: "udp://127.0.0.1:10001"
    syslog-format: rfc5424
    tag: <SERVICE NAME>
```

and run `sysloco` with the provided docker-compose file:

```
docker-compose up -d
```

Boom, log output on port 80!

![sysloco Screensho](./docs/screenshot.png)

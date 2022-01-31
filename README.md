# go-tokenbucket-ws
A tokenbucket service written in Go.

Since it's a webservice, it doesn't really matter that it's written in Go if you just want to use it.

# What's a token bucket?

The simple answer is that it's a widely used algorithm for rate limiting.
The wikipedia entry is pretty descriptive: https://en.wikipedia.org/wiki/Token_bucket

# Target use cases?

On the surface, it seems pretty silly to expose an algorithm as a webservice
that is already implemented in networking equipment.

However, there are additional use cases for rate limiting beyond networking.
Some examples:

* Applying some sort of restriction after too many failed login attempts
* A service offering that includes something like "x downloads per hour"

In the case of online offerings, there may be many servers offering services
in parallel, thus tracking these rates must be done centrally, and that's
what this service does.

# Usage

Once you have the program built, edit the config file and run the program.
It will listen on port 8080 by default. An example `wget` using the default
example config could be:

``wget -q -O - http://localhost:8080/example0/meow/5``

In this example, the config name is "example0" and the bucket is named "meow".

Configs must be configured in the config file. Attempts to access nonexistent configs
will result in an HTTP 404.

Buckets are created as-needed. Thus the trust model is that anyone that can
access the service can create and fetch data from buckets. i.e. the service
is not appropriate for use by untrusted sources.

The final component is the number of tokens requested. This must be an integer
no larger than a 64 bit int.

The GET will return a simple JSON response like this:

``{"Used":5,"Remaining":995}``

Where "Used" is the number of tokens used of those requested. This value should
be checked, as a request for 5 tokens is not guaranteed to have 5 tokens available
to use. In fact, that's the entire point of the service ... to tell another
service when tokens have run out.

"Remaining" is informational and tells how many tokens remain in the bucket after
using the number specified. This number will increase based on the configured
replenishment schedule.

Thus the recommended usage is to request tokens from a bucket and
take action as appropriate based on whether the tokens are available or not.

Bucket names can basically be any string. Thus buckets can be associated with
whatever type of ID you have for the service you're limiting (i.e. integers, UUIDS,
or anything else that is a string or can be converted to one)

# Configuration

The program accepts two parameters on startup.

``-l`` specifies the listen address and port, as would be appropriate to
pass into ``http.ListenAndServe()``. For example, ``-l localhost:8080`` will
cause the service to only respond to HTTP requests on 127.0.0.1 port 8080.

Also, you can specify the location of the config file with
the ``-c`` flag.

On startup, the service reads its configuration file. (Which defaults to
``tokenbucket.conf`` in the same directory as the binary) The service will
not start if it can not access this file.

This file is in CSV format, because I thought it would be funny to hear what
kind of comments I got about that.

The file format is simple. Each line represents a bucket configuration.
* The first column is the name.
* The second column is the initial number of tokens in a new bucket.
* The third column is the maximum number of tokens per bucket.
* The final column is the refresh interval. One token will be added to
the bucket (up to the maximum configured) on this interval. This can be
in any format that will work for Go's ``time.ParseDuration()`` (https://pkg.go.dev/time#ParseDuration)
For example: "1s", "15m", "1ms"

# Roadmap

_Things we want to do in the future. You can also think of this document as a "desired feature set" for ground._

Following is the set of features that ground is supposed to provide to its callers/dependents

## Web Server

We should be able to provide the basic Web Server (or Services) related features. This shall include: 

1. FastHTTP Server instance
2. Router for the FastHTTP server
3. Facility to add middleware
   - In future we want to be able to sub-mux our router so that middlewares could be grouped with APIs
   - Like how gorilla mux does it
4. Basic static file serving ability
5. Web Socket support?

## Databases

1. Provide interfaces which can be useful for setting up projects with a DB (most projects would need it)
2. SQLx _ lib/pq support for database - we start with PostgreSQL but since SQLx is DB agnostic, we can easily add support for MySQL and others later.
3. Code Generator to produce quality code that covers the general usecases around basic querying; support for indexes (indices) and constraints.


Some (most) projects are hardcore-dependent on a DB but some are not and there can be varying levels of needs around that. So we should allow an application to express one of the following demands:

1. Hardcore requirement - if the program does not connect to DB, it exits
2. Softcore requirement - If the program can't connect at first, it should keep retrying till it gets a connection
3. No requiement - disable the ability to connect to a DB.

And in addition, some projects might want to connect to multiple databases at once!

## Some good types and utility functions

Support for handling types

**NOTE**: This is not going to be super performant as a types library but as utilities, they would be handy

1. Strings
2. Integers
3. Floats
4. JsonObject

## Routine functions (Cron jobs)

A number of times, we need to process things periodically - either every few seconds or minutes (or some other duration) or by following a cron expression. We need to have support for such routines. 

1. A mini framework which allows you to create manageable routines which can be started, stopped, paused, resumed at any time (while ensuring that the atomicity of a single operation remains intact). 
2. Support for background this entire system of job management. 
3. support for making these routines follow a cron style execution (parsing a cron expression and running the function at the right time.)
4. The routines must be able to communicate in and out (provide channels for input and output for these routines)

## Support for caching
Some projects are hardcore-dependent on a cache but some are not and there can be varying levels of needs around that. So we should allow an application to express one of the following demands:

1. Hardcore requirement - if the program does not connect to cache, it exits
2. Softcore requirement - If the program can't connect at first, it should keep retrying till it gets a connection
3. No requiement - disable the ability to connect to a cache.

## Configuration management
There are multiple levels from which a configuration can come. In the decreasing order of importance, the sources are: 

1. Environment Variable
2. Config File
3. Default (hardcoded) value in code

## Logging
We will use bark. because it can: 

1. Allow you to log easily to STDOUT
2. Allow you to send your log messages to a remote server
3. (IN FUTURE) allow you to create a client which also hosts the DB and writes the logs to the database without having to start a separate server and then connect to it!

## Maintenance mode for the service (mostly for web server projects)
Sometimes you need to bring your services down - gracefully and in a controlled and planned manner. 

Gracefully - the requests that were already received are to be fulfilled. 
Controlled - No future requests are accepted
Planned - At a given point of time, the action (of stopping the acceptance of new requests and watiting to finish the remaining ones) starts

## Support for liveness and readiness probe
For kubernetes deployement. We can use the status and health check APis for that (see below)

## Support for status and health check APIs
- So that we can hit these endpoints and know that the server is up and running - I typically also include the UTC timestamp in response to these APIs to reflect if the service is actually running (and you are not getting a cached response)
- Health check APIs have to show the state of resources (is the DB connected, what is the number of requests waiting to be processed, is the cache engaged, what is the current memory footprint of the server and so on)

#### What does JsonObject do? (simple, rough example)
```json
{
   "some": "value",
   "key": "value2",
   "intval": 1234,
   "floatval": 12.34,
   "arrObjs": [
      {
         "k1": "v1",
         "k2": "v2"
      },
      {
         "name": "vaibhav",
         "where": "bangalore",
         "contact": {
            "country": "india",
            "state": "karnataka",
            "addresses": [
               {
                  "type": "office",
                  "region": "kr puram"
               },
               {
                  "type": "home",
                  "region": "electronic city"
               }
            ]
         }
      },
      1234,
      "some string"
   ],
   "obj": {
      "some": "thing"
   }
}
```

`arrObjs.[1].contact.addresses[0].region`

```go

var someJson map[string]any
```
if val, ok := someJson["arrObjs]; ok {
   if reflect.Type(val) == reflect.Array {
      arrVal := val.([]any)
      if len(arrVal) > 1 {
         firstElement := arrVal[0]
      }
   }
}



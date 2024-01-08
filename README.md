# Redis Server Lite
A simplified Redis server implementation written in Golang. This is my solution for [John Crickett's Write Your Own Redis Server Coding Challenge](https://codingchallenges.fyi/challenges/challenge-redis).

## Supported Commands
A list of the server's supported commands and their usage syntax.

### PING
Returns `PONG` if no argument is provided, otherwise return a copy of the argument as a bulk.
```
PING [message]
```

### ECHO
Returns `message`.
```
ECHO message
```

### SET
Set `key` to hold the string `value`. If `key` already holds a value, it is overwritten, regardless of its type.
```
SET key value [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds]
```

### GET
Returns the value of `key`. If the key does not exist the special value `nil` is returned.
```
GET key
```

### EXISTS
Returns the number of keys that exist from those specified as arguments.
```
EXISTS key [key ...]
```

## How It Works
The server listens for clients. Once a client connects, a go routine is created to handle the client's session. During a session the client can send commands with RESP (Redis serialization protocol) messages. The server parses a message and then sends a response in the same RESP message format.

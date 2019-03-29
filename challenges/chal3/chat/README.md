Simple IRC Server
=================
Implement a simple [IRC](https://en.wikipedia.org/wiki/Internet_Relay_Chat) server with a client-server architecture.


| Sub-command         | description                                             |
|---------------------|---------------------------------------------------------|
| `/users`            | To list all connected users                             |
| `/msg <user> <msg>` | To send a direct message to the specified `user`        |
| `/time`             | Get IRC Server's localtime                              |
| `/user <user>`      | Get more details about `user` (Username and IP address) |



To run Server
------------------------


```
# go run server.go -host localhost -port 9000

```

To run clients
---------------------------
```
# go run client.go -user user1 -server localhost:9000
```

```
# go run client.go -user user2 -server localhost:9000
```


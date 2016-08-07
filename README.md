# oauth2go
Just playground for oauth2 in go using https://github.com/RangelReale/osin


# How to run

## Create table(with database name = `mikamikuh`)
* Execute SQL: https://github.com/ory-am/osin-storage/blob/master/storage/postgres/postgres.go#L14

## Add client in the table
* `INSERT INTO client (id, secret, redirect_uri, extra) VALUES (1111, 'aabbccdd', 'http://localhost:14001/appauth', '');`

## Run client
```
$ go install github.com/mikamikuh/oauth2go/client
$ client
```
## Run server
```
$ go install github.com/mikamikuh/oauth2go/server
$ server
```

# How to play
* Access http://localhost:14001 (client)
* Click "Log in" button on the page
* Show authentication page on localhost:14000 (server)
* Click "Allow" -> redirect to localhost:14001/appauth with token
* Click "Deny" -> redirect to localhost:14001/appauth with error msg param

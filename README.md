# Golang ticket api

## Status

In progress

## Curl commands for available endpoints

### Create a new ticket option

```
curl -H 'Content-Type: application/json' \
      -d '{ "name":"Event","desc":"Event Description", "allocation": 10}' \
      -X POST \
      http://localhost:3000/ticket_options
```

### Get ticket option using ID

```
curl -H 'Content-Type: application/json' http://localhost:3000/ticket_options/38f13ced-8148-4d95-b5f6-28bf66f1eada
```

### Purchase n number of tickets from the ticket option

```
curl -H 'Content-Type: application/json' \
      -d '{ "quantity": 2,"user_id":"38f13ced-8148-4d95-b5f6-28bf66f1eadb"}' \
      -X POST \
      http://localhost:3000/ticket_options/38f13ced-8148-4d95-b5f6-28bf66f1eada/purchases
```

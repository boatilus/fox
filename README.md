# fox
[![Go Report Card](https://goreportcard.com/badge/github.com/boatilus/fox)](https://goreportcard.com/report/github.com/boatilus/fox) [![GoDoc](https://godoc.org/github.com/boatilus/fox?status.svg)](https://godoc.org/github.com/boatilus/fox)

A simple, dependency-free Go client for the Twilio programmatic fax API.

__fox__ seeks to implement all the functions associated with Twilio's "Faxes" endpoint, but to keep the library tight, does not facilitate, for example, E.164 phone number parsing and validation, or handling Twilio status callbacks.

## Getting started
To get started, construct a new `Client` with your Twilio account SID and auth token:

```go
c := fox.NewClient("YOUR_TWILIO_ACCOUNT_SID", "YOUR_TWILIO_AUTH_TOKEN")
```

Optionally, you can also pass a pointer to a `SendOptions` object to `NewClient` to specify custom
send options (to, for example, tell Twilio *not* to store fax media):

```go
opts := SendOpts{StoreMedia: false}
c := fox.NewClient("YOUR_TWILIO_ACCOUNT_SID", "YOUR_TWILIO_AUTH_TOKEN", &opts)
```

The `Cancel`, `Delete`, `Get`, `List` and `Send` methods on the returned `Client` are used to make the API calls as described by Twilio's API reference. For example, to retrieve a fax's data by its SID:

```go
res, _ := c.Get("FXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
```

## Implementation status
- ✅ Get a fax instance by its SID
- ✅ List all faxes instances in an account
- ✅ Send a fax
- ✅ Cancel (update) a fax by its SID
- ✅ Delete a fax instance by its SID
- ❌ Get a fax's media resource by its SID
- ❌ List all fax media resources in an account

## Running tests
First, grab the [testify](https://github.com/stretchr/testify) package:

    go get -u github.com/stretchr/testify
  
Then, run the tests (the `-v` flag specifies verbose output):

    go test -v

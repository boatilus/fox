# fox
A simple, dependency-free Go wrapper for the Twilio programmatic fax API.

__fox__ seeks to implement all the functions associated with Twilio's "Faxes" endpoint, but to keep the library tight, does not faciliate, for example, E.164 phone number parsing and validation, or handling Twilio status callbacks.

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

The `Get`, `List` and `Send` methods on the returned `Client` are used to make the API calls as described
by Twilio's API reference.

## Implementation status
- ✅ Get a fax by its SID
- ✅ List all faxes in an account
- ✅ Send a fax
- ❌ Get a fax resource by its SID
- ❌ List all fax resources in an account
- ❌ Create a new fax resource

## Running tests
First, grab the assert library:

    go get -u github.com/stretchr/testify/assert
  
Then, run the tests (the `-v` flag specifies verbose output):

    go test -v

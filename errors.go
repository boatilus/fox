package fox

import "errors"

var (
	// ErrNotAuthenticated indicates that the account SID and/or auth token are unspecified.
	ErrNotAuthenticated = errors.New("fox: account SID and/or auth token not specified")
	// ErrInvalidFaxNumber indicates that the fax number provided is invalid.
	ErrInvalidFaxNumber = errors.New("fox: fax number supplied is invalid")
	// ErrMissingSID indicates that a SID is required but was not supplied.
	ErrMissingSID = errors.New("fox: SID is required")
)

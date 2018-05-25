package fox

import "errors"

var (
	// ErrNotAuthenticated indicates that the account SID and/or auth token are unspecified.
	ErrNotAuthenticated = errors.New("fox: account SID and/or auth token not specified")
	// ErrInvalidFaxNumber indicates that the fax number provided is invalid.
	ErrInvalidFaxNumber = errors.New("fox: fax number supplied is invalid")
	// ErrMissingSID indicates that a SID is required but was not supplied.
	ErrMissingSID = errors.New("fox: SID is required")
	// ErrMissingToNumber indicates that a to number is required but was not supplied.
	ErrMissingToNumber = errors.New("fox: to number is required")
	// ErrMissingFromNumber indicates that a from number is required but was not supplied.
	ErrMissingFromNumber = errors.New("fox: from number is required")
	// ErrMissingMediaURL indicates that a media URL is required but was not supplied.
	ErrMissingMediaURL = errors.New("fox: media URL is required")
)

package fox

import "errors"

var (
	// ErrNotAuthenticated indicates that the AccountSID and/or AuthToken variables aren't set.
	ErrNotAuthenticated = errors.New("fox: account SID and/or auth token not specified")
	// ErrInvalidFaxNumber indicates that the fax number provided is invalid.
	ErrInvalidFaxNumber = errors.New("fox: fax number supplied is invalid")
)

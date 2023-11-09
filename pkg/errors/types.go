package errors

var (
	// ErrAuthentication indicates failure occurred while authenticating the entity.
	ErrAuthentication = New("failed to perform authentication over the entity")

	// ErrAuthorization indicates failure occurred while authorizing the entity.
	ErrAuthorization = New("failed to perform authorization over the entity")

	// ErrMalformedEntity indicates a malformed entity specification.

	//ErrUniqueID indicates an error in generating a unique ID
	ErrUniqueID = New("failed to generate unique identifier")

	//ErrFailedOpDB indicates a failure in a database operation
	ErrFailedOpDB = New("operation on db element failed")
)

package errors

var (
	// ErrAuthentication indicates failure occurred while authenticating the entity.
	ErrAuthentication = New("failed to perform authentication over the entity")

	// ErrAuthorization indicates failure occurred while authorizing the entity.
	ErrAuthorization = New("failed to perform authorization over the entity")

	//ErrUniqueID indicates an error in generating a unique ID
	ErrUniqueID = New("failed to generate unique identifier")

	//ErrFailedOpDB indicates a failure in a database operation
	ErrFailedOpDB = New("operation on db element failed")

	//ErrUnsupportedContentType indicates an invalid content type.
	ErrUnsupportedContentType = New("invalid content type")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = New("entity not found")

	// ErrConflict indicates that entity already exists.
	ErrConflict = New("entity already exists")

	// ErrMalformedEntity indicates a malformed entity specification.
	ErrMalformedEntity = New("malformed entity specification")
)

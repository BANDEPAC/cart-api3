package cart_errors

import "errors"

// Custom errors for handling specific scenarios.
var (
	ErrInvalidRequestMethod      = errors.New("invalid request method")
	ErrCartDoesNotExist          = errors.New("cart doesn't exist")
	ErrInvalidRequestBody        = errors.New("invalid request body")
	ErrInvalidQuery              = errors.New("invalid query")
	ErrCartIDRequired            = errors.New("cartID is required")
	ErrItemIDRequired            = errors.New("itemID is required")
	ErrNotFound                  = errors.New("failed to find a cart")
	ErrFailedToRetrieveCart      = errors.New("failed to retrieve cart")
	ErrFailedToRetrieveCartItems = errors.New("failed to retrieve cart items")
	ErrFailedPostgresOpperation  = errors.New("failed to make a query to postgres")
)

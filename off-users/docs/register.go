package docs

import (
	"github.com/juanjoss/off-users-service/ports"
)

// swagger:route POST /register post-requests endpoints
// Register adds a user and its devices.
// responses:
//   200: registerResponse
//	 default: registerErrorResponse

// Registration success.
// swagger:response registerResponse
type registrationResponseWrapper struct {
	// in:body
	Msg string
}

// Registration error.
// swagger:response registerErrorResponse
type registerErrorResponseWrapper struct {
	// in:body
	Msg string
}

// swagger:parameters endpoints
type registerRequestWrapper struct {
	// Send user data as an object and devices as an array of ssd objects.
	// in:body
	Body ports.RegisterRequest
}

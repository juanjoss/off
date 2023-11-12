package docs

import "github.com/juanjoss/off-orders-service/ports"

// swagger:route POST /orders post-requests endpoints
// Create a product order for a given product barcode and device id.
// responses:
//   200: createProductOrderResponse
//	 default: createProductOrderErrorResponse

// Product order created successfully.
// swagger:response createProductOrderResponse
type createProductOrderResponseWrapper struct {
	// in:body
	Msg string
}

// Error creating product order.
// swagger:response createProductOrderErrorResponse
type createProductOrderErrorResponseWrapper struct {
	// in:body
	Msg string
}

// swagger:parameters endpoints
type createProductOrderRequestWrapper struct {
	// Send the user's device id, product barcode and quantity.
	// in:body
	Body ports.CreateProductOrderRequest
}

package helpers

// NewRateResponse is a helper function that creates a new RateResponse struct
func NewRateResponse(price float64) interface{} {
	type RateResponse struct {
		Price float64 `json:"price"`
	}
	return RateResponse{Price: price}
}

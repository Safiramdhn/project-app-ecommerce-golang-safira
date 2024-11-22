package helper

import "math"

func CalculateDiscountPrice(price, discount float64) float64 {
	// Calculate the promo price
	priceAfterDiscount := price * ((100.00 - discount) / 100)

	// Round to 2 decimal places
	return math.Round(priceAfterDiscount*100) / 100
}

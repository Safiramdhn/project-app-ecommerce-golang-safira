package helper

import "math"

func CalculateDiscountPrice(price, discount float64) float64 {
	// Calculate the promo price
	priceAfterDiscount := price * ((100.00 - discount) / 100)

	// Round to 2 decimal places
	return math.Round(priceAfterDiscount*100) / 100
}

func CalculateCartPrice(price, additionalPrice, discount, promoDiscount float64, amount int) float64 {
	totalPrice := (price + additionalPrice) * ((100.00 - discount) / 100) * float64(amount)
	totalPrice -= totalPrice * ((100.00 - promoDiscount) / 100)

	return math.Round(totalPrice*100) / 100
}

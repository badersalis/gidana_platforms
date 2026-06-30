package config

// SeekerPlanPrices maps seeker plan → currency → price.
var SeekerPlanPrices = map[string]map[string]float64{
	"essential": {"XOF": 1500, "USD": 2.50},
	"pro":       {"XOF": 4000, "USD": 4.00},
}

// LandlordPlanPrices maps landlord plan → currency → price.
var LandlordPlanPrices = map[string]map[string]float64{
	"standard": {"XOF": 3000, "USD": 5.00},
	"agency":   {"XOF": 10000, "USD": 16.00},
}

// ListingLimits maps landlord plan → max active listings (-1 = unlimited).
var ListingLimits = map[string]int{
	"free":     1,
	"standard": 3,
	"agency":   -1,
}

func GetSeekerPrice(plan, currency string) (float64, bool) {
	prices, ok := SeekerPlanPrices[plan]
	if !ok {
		return 0, false
	}
	p, ok := prices[currency]
	return p, ok
}

func GetLandlordPrice(plan, currency string) (float64, bool) {
	prices, ok := LandlordPlanPrices[plan]
	if !ok {
		return 0, false
	}
	p, ok := prices[currency]
	return p, ok
}

func GetListingLimit(landlordPlan string) int {
	if limit, ok := ListingLimits[landlordPlan]; ok {
		return limit
	}
	return 1
}

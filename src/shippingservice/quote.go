// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"math"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/shippingservice/genproto"
)

// Quote represents a currency value.
type Quote struct {
	Dollars uint32
	Cents   uint32
}

// String representation of the Quote.
func (q Quote) String() string {
	return fmt.Sprintf("$%d.%02d", q.Dollars, q.Cents)
}

// CreateQuoteFromItems takes a number of items and returns a Price struct.
func CreateQuoteFromItems(items []*pb.CartItem) Quote {
	log.Info("[CreateQuoteFromItems] received request")
	defer log.Info("[CreateQuoteFromItems] completed request")

	for _, item := range items {
		quantity := item.Quantity
		if quantity < 1 {
			return Quote{Dollars: 0, Cents: 0}
		}
	}

	var totalItems int32
	for _, item := range items {
		totalItems += item.GetQuantity()
	}

	totalCost := 0.0
	if totalItems > 0 {
		// Flat $10 fee plus $0.50 per item
		totalCost = 10.0 + (float64(totalItems) * 0.5)
	}

	return CreateQuoteFromFloat(totalCost)
}

// CreateQuoteFromFloat takes a price represented as a float and creates a Price struct.
func CreateQuoteFromFloat(value float64) Quote {
	if value < 0 {
		// Decide how to handle negative input, maybe return zero or error
		return Quote{Dollars: 0, Cents: 0}
	}
	units, fraction := math.Modf(value)

	// Round cents to the nearest whole number
	roundedCentsFloat := math.Round(fraction * 100)

	// Calculate total cents including dollars
	totalCents := uint64(units*100) + uint64(roundedCentsFloat)

	// Extract final dollars and cents
	finalDollars := uint32(totalCents / 100)
	finalCents := uint32(totalCents % 100)

	return Quote{
		Dollars: finalDollars,
		Cents:   finalCents,
	}
}

// Server Handler
package handler

import (
	"strings"
)

var ServerOrders []string = []string{"banana", "apple", "orange", "grape", "red apple", "kiwi", "mango", "pear", "cherry", "green apple"}

func FindOrderByItemName(itemName string) (bool, []string) {
	var orders []string
	found := false
	for _, serverOrder := range ServerOrders {
		if strings.Contains(serverOrder, itemName) {
			orders = append(orders, serverOrder)
			found = true
		}
	}
	return found, orders
}

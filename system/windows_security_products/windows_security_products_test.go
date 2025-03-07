package windows_security_products

import (
	"fmt"
	"testing"
)

func TestGenWindowsSecurityProducts(t *testing.T) {
	products, err := GenWindowsSecurityProducts()
	if err != nil {
		t.Fatalf("Failed to get security products: %v", err)
	}

	// print products
	fmt.Println(products)
}

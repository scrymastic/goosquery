package windows_security_products

import (
	"fmt"
	"testing"

	"github.com/scrymastic/goosquery/sql/sqlctx"
)

func TestGenWindowsSecurityProducts(t *testing.T) {
	ctx := sqlctx.NewContext()
	products, err := GenWindowsSecurityProducts(ctx)
	if err != nil {
		t.Fatalf("Failed to get security products: %v", err)
	}

	// print products
	fmt.Println(products)
}

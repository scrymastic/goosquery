package windows_security_products

import (
	"testing"
)

func TestGenWindowsSecurityProducts(t *testing.T) {
	products, err := GenWindowsSecurityProducts()
	if err != nil {
		t.Fatalf("Failed to get security products: %v", err)
	}

	// Verify we got some products
	if len(products) == 0 {
		t.Error("Expected at least one security product, got none")
	}

	// Check each product has required fields
	for _, product := range products {
		if product.Type == "" {
			t.Error("Product type is empty")
		}
		if product.Name == "" {
			t.Error("Product name is empty")
		}
		if product.State == "" {
			t.Error("Product state is empty")
		}
		if product.RemediationPath == "" {
			t.Error("Product remediation path is empty")
		}
		if product.SignaturesUpToDate != 0 && product.SignaturesUpToDate != 1 {
			t.Errorf("Invalid signatures up to date value: %d", product.SignaturesUpToDate)
		}
	}
}

func TestSecurityProviderTypes(t *testing.T) {
	expectedTypes := []string{"Firewall", "Antivirus", "Antispyware"}
	for _, expectedType := range expectedTypes {
		found := false
		for _, actualType := range securityProviderTypes {
			if actualType == expectedType {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected security provider type %s not found", expectedType)
		}
	}
}

func TestSecurityProviderStates(t *testing.T) {
	expectedStates := []string{"On", "Off", "Snoozed", "Expired"}
	for _, expectedState := range expectedStates {
		found := false
		for _, actualState := range securityProviderStates {
			if actualState == expectedState {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected security provider state %s not found", expectedState)
		}
	}
}

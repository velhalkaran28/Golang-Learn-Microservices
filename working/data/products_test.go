package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "Karan",
		Price: 100,
		SKU:   "abc-abc-abc",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}

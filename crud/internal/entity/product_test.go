package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProdu(t *testing.T) {
	product, err := NewProduct("product1", 10)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, "product1", product.Name)
	assert.Equal(t, 10, product.Price)
}

func TestProductNameRequired(t *testing.T) {
	product, err := NewProduct("", 10)

	assert.NotNil(t, err)
	assert.Nil(t, product)
	assert.Equal(t, ErrNameRequired, err)
}

func TestProductPriceIsEquired(t *testing.T) {
	product, err := NewProduct("product1", 0)

	assert.NotNil(t, err)
	assert.Nil(t, product)
	assert.Equal(t, ErrPriceRequired, err)
}

func TestProductPriceInvalid(t *testing.T) {
	product, err := NewProduct("product1", -10)

	assert.NotNil(t, err)
	assert.Nil(t, product)
	assert.Equal(t, ErrInvalidPrice, err)
}

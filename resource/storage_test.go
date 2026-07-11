package resource_test

import (
	"testing"

	"github.com/bits-engine/bits-ecs/resource"
	"github.com/stretchr/testify/assert"
)

type RA struct {
	X int
	Y float64
}

type RB struct {
	X int
}

type RC struct {
	value string
}

func TestResourceStorage_Register(t *testing.T) {
	rs := resource.NewResourceStorage()
	RAT := resource.Register[RA](rs)
	RBT := resource.Register[RB](rs)

	assert.Equal(t, RAT.ID(), resource.ResourceID(0))
	assert.Equal(t, RBT.ID(), resource.ResourceID(1))
}

func TestResourceStorage_SetHasRemove(t *testing.T) {
	rs := resource.NewResourceStorage()
	RAT := resource.Register[RA](rs)
	RBT := resource.Register[RB](rs)
	RCT := resource.Register[RC](rs)

	assert.False(t, resource.Has(rs, RAT))
	assert.False(t, resource.Has(rs, RBT))
	assert.False(t, resource.Has(rs, RCT))

	resource.Set(rs, RAT, RA{X: 13})
	resource.Set(rs, RBT, RB{X: 13})

	assert.True(t, resource.Has(rs, RAT))
	assert.True(t, resource.Has(rs, RBT))
	assert.False(t, resource.Has(rs, RCT))

	resource.Remove(rs, RBT)

	assert.True(t, resource.Has(rs, RAT))
	assert.False(t, resource.Has(rs, RBT))
	assert.False(t, resource.Has(rs, RCT))
}

func TestResourceStorage_SetGetUpdate(t *testing.T) {
	rs := resource.NewResourceStorage()
	RAT := resource.Register[RA](rs)
	RBT := resource.Register[RB](rs)
	RCT := resource.Register[RC](rs)

	resource.Set(rs, RAT, RA{X: 13})
	resource.Set(rs, RBT, RB{X: 76})

	ra, exists := resource.Get(rs, RAT)
	assert.True(t, exists)
	assert.Equal(t, ra.X, 13)
	assert.Equal(t, ra.Y, 0.0)

	rb, exists := resource.Get(rs, RBT)
	assert.True(t, exists)
	assert.Equal(t, rb.X, 76)

	_, exists = resource.Get(rs, RCT)
	assert.False(t, exists)

	ra.X = 100
	ra.Y = 78.78

	ra, exists = resource.Get(rs, RAT)
	assert.True(t, exists)
	assert.Equal(t, ra.X, 100)
	assert.Equal(t, ra.Y, 78.78)
}

package bitset_test

import (
	"testing"

	"github.com/bits-engine/bits-ecs/common/bitset"
	"github.com/stretchr/testify/assert"
)

func TestBitSet_IsSet(t *testing.T) {
	bs := bitset.BitSet{}

	bs.Set(1)
	bs.Set(9)
	bs.Set(10)

	assert.True(t, bs.IsSet(1))
	assert.True(t, bs.IsSet(9))
	assert.True(t, bs.IsSet(10))
}

func TestBitSet_Clear(t *testing.T) {
	bs := bitset.BitSet{}

	bs.Set(1)
	bs.Set(9)
	bs.Set(10)

	assert.True(t, bs.IsSet(1))
	assert.True(t, bs.IsSet(9))
	assert.True(t, bs.IsSet(10))

	bs.Clear(9)
	assert.True(t, bs.IsSet(1))
	assert.False(t, bs.IsSet(9))
	assert.True(t, bs.IsSet(10))
}

func TestBitSet_Key(t *testing.T) {
	bs1 := bitset.BitSet{}

	bs1.Set(1)
	bs1.Set(243)
	bs1.Set(10)

	key1 := bs1.Key()

	bs2 := bitset.BitSet{}

	bs2.Set(1)
	bs2.Set(243)
	bs2.Set(10)

	key2 := bs1.Key()

	bs3 := bitset.BitSet{}

	bs3.Set(1)
	bs3.Set(243)
	bs3.Set(11)

	key3 := bs1.Key()

	assert.Equal(t, key1, bs1.Key())
	assert.Equal(t, key1, key2)
	assert.Equal(t, bs2.Key(), bs1.Key())
	assert.NotEqual(t, key1, key3)
	assert.NotEqual(t, key2, key3)
	assert.NotEqual(t, key2, bs3.Key())
}

func TestBitSet_Clone(t *testing.T) {
	bs := bitset.BitSet{}

	bs.Set(1)
	bs.Set(243)
	bs.Set(10)

	cloned := bs.Clone()

	assert.True(t, bs.HasAll(cloned))

	cloned.Set(2)
	assert.False(t, bs.HasAll(cloned))
}

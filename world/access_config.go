package world

import (
	"github.com/bits-engine/bits-ecs/common/bitset"
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/resource"
)

type AccessConfig struct {
	compRead  []component.ComponentID
	compWrite []component.ComponentID
	resRead   []resource.ResourceID
	resWrite  []resource.ResourceID
	exclusive bool
}

type componentIdentifier interface {
	ID() component.ComponentID
}

type resourceIdentifier interface {
	ID() resource.ResourceID
}

func (ac *AccessConfig) Reads(cid componentIdentifier) *AccessConfig {
	ac.compRead = append(ac.compRead, cid.ID())
	return ac
}

func (ac *AccessConfig) Writes(cid componentIdentifier) *AccessConfig {
	ac.compWrite = append(ac.compWrite, cid.ID())
	return ac
}

func (ac *AccessConfig) ReadsRes(rid resourceIdentifier) *AccessConfig {
	ac.resRead = append(ac.resRead, rid.ID())
	return ac
}

func (ac *AccessConfig) WritesRes(rid resourceIdentifier) *AccessConfig {
	ac.resWrite = append(ac.resWrite, rid.ID())
	return ac
}

func (ac *AccessConfig) Exclusive(isExclusive bool) *AccessConfig {
	ac.exclusive = isExclusive
	return ac
}

func (ac *AccessConfig) Compile() *compiledAccessConfig {
	reads := &bitset.BitSet{}
	writes := &bitset.BitSet{}

	// building reads and writes set, where odd bits - is for resources, and even bits - for components
	for _, compID := range ac.compRead {
		reads.Set(uint64(compID) * 2)
	}

	for _, compID := range ac.compWrite {
		writes.Set(uint64(compID) * 2)
	}

	for _, resID := range ac.resRead {
		reads.Set(uint64(resID)*2 + 1)
	}

	for _, resID := range ac.resWrite {
		writes.Set(uint64(resID)*2 + 1)
	}

	return &compiledAccessConfig{
		reads:     reads,
		writes:    writes,
		exclusive: ac.exclusive,
	}
}

type compiledAccessConfig struct {
	reads     *bitset.BitSet
	writes    *bitset.BitSet
	exclusive bool
}

func (c *compiledAccessConfig) Conflicts(o *compiledAccessConfig) bool {
	if c.exclusive || o.exclusive {
		return true
	}

	return c.writes.HasAny(o.writes) || c.writes.HasAny(o.reads) || c.reads.HasAny(o.writes)
}

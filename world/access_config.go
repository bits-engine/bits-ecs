package world

import (
	"github.com/bits-engine/bits-ecs/common/bitset"
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/resource"
)

// AccessConfig defines list of accessed components and resources by system.
//
// It is important to define correct list of accessed resources/components, because this information is used by scheduler to schedule
// systems to run in parallel locking-free.
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

// Reads defines read access to component
func (ac *AccessConfig) Reads(cid componentIdentifier) *AccessConfig {
	ac.compRead = append(ac.compRead, cid.ID())
	return ac
}

// Writes defines write access to component
func (ac *AccessConfig) Writes(cid componentIdentifier) *AccessConfig {
	ac.compWrite = append(ac.compWrite, cid.ID())
	return ac
}

// ReadsRes defines read access to resource
func (ac *AccessConfig) ReadsRes(rid resourceIdentifier) *AccessConfig {
	ac.resRead = append(ac.resRead, rid.ID())
	return ac
}

// WritesRes defines write access to resource
func (ac *AccessConfig) WritesRes(rid resourceIdentifier) *AccessConfig {
	ac.resWrite = append(ac.resWrite, rid.ID())
	return ac
}

// Exclusive defines that system will be run with exclusive access to world.
//
// Use only when you creating/deleting entire entities, or manipulating scheduling / world.
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

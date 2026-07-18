package world

import (
	"fmt"
	"slices"

	"github.com/bits-engine/bits-ecs/component"
)

// FilterRegistry incapsulates access to component.ComponentStorage inside Access() method of system
type FilterRegistry interface {
	Register(ac *AccessConfig, f *component.Filter) component.FilterID
}

type csFilterRegistry struct {
	cs *component.ComponentStorage
}

// Register registrates new [Filter] and checks access config for allowance of used components.
//
// Returned [FilterID] is used to actually query with this filter.
func (fr *csFilterRegistry) Register(
	ac *AccessConfig,
	f *component.Filter,
) component.FilterID {
	for _, compID := range f.Required() {
		if !slices.Contains(ac.compWrite, compID) && !slices.Contains(ac.compRead, compID) && !ac.exclusive {
			panic(fmt.Sprint("filter can not contain componentID", compID, "because it is not in AccessConfig set"))
		}
	}

	return component.RegisterFilter(fr.cs, f)
}

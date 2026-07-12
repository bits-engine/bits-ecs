package world

import (
	"fmt"
	"slices"

	"github.com/bits-engine/bits-ecs/component"
)

type FilterRegistry interface {
	Register(ac *AccessConfig, f *component.Filter) component.FilterID
}

type csFilterRegistry struct {
	cs *component.ComponentStorage
}

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

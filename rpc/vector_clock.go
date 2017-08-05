package rpc

import "github.com/la0rg/test_tasks/vector_clock"

// Go transforms proto struct to original go struct
func (vc *VC) Go() *vector_clock.VC {
	return vector_clock.NewVcWithExistingStore(vc.GetStore())
}

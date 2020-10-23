package accumulator

import (
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/cuebernetes/cuebectl/pkg/ensure"
	"github.com/cuebernetes/cuebectl/pkg/identity"
)

// LocationAccumulator tracks locations (in a cluster) of objects that have been created / synced from an instance.
type LocationAccumulator struct {
	ensurer ensure.Interface

	// locators to lookup values that have been synced at least once. sync.Map because it is a grow-only cache
	locators sync.Map
}

func NewLocationAccumulator(ensurer ensure.Interface) *LocationAccumulator {
	return &LocationAccumulator{
		ensurer:  ensurer,
		locators: sync.Map{},
	}
}

// Sync attempts to create an unstructured object identified by []path in instance.
// if successful, it returns a locator that can be used to lookup the object in the cluster later.
func (a *LocationAccumulator) Sync(obj *unstructured.Unstructured, path ...string) (*identity.Locator, error) {
	_, locator, err := a.ensurer.EnsureUnstructured(obj)
	if err != nil {
		return nil, err
	}
	locator.Path = path

	a.locators.Store(strings.Join(path, "."), &locator)

	return &locator, nil
}

// Locators returns the list of locators for concrete values
func (a *LocationAccumulator) Locators() (locators []*identity.Locator) {
	locators = make([]*identity.Locator, 0)
	a.locators.Range(func(key, value interface{}) bool {
		l := value.(*identity.Locator)
		locators = append(locators, l)
		return true
	})
	return
}
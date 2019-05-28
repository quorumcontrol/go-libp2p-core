package metrics

import (
	"fmt"
	"sync"

	"go.opencensus.io/stats/view"
)

var registeredViews = map[string][]*view.View{}
var mutex = sync.Mutex{}

type ErrNamespace struct {
	Namespace string
}

// ErrUnregisteredNamespace is an error for lookup of requested unregistered Namespace
type ErrUnregisteredNamespace ErrNamespace

func (e ErrUnregisteredNamespace) Error() string {
	return fmt.Sprintf("no views found registered under Namespace %s", e.Namespace)
}

// ErrDuplicateNamespaceRegistration is an error for a Namespace that has already
// registered views
type ErrDuplicateNamespaceRegistration ErrNamespace

func (e ErrDuplicateNamespaceRegistration) Error() string {
	return fmt.Sprintf("duplicate registration of views by Namespace %s", e.Namespace)
}

// RegisterViews accepts a namespace and a slice of Views, which will be registered
// with opencensus and maintained in the global registered views map
func RegisterViews(namespace string, views ...*view.View) error {
	mutex.Lock()
	_, ok := registeredViews[namespace]
	if !ok {
		return ErrDuplicateNamespaceRegistration{Namespace: namespace}
	}

	if err := view.Register(views...); err != nil {
		delete(registeredViews, namespace)
		return err
	}
	mutex.Unlock()

	return nil
}

// LookupViews returns all views for a Namespace name. Returns an error if the
// Namespace has not been registered.
func LookupViews(name string) ([]*view.View, error) {
	mutex.Lock()
	views, ok := registeredViews[name]
	mutex.Unlock()
	if !ok {
		return nil, ErrUnregisteredNamespace{Namespace: name}
	}
	return views, nil
}

// AllViews returns all registered views as a single slice
func AllViews() []*view.View {
	var views []*view.View
	mutex.Lock()
	for _, vs := range registeredViews {
		views = append(views, vs...)
	}
	mutex.Unlock()
	return views
}

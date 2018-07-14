package app

import (
	"net/url"
	"strings"
)

// Compo is the interface that describes a component.
// Must be implemented on a non empty struct pointer.
type Compo interface {
	// Render must return HTML 5.
	// It supports standard Go html/template API.
	// The pipeline is based on the component struct.
	// See https://golang.org/pkg/text/template and
	// https://golang.org/pkg/html/template for template usage.
	Render() string
}

// Mounter is the interface that wraps OnMount method.
type Mounter interface {
	Compo

	// OnMount is called when a component is mounted.
	// App.Render should not be called inside.
	OnMount()
}

// Dismounter is the interface that wraps OnDismount method.
type Dismounter interface {
	Compo

	// OnDismount is called when a component is dismounted.
	// App.Render should not be called inside.
	OnDismount()
}

// Navigable is the interface that wraps OnNavigate method.
type Navigable interface {
	Compo

	// OnNavigate is called when a component is loaded or navigated to.
	// It is called just after the component is mounted.
	OnNavigate(u *url.URL)
}

// Subscriber is the interface that describes a component that subscribes to
// events generated from actions.
type Subscriber interface {
	// Subscribe is called when a component is mounted.
	// The returned event subscriber is used to subscribe to events generated
	// from actions.
	// All the event subscribed are automatically unsuscribed when the component
	// is dismounted.
	Subscribe() EventSubscriber
}

// CompoWithExtendedRender is the interface that wraps Funcs method.
type CompoWithExtendedRender interface {
	Compo

	// Funcs returns a map of funcs to use when rendering a component.
	// Funcs named raw, json and time are reserved.
	// They handle raw html code, json conversions and time format.
	// They can't be overloaded.
	// See https://golang.org/pkg/text/template/#Template.Funcs for more details.
	Funcs() map[string]interface{}
}

// ZeroCompo is the type to use as base for an empty compone.
// Every instances of an empty struct is given the same memory address, which
// causes problem for indexing components.
// ZeroCompo have a placeholder field to avoid that.
type ZeroCompo struct {
	placeholder byte
}

// CompoNameFromURL is a helper function that returns the component name
// targeted by the given URL.
func CompoNameFromURL(u *url.URL) string {
	if len(u.Scheme) != 0 && u.Scheme != "compo" {
		return ""
	}

	path := u.Path
	path = strings.TrimPrefix(path, "/")

	paths := strings.SplitN(path, "/", 2)
	if len(paths[0]) == 0 {
		return ""
	}
	return normalizeCompoName(paths[0])
}

// CompoNameFromURLString is a helper function that returns the component
// name targeted by the given URL.
func CompoNameFromURLString(rawurl string) string {
	u, _ := url.Parse(rawurl)
	return CompoNameFromURL(u)
}

func normalizeCompoName(name string) string {
	name = strings.ToLower(name)
	if pkgsep := strings.IndexByte(name, '.'); pkgsep != -1 {
		if name[:pkgsep] == "main" {
			name = name[pkgsep+1:]
		}
	}
	return name
}

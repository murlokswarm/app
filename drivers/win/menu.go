// +build windows

package win

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/murlokswarm/app"
	"github.com/murlokswarm/app/internal/core"
	"github.com/murlokswarm/app/internal/dom"
	"github.com/pkg/errors"
)

// Menu implements the app.Menu interface.
type Menu struct {
	core.Menu

	id             string
	dom            dom.Engine
	typ            string
	compo          app.Compo
	keepWhenClosed bool
}

func newMenu(c app.MenuConfig, typ string) *Menu {
	m := &Menu{
		id: uuid.New().String(),
		dom: dom.Engine{
			Factory:   driver.Factory,
			Resources: driver.Resources,
			AllowedNodes: []string{
				"menu",
				"menuitem",
			},
			UI: driver.UI,
		},
		typ: typ,
	}

	m.dom.Sync = m.render

	if err := driver.Platform.Call("menus.New", nil, struct {
		ID string
	}{
		ID: m.id,
	}); err != nil {
		m.SetErr(err)
		return m
	}

	driver.Elems.Put(m)

	if len(c.URL) != 0 {
		m.Load(c.URL)
	}

	return m
}

// ID satisfies the app.Menu interface.
func (m *Menu) ID() string {
	return m.id
}

// Load satisfies the app.Menu interface.
func (m *Menu) Load(urlFmt string, v ...interface{}) {
	var err error
	defer func() {
		m.SetErr(err)
	}()

	u := fmt.Sprintf(urlFmt, v...)
	n := core.CompoNameFromURLString(u)

	var c app.Compo
	if c, err = driver.Factory.NewCompo(n); err != nil {
		return
	}

	m.compo = c

	if err = driver.Platform.Call("menus.Load", nil, struct {
		ID string
	}{
		ID: m.id,
	}); err != nil {
		return
	}

	err = m.dom.New(c)
	if err != nil {
		return
	}

	if nav, ok := c.(app.Navigable); ok {
		navURL, _ := url.Parse(u)
		nav.OnNavigate(navURL)
	}
}

// Compo satisfies the app.Menu interface.
func (m *Menu) Compo() app.Compo {
	return m.compo
}

// Contains satisfies the app.Menu interface.
func (m *Menu) Contains(c app.Compo) bool {
	return m.dom.Contains(c)
}

// Render satisfies the app.Menu interface.
func (m *Menu) Render(c app.Compo) {
	m.SetErr(m.dom.Render(c))
}

func (m *Menu) render(changes interface{}) error {
	b, err := json.Marshal(changes)
	if err != nil {
		return errors.Wrap(err, "encode changes failed")
	}

	return driver.Platform.Call("menus.Render", nil, struct {
		ID      string
		Changes string
	}{
		ID:      m.id,
		Changes: string(b),
	})
}

// Type satisfies the app.Menu interface.
func (m *Menu) Type() string {
	return m.typ
}

func onMenuCallback(m *Menu, in map[string]interface{}) {
	mappingStr := in["Mapping"].(string)

	var mapping dom.Mapping
	if err := json.Unmarshal([]byte(mappingStr), &mapping); err != nil {
		app.Logf("menu callback failed: %s", err)
		return
	}

	c, err := m.dom.CompoByID(mapping.CompoID)
	if err != nil {
		app.Logf("menu callback failed: %s", err)
		return
	}

	var f func()
	if f, err = mapping.Map(c); err != nil {
		app.Logf("menu callback failed: %s", err)
		return
	}

	if f != nil {
		f()
		return
	}

	app.Render(c)
}

func onMenuClose(m *Menu, in map[string]interface{}) {
	driver.Elems.Delete(m)
}

func handleMenu(h func(m *Menu, in map[string]interface{})) core.GoHandler {
	return func(in map[string]interface{}) {
		id, _ := in["ID"].(string)
		e := driver.Elems.GetByID(id)

		switch m := e.(type) {
		case *Menu:
			h(m, in)

		default:
			app.Panic("menu not supported")
		}
	}
}
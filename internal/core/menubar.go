package core

import (
	"encoding/json"
	"net/url"

	"github.com/murlokswarm/app"
)

func init() {
	app.Import(&MenuBar{})
}

// MenuBar is a component that represents a menu bar.
type MenuBar struct {
	AppName    string
	AppURL     string
	CustomURLs []string
	EditURL    string
	FileURL    string
	HelpURL    string
	WindowURL  string
}

// OnNavigate satisfies the app.Navigable interface.
func (m *MenuBar) OnNavigate(u *url.URL) {
	m.AppName = app.Name()
	app.Render(m)
}

// Render satisfies the app.Compo interface.
func (m *MenuBar) Render() string {
	return `
<menu label="{{.AppName}}">
	{{if .AppURL}}
	{{compo .AppURL}}
	{{else}}
	<menu>
		<menuitem label="About {{.AppName}}" role="about">
		<menuitem separator>
		<menuitem label="Preferences…" keys="cmdorctrl+," onclick="OnPreferences">
		<menuitem separator>
		<menuitem label="Hide {{.AppName}}" keys="cmdorctrl+h" role="hide">
		<menuitem label="Hide Others" keys="cmdorctrl+alt+h" role="hideOthers">
		<menuitem label="Show All" role="unhide">
		<menuitem separator>
		<menuitem label="Quit {{.AppName}}" keys="cmdorctrl+q" role="quit">
	</menu>
	{{end}}

	{{if .FileURL}}
	{{compo .FileURL}}
	{{end}}

	{{if .EditURL}}
	{{compo .EditURL}}
	{{else}}
	<menu label="Edit">
		<menuitem label="Undo" keys="cmdorctrl+z" role="undo">
		<menuitem label="Redo" keys="cmdorctrl+shift+z" role="redo">
		<menuitem separator>
		<menuitem label="Cut" keys="cmdorctrl+x" role="cut">
		<menuitem label="Copy" keys="cmdorctrl+c" role="copy">
		<menuitem label="Paste" keys="cmdorctrl+v" role="paste">
		<menuitem label="Delete" role="delete">
		<menuitem label="Select All" keys="cmdorctrl+a" role="selectAll">
	</menu>
	{{end}}

	{{range .CustomURLs}}
	{{compo .}}
	{{end}}

	{{if .WindowURL}}
	{{compo .WindowURL}}
	{{else}}
	<menu label="Window">
		<menuitem label="Minimize" keys="cmdorctrl+m" role="minimize">
		<menuitem label="Zoom" role="zoom">
		<menuitem separator>
		<menuitem label="Bring All to Front" role="arrangeInFront">
		<menuitem label="Close" keys="cmdorctrl+w" role="close">
	</menu>
	{{end}}

	{{if .WindowURL}}
	{{compo .WindowURL}}
	{{else}}
	<menu label="Help">
		<menuitem label="Built with github.com/murlokswarm/app" onclick="OnBuiltWith">
	</menu>
	{{end}}
</menu>
		`
}

// OnPreferences is the function called when the default app menu preferences
// button is clicked.
func (m *MenuBar) OnPreferences() {
	app.Emit(app.PreferencesRequested)
}

// OnBuiltWith is called when the On Built with button is clicked. It opens the
// app repository on the default browser.
func (m *MenuBar) OnBuiltWith() {
	app.OpenDefaultBrowser("https://github.com/murlokswarm/app")
}

func menuBarConfigToAddr(c app.MenuBarConfig) string {
	u, _ := url.Parse(app.CompoName(&MenuBar{}))
	u.Query().Set("AppURL", c.AppURL)
	u.Query().Set("EditURL", c.EditURL)
	u.Query().Set("FileURL", c.FileURL)
	u.Query().Set("HelpURL", c.HelpURL)
	u.Query().Set("WindowURL", c.WindowURL)

	customURLs, _ := json.Marshal(c.CustomURLs)
	u.Query().Set("CustomURLs", string(customURLs))

	return u.String()
}

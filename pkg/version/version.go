package version

import (
	"bytes"
	"runtime"
	"strings"
	"text/template"
)

// Build information
var (
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	AppName   = "sketch"
	GoVersion = runtime.Version()
)

var versionInfoTmpl = `
{{.program}}, version {{.version}} (branch: {{.branch}}, revision: {{.revision}})
  build user:       {{.buildUser}}
  build date:       {{.buildDate}}
  go version:       {{.goVersion}}
  platform:         {{.platform}}
`

// Info returns version and branch information
func Info() map[string]string {
	return map[string]string{
		"version":   Version,
		"branch":    Branch,
		"buildUser": BuildUser,
		"goVersion": GoVersion,
	}
}

// Print returns version information.
func Print() string {
	m := map[string]string{
		"program":   AppName,
		"version":   Version,
		"revision":  Revision,
		"branch":    Branch,
		"buildUser": BuildUser,
		"buildDate": BuildDate,
		"goVersion": GoVersion,
		"platform":  runtime.GOOS + "/" + runtime.GOARCH,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}

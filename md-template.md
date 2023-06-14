<!-- TOC -->
{{ range $ -}}
- [{{ .Name -}}](#{{ .Name -}})
{{ end }}

{{ range $ -}}
{{ .Anchor }}
## {{ .Name }}

{{ .Doc -}}
{{ if .Items }}

{{ .TableFieldName }} | {{ .TableFieldDoc }} | {{ .TableFieldRawType }}
{{ .TableFieldNameDashSize }} | {{ .TableFieldDocDashSize }} | {{ .TableFieldRawTypeDashSize }}
{{ end }}
{{- range .Items -}}
`{{ .Name }}` | {{ .Doc }}{{ if .Mandatory }} - *mandatory* {{ end }} {{ if .Inline }} - *inline* {{ end }} | {{ .RawType }}
{{ end }}
{{ end -}}

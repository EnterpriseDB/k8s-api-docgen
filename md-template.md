{{range $}}
#### {{.NameWithAnchor }}
{{.Doc -}}
{{if .Items }}

{{ .TableFieldName}} | {{ .TableFieldDoc}} | {{ .TableFieldRawType}} | Mandatory
{{ .TableFieldNameDashSize}} | {{ .TableFieldDocDashSize}} | {{ .TableFieldRawTypeDashSize}} | ---------
{{end}}
{{- range .Items -}}
{{.Name }} | {{.Doc }} | {{.RawType }} | {{.Mandatory }}
{{end}}
{{end}}
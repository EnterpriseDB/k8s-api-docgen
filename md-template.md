<!-- TOC -->
{{ range $ -}}
- [{{.Name -}}](#{{.Name -}})
{{end}}

{{ range $ -}}
## {{.NameWithAnchor }}

{{.Doc -}}
{{if .Items }}

{{ .TableFieldName}} | {{ .TableFieldDoc}} | {{ .TableFieldRawType}} |  {{ .TableFieldMandatory}}
{{ .TableFieldNameDashSize}} | {{ .TableFieldDocDashSize}} | {{ .TableFieldRawTypeDashSize}} | ---------
{{end}}
{{- range .Items -}}
{{.Name }} | {{.Doc }} | {{.RawType }} | {{.Mandatory }}
{{end}}
{{end -}}

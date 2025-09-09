export interface {{.TableName}} {
{{range .Fields}}  /** {{.Comment}} */
  {{.Name}}: {{.TsType}};
{{end}}
}

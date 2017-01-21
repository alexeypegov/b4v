{{- define "note"}}
<div class="note">
  <div class="title">{{.Title}}</div>
  <div class="date">{{timestamp .CreatedAt}}</div>
  {{- if .Tags}}
  <div class="tags">
    {{- range .Tags -}}
    <div class="tag">{{.}}</div>
    {{- end -}}
  </div>
  {{- end}}
  <div class="body">{{html .Content}}</div>
</div>
{{end -}}
<!doctype html>
<html>
<head>
  <title>{{.Vars.Title}}</title>
  <link rel="stylesheet" type="text/css" href="/css/blog.css"/>
</head>
<body>
{{- if .Note -}}
{{- template "note" .Note -}}
{{- else if .Notes -}}
{{- range .Notes -}}
{{- template "note" . -}}
{{- end -}}
{{- end -}}
</body>
</html>

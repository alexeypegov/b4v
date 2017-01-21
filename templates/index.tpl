{{- define "note"}}
<div class="note">
  <div class="title"><a href="/note/{{.UUID}}">{{.Title}}</a></div>
  <div class="date">{{timestamp .CreatedAt}}</div>
  <div class="body">{{html .Content}}</div>
  {{- if .Tags}}
  <div class="tags">
    {{- range .Tags -}}
    <div class="tag">{{.}}</div>
    {{- end -}}
  </div>
  {{- end}}
</div>
{{end -}}
{{- define "notes"}}
{{- range . -}}
{{- template "note" . -}}
{{- end -}}
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
{{- template "notes" .Notes -}}
{{- end -}}
<div class="footer">&copy;&nbsp;{{.Vars.Copyright}}</div>
</body>
</html>

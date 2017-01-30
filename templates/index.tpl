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
{{- define "paging-start" -}}
{{- if eq .Paging.Current 2 -}}
<div class="paging"><a href="/">{{- .Vars.previous_page -}}</a></div>
{{ else -}}
{{- if gt .Paging.Current 2 -}}
<div class="paging"><a href="/page/{{ minus .Paging.Current 1 }}">{{- .Vars.previous_page -}}</a></div>
{{ end -}}
{{- end -}}
{{- end -}}
{{- define "paging-end" -}}
{{- if lt .Paging.Current .Paging.Total -}}
<div class="paging"><a href="/page/{{ plus .Paging.Current 1 }}">{{- .Vars.next_page -}}</a></div>
{{ end -}}
{{- end -}}
<!doctype html>
<html>
<head>
  <title>{{.Vars.title}}</title>
  <link rel="stylesheet" type="text/css" href="/css/blog.css"/>
</head>
<body>
{{- if .Note -}}
{{- template "note" .Note -}}
{{- else if .Notes -}}
{{- template "paging-start" . -}}
{{- template "notes" .Notes -}}
{{- template "paging-end" . -}}
{{- end -}}
<div class="footer">&copy;&nbsp;{{.Vars.copyright}}</div>
</body>
</html>

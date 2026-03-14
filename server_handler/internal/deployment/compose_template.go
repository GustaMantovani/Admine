package deployment

import _ "embed"

//go:embed docker-compose.yaml.tmpl
var composeTemplate string

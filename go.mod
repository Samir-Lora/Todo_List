module todo_list

go 1.16

require (
	github.com/gobuffalo/buffalo v0.17.2
	github.com/gobuffalo/buffalo-pop/v2 v2.3.0
	github.com/gobuffalo/envy v1.9.0
	github.com/gobuffalo/mw-csrf v1.0.0
	github.com/gobuffalo/mw-paramlogger v1.0.0
	github.com/gobuffalo/pop v4.13.1+incompatible
	github.com/gobuffalo/pop/v5 v5.3.4
	github.com/gobuffalo/suite/v3 v3.0.2
	github.com/gobuffalo/validate v2.0.4+incompatible
	github.com/gobuffalo/validate/v3 v3.1.0
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/paganotoni/fsbox v1.1.3
	github.com/pkg/errors v0.9.1
	github.com/wawandco/oxpecker v1.5.4
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
)

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b

package sqlite

import _ "embed"

// embed the sqlite schema within the binary to create the tables at runtime.
//
//go:embed schema.sql
var Schema string

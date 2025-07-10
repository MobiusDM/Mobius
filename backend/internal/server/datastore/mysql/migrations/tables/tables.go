// Package tables provides database table migrations for Mobius
package tables

import (
	"github.com/notawar/mobius/internal/server/goose"
)

// MigrationClient for table migrations
var MigrationClient = goose.New("goose_db_version", &goose.MySqlDialect{})

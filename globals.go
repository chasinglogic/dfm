package cli

import (
	"github.com/chasinglogic/dfm/backend"
)

// DRYRUN is used to globally set whether this is a dry run
var DRYRUN = false

// Backend is the selected backend.
var Backend backend.Backend

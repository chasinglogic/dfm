package cli

import "github.com/chasinglogic/dfm/backend"

// DRYRUN indicates whether this is a dry run or not.
var DRYRUN bool

// Verbose controls the verbosity of information that dfm prints
var Verbose bool

// Backend is the currently selected DFM backend.
var Backend backend.Backend

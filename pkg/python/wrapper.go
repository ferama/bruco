package python

import "embed"

//go:embed wrapper/*
var pythonWrapper embed.FS

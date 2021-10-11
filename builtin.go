package main

import "embed"

//go:embed builtin/**
var builtin embed.FS

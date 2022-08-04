package main

import (
	"io/fs"
	"os"
	"time"
)

type PreviousDir struct{}

func (p *PreviousDir) IsDir() bool        { return true }
func (p *PreviousDir) Name() string       { return ".." }
func (p *PreviousDir) Size() int64        { return 0 }
func (p *PreviousDir) Mode() fs.FileMode  { return os.FileMode(0) }
func (p *PreviousDir) ModTime() time.Time { return time.Time{} }
func (p *PreviousDir) Sys() any           { return nil }

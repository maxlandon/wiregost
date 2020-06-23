package generic

import corepb "github.com/maxlandon/wiregost/proto/v1/gen/go/ghost/core"

// FileSystem - The base ghost type implements this interface
type FileSystem interface {
	Ls(path string) (ls *corepb.Ls)
	Cd(path string)
	Pwd() (pwd *corepb.Pwd)
	Rm(path string) (rm *corepb.Rm)
	Cat(file string) (cat *corepb.Download)
	Mkdir(name string) (dir *corepb.Mkdir)
	Download(file string) (dl *corepb.Download)
	Upload(file string) (upl *corepb.Upload)
}

// Ls - Returns the contents of a directory
func (g *Ghost) Ls(path string) (ls *corepb.Ls) {
	return
}

// Cd - Change the working directory of the implant
func (g *Ghost) Cd(path string) {
	return
}

// Pwd - Print the implant working directory
func (g *Ghost) Pwd() (pwd *corepb.Pwd) {
	return
}

// Rm - Remove a directory/file on target
func (g *Ghost) Rm(path string) (rm *corepb.Rm) {
	return
}

// Cat - Download and print a file on target
func (g *Ghost) Cat(file string) (cat *corepb.Download) {
	return
}

// Mkdir - Make a new directory on target
func (g *Ghost) Mkdir(name string) (dir *corepb.Mkdir) {
	return
}

// Download - Download a file from the target
func (g *Ghost) Download(file string) (dl *corepb.Download) {
	return
}

// Upload - Upload a local file onto the target
func (g *Ghost) Upload(file string) (upl *corepb.Upload) {
	return
}

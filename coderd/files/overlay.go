package files

import (
	"io/fs"
	"path"
	"strings"
)

// overlayFS allows you to "join" together the template files tar file fs.FS
// with the Terraform modules tar file fs.FS. We could potentially turn this
// into something more parameterized/configurable, but the requirements here are
// a _bit_ odd, because every file in the modulesFS includes the
// .terraform/modules/ folder at the beginning of it's path.
type overlayFS struct {
	baseFS   fs.FS
	overlays []Overlay
}

type Overlay struct {
	Path string
	fs.FS
}

func NewOverlayFS(baseFS fs.FS, overlays []Overlay) fs.FS {
	return overlayFS{
		baseFS:   baseFS,
		overlays: overlays,
	}
}

func (f overlayFS) Open(p string) (fs.File, error) {
	for _, overlay := range f.overlays {
		if strings.HasPrefix(path.Clean(p), overlay.Path) {
			return overlay.FS.Open(p)
		}
	}
	return f.baseFS.Open(p)
}

func (f overlayFS) ReadDir(p string) ([]fs.DirEntry, error) {
	for _, overlay := range f.overlays {
		if strings.HasPrefix(path.Clean(p), overlay.Path) {
			return overlay.FS.(fs.ReadDirFS).ReadDir(p)
		}
	}
	return f.baseFS.(fs.ReadDirFS).ReadDir(p)
}

func (f overlayFS) ReadFile(p string) ([]byte, error) {
	for _, overlay := range f.overlays {
		if strings.HasPrefix(path.Clean(p), overlay.Path) {
			return overlay.FS.(fs.ReadFileFS).ReadFile(p)
		}
	}
	return f.baseFS.(fs.ReadFileFS).ReadFile(p)
}

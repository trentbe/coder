package files_test

import (
	"io/fs"
	"testing"

	"github.com/coder/coder/v2/coderd/files"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestWoohoo(t *testing.T) {
	a := afero.NewMemMapFs()
	afero.WriteFile(a, "main.tf", []byte("terraform {}"), 0o644)
	b := afero.NewMemMapFs()
	afero.WriteFile(b, ".terraform/modules/modules.json", []byte("{}"), 0o644)
	it := files.NewOverlayFS(afero.NewIOFS(a), afero.NewIOFS(b), ".terraform/modules")

	content, err := fs.ReadFile(it, "main.tf")
	require.NoError(t, err)
	require.Equal(t, "terraform {}", string(content))
	content, err = fs.ReadFile(it, ".terraform/modules/modules.json")
	require.NoError(t, err)
	require.Equal(t, "{}", string(content))
}

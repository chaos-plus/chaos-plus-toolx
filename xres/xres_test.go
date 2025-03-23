package xres

import (
	"embed"
	"fmt"
	"testing"

	"github.com/chaos-plus/chaos-plus-toolx/xfile"
	"github.com/stretchr/testify/require"
)

//go:embed all:data_test
var testFs embed.FS

func TestRes(t *testing.T) {
	res := New(testFs)

	items, err := res.ScanAll()

	dirs := []string{}
	files := []string{}
	for _, item := range items {
		fmt.Println(item)
		if item.IsDir {
			dirs = append(dirs, item.Path)
		} else {
			files = append(files, item.Path)
		}
	}

	require.NoError(t, err)
	require.Equal(t, 3, len(dirs))
	require.Equal(t, 3, len(files))

	items, err = res.ScanDirFile("data_test/public", "*", true)
	require.NoError(t, err)
	require.Equal(t, 2, len(items))

	sql, err := res.GetContent("data_test/sql/kv.sql")
	require.NoError(t, err)
	require.Contains(t, string(sql), "CREATE TABLE kv")

	xfile.RemoveAll("export_test")

	err = res.Export("data_test/public", "export_test", true)
	require.NoError(t, err)

	err = res.Export("data_test/public", "export_test/public", true)
	require.NoError(t, err)

	err = res.Export("data_test", "export_test", true)
	require.NoError(t, err)

	err = res.Export("data_test", "export_test", false)
	require.Error(t, err)
}

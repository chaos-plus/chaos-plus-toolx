package xres

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/chaos-plus/chaos-plus-toolx/xfile"
)

type ResFile struct {
	Path  string
	Name  string
	Info  fs.FileInfo
	IsDir bool
}

type ResFiles struct {
	Root embed.FS
}

func New(root embed.FS) *ResFiles {
	return &ResFiles{Root: root}
}

func (r *ResFiles) DumpAll() error {
	items, err := r.ScanAll()
	if err != nil {
		return err
	}
	for _, entry := range items {
		fmt.Println(entry)
	}
	return nil
}

func (r *ResFiles) ScanAll() ([]ResFile, error) {
	return r.Scan(".", true)
}

func (r *ResFiles) Scan(path string, recursive bool) ([]ResFile, error) {
	fs, err := r.Root.Open(path)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	stat, err := fs.Stat()
	if err != nil {
		return nil, err
	}
	list := []ResFile{}
	if path != "." {
		list = append(list, ResFile{
			Name:  stat.Name(),
			Path:  path,
			Info:  stat,
			IsDir: stat.IsDir(),
		})
	}
	if !stat.IsDir() {
		return list, nil
	}
	if !recursive && len(list) > 1 {
		return list, nil
	}
	items, err := r.Root.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range items {
		epath := ""
		if path == "." {
			epath = entry.Name()
		} else {
			epath = path + "/" + entry.Name()
		}
		efiles, err := r.Scan(epath, entry.IsDir() && recursive)
		if err != nil {
			return nil, err
		}
		list = append(list, efiles...)
	}

	return list, nil
}

func (r *ResFiles) ScanDirFile(path string, pattern string, recursive bool) ([]ResFile, error) {
	items, err := r.Scan(path, recursive)
	if err != nil {
		return nil, err
	}
	list := []ResFile{}
	for _, item := range items {
		if item.IsDir {
			continue
		}
		match, err := filepath.Match(pattern, item.Name)
		if err != nil {
			return nil, err
		}
		if match {
			list = append(list, item)
		}
	}
	return list, nil
}

func (r *ResFiles) GetDirs() ([]string, error) {
	items, err := r.ScanAll()
	if err != nil {
		return nil, err
	}
	dirs := []string{}
	for _, item := range items {
		if item.IsDir {
			dirs = append(dirs, item.Path)
		}
	}
	return dirs, nil
}

func (r *ResFiles) GetFiles() ([]string, error) {
	items, err := r.ScanAll()
	if err != nil {
		return nil, err
	}
	files := []string{}
	for _, item := range items {
		if !item.IsDir {
			files = append(files, item.Path)
		}
	}
	return files, nil
}

func (r *ResFiles) GetFileInfo(path string) (fs.FileInfo, error) {
	fs, err := r.Root.Open(path)
	if err != nil {
		return nil, err
	}
	return fs.Stat()
}

func (r *ResFiles) GetContent(path string) ([]byte, error) {
	file, err := r.Root.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func (r *ResFiles) Export(from, to string, overwrite bool) error {
	files, err := r.Scan(from, true)
	if err != nil {
		return err
	}
	prefix := xfile.GetParentDirectory(from)
	for _, file := range files {
		toFile := filepath.Join(to, strings.TrimPrefix(file.Path, prefix))
		if file.IsDir {
			xfile.MkdirAll(toFile)
			continue
		}
		reader, err := r.Root.Open(file.Path)
		if err != nil {
			return err
		}
		err = xfile.CopyToFile(reader, toFile, overwrite)
		if err != nil {
			return err
		}
	}
	return nil
}

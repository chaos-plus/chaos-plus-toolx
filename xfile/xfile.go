package xfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func GetTempPath(path ...string) string {
	return filepath.Join(os.TempDir(), filepath.Join(path...))
}

func GetUserPath(path ...string) (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filepath.Join(path...)), nil
}

func GetUserCachePath(path ...string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filepath.Join(path...)), nil
}

func GetUserConfigPath(path ...string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filepath.Join(path...)), nil
}

func GetExecutableTempPath(path ...string) string {
	path = append([]string{GetExecutableName()}, path...)
	return GetTempPath(path...)
}

func GetExecutableUserPath(path ...string) (string, error) {
	path = append([]string{GetExecutableName()}, path...)
	return GetUserPath(path...)
}

func GetExecutableUserCachePath(path ...string) (string, error) {
	path = append([]string{GetExecutableName()}, path...)
	return GetUserCachePath(path...)
}

func GetExecutableUserConfig(path ...string) (string, error) {
	path = append([]string{GetExecutableName()}, path...)
	return GetUserConfigPath(path...)
}

func GetFileName(path string) string {
	return filepath.Base(path)
}

func GetFileNameWithoutExt(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	return name[:len(name)-len(ext)]
}

func GetParentDirectory(path string) string {
	return filepath.Dir(path)
}

func NormalizePath(path string) string {
	return filepath.ToSlash(path)
}

func GetFileSize(path string) int64 {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	return fileInfo.Size()
}

func GetExecutablePath() string {
	path, _ := os.Executable()
	return path
}

func GetExecutableName() string {
	return GetFileName(GetExecutablePath())
}

func GetExecutableNameWithoutExt() string {
	return GetFileNameWithoutExt(GetExecutablePath())
}

func GetExecutableDirectory() string {
	return GetParentDirectory(GetExecutablePath())
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func IsFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func MkdirParent(path ...string) error {
	if len(path) == 0 {
		return nil
	}
	for _, p := range path {
		parent := GetParentDirectory(p)
		if err := os.MkdirAll(parent, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func MkdirAll(path ...string) error {
	if len(path) == 0 {
		return nil
	}
	for _, p := range path {
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func RemoveAll(path ...string) error {
	if len(path) == 0 {
		return nil
	}
	for _, p := range path {
		if err := os.RemoveAll(p); err != nil {
			return err
		}
	}
	return nil
}

func RemoveFile(path ...string) error {
	if len(path) == 0 {
		return nil
	}
	for _, p := range path {
		if err := os.Remove(p); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
	}
	return nil
}

// CopyDir recursively copies a directory from srcDir to dstDir
func CopyDir(srcDir, dstDir string, overwrite bool) error {
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("failed to get source directory info: %v", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source path %s is not a directory", srcDir)
	}

	err = os.MkdirAll(dstDir, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return CopyFile(path, destPath, overwrite)
	})
	if err != nil {
		return fmt.Errorf("error copying directory: %v", err)
	}
	return nil
}

func CopyFile(srcFile, dstFile string, overwrite bool) error {
	if IsDirectory(srcFile) {
		return fmt.Errorf("source path %s is a directory", srcFile)
	}
	if IsDirectory(dstFile) {
		return fmt.Errorf("destination path %s is a directory", dstFile)
	}
	src, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %v", srcFile, err)
	}
	defer src.Close()
	err = CopyToFile(src, dstFile, overwrite)
	if err != nil {
		return err
	}
	srcInfo, err := os.Stat(srcFile)
	if err != nil {
		return err
	}
	return os.Chmod(dstFile, srcInfo.Mode())
}

func CopyToFile(srcFile io.Reader, dstFile string, overwrite bool) error {
	if IsExist(dstFile) {
		if overwrite {
			err := RemoveFile(dstFile)
			if err != nil {
				return fmt.Errorf("failed to remove destination file %s: %v", dstFile, err)
			}
		} else {
			return fmt.Errorf("dest file is already exists: %s", dstFile)
		}
	}

	dst, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %v", dstFile, err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file from %s to %s: %v", srcFile, dstFile, err)
	}
	return nil
}

/*
Copyright © 2025 this guy Labs <thisguy@thisguylabs.com>

This file is part of GVT (Guy's Versioning Tool).

Do not remove or modify this notice.
*/

package storage

import (
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DecompressToFile decompresses a .zlib file to the destination path.
func DecompressToFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	r, err := zlib.NewReader(srcFile)
	if err != nil {
		return fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer r.Close()

	err = os.MkdirAll(getDir(dstPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, r)
	if err != nil {
		return fmt.Errorf("failed to write decompressed data: %w", err)
	}

	return nil
}

// getDir returns the directory portion of a path
func getDir(path string) string {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return filepath.Dir(path)
	}
	return path
}

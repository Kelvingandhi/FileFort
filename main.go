package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Check if a directory exists
func checkDirectoryExists(dir string) error {
	info, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return fmt.Errorf("%s is not a directory", dir)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}
	return nil
}

// check if program has read permissions for a file
func hasReadPermissions(info os.FileInfo) bool {
	return info.Mode().Perm()&0400 != 0
}

func copyFile(src, dst string) error {

	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", src, err)
	}

	// check read permissions
	if !hasReadPermissions(info) {
		return fmt.Errorf("source file %s does not have read permissions", src)
	}

	in_file, err := os.Open(src)

	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer in_file.Close()

	out_file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer out_file.Close()

	_, err = io.Copy(out_file, in_file)

	if err != nil {
		return fmt.Errorf("error while copying file %s to %s: %w", src, dst, err)
	}

	return nil
}

func backup_files(source_dir, backup_dir, filename string) error {

	// Ensure source directory exists
	if err := checkDirectoryExists(source_dir); err != nil {
		return err
	}

	// Ensure backup directory exists or create it
	if err := os.MkdirAll(backup_dir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory %s: %w", backup_dir, err)
	}

	// backup all the files under source_dir
	err := filepath.Walk(source_dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}
		if !info.IsDir() {
			dst := filepath.Join(backup_dir, info.Name())
			err := copyFile(path, dst)
			if err != nil {
				fmt.Printf("Error backing up %s: %v\n", info.Name(), err)
				return err
			}
			fmt.Printf("Copied: %s to %s\n", path, dst)
		}
		return nil
	})
	return err

}

func main() {
	fmt.Println("Hello world!")

	source_dir := "/home/kelvin/Documents/VS Code Projects/go language/Project Files/source"
	backup_dir := "/home/kelvin/Documents/VS Code Projects/go language/Project Files/backup"
	filename := "hello.txt"

	err := backup_files(source_dir, backup_dir, filename)

	if err != nil {
		fmt.Println("Error backing up files: ", err)

	} else {
		fmt.Println("File(s) backed up successfully!")
	}

}

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func backup_files(sourceDir, backupDir, fileFilter, typeFilter string) error {

	// Ensure source directory exists
	if err := checkDirectoryExists(sourceDir); err != nil {
		return err
	}

	// Ensure backup directory exists or create it
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory %s: %w", backupDir, err)
	}

	// Track the number of files successfully backed up
	filesBackedUp := 0

	// backup all the files under source_dir
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Apply file filters
		if fileFilter != "" && info.Name() != fileFilter {
			return nil
		}

		if typeFilter != "" && !strings.HasSuffix(info.Name(), typeFilter) {
			return nil
		}

		// Construct destination path
		dst := filepath.Join(backupDir, info.Name())

		// Copy file
		err = copyFile(path, dst)
		if err != nil {
			fmt.Printf("Error backing up %s: %v\n", info.Name(), err)
			return err
		}
		fmt.Printf("Copied: %s to %s\n", path, dst)
		filesBackedUp++
		return nil
	})

	if err != nil {
		return fmt.Errorf("backup process encountered an error: %w", err)
	}

	if filesBackedUp == 0 && fileFilter != "" {
		fmt.Printf("%d files backed up.\n", filesBackedUp)
		return fmt.Errorf("no file matching the filter %s found in %s", fileFilter, sourceDir)
	}

	if filesBackedUp == 0 && typeFilter != "" {
		fmt.Printf("%d files backed up.\n", filesBackedUp)
		return fmt.Errorf("no file matching the filter %s found in %s", typeFilter, sourceDir)
	}

	fmt.Printf("%d file(s) backed up.\n", filesBackedUp)

	return nil

}

// schedule for backup at regular intervals
func scheduleBackup(source_dir, backup_dir, fileFilter, typeFilter string, interval time.Duration) {

	for {
		fmt.Printf("Starting backup at %s ...\n", time.Now().Format(time.RFC1123))
		err := backup_files(source_dir, backup_dir, fileFilter, typeFilter)

		if err != nil {
			fmt.Println("Error backing up files: ", err)

		} else {
			fmt.Println("File(s) backed up successfully!")
		}
		fmt.Printf("Next backup scheduled at %s ...\n", time.Now().Add(interval).Format(time.RFC1123))
		time.Sleep(interval)
	}
}

func main() {
	fmt.Println("Hello world!")

	sourceDir := flag.String("source", "", "Path to the Source directory")
	backupDir := flag.String("backup", "", "Path to the Backup directory")
	interval := flag.Int("interval", 0, "Interval for backup")
	fileFilter := flag.String("file", "", "Specific file to backup")
	typeFilter := flag.String("type", "", "Specific file type to backup (e.g. .txt, .pdf)")

	// source_dir := "/home/kelvin/Documents/VS Code Projects/go language/Project Files/source"
	// backup_dir := "/home/kelvin/Documents/VS Code Projects/go language/Project Files/backup"
	// filename := "hello.txt"

	flag.Parse()

	// validate flags
	if *sourceDir == "" || *backupDir == "" {
		fmt.Println("Source and Backup directories are required")
		flag.Usage()
		os.Exit(1)
	}

	if *fileFilter != "" && *typeFilter != "" {
		fmt.Println("Cannot specify both file and type filters together. Specify Either")
		flag.Usage()
		os.Exit(1)
	}

	// Convert interval to time.Duration
	backupInterval := time.Duration(*interval) * time.Second

	scheduleBackup(*sourceDir, *backupDir, *fileFilter, *typeFilter, backupInterval)

}

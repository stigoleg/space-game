package systems

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// UpdateInstaller handles platform-specific binary replacement
type UpdateInstaller struct {
	platform string
	arch     string
}

// NewUpdateInstaller creates a new update installer
func NewUpdateInstaller() *UpdateInstaller {
	return &UpdateInstaller{
		platform: runtime.GOOS,
		arch:     runtime.GOARCH,
	}
}

// InstallUpdate installs the downloaded update file
func (ui *UpdateInstaller) InstallUpdate(downloadedFilePath string) error {
	switch ui.platform {
	case "darwin":
		return ui.installUpdateDarwin(downloadedFilePath)
	case "windows":
		return ui.installUpdateWindows(downloadedFilePath)
	case "linux":
		return ui.installUpdateLinux(downloadedFilePath)
	default:
		return fmt.Errorf("unsupported platform: %s", ui.platform)
	}
}

// installUpdateDarwin handles macOS DMG installation
func (ui *UpdateInstaller) installUpdateDarwin(dmgPath string) error {
	log.Println("Installing macOS update...")

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Resolve symlinks
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	// Check if running from .app bundle
	isAppBundle := strings.Contains(execPath, ".app/Contents/MacOS")

	if isAppBundle {
		return ui.installDMGToAppBundle(dmgPath, execPath)
	}

	// For standalone binary, we can't easily extract from DMG
	// So we'll just notify the user
	return fmt.Errorf("standalone macOS binary update not yet implemented - please download manually")
}

// installDMGToAppBundle installs DMG to .app bundle location
func (ui *UpdateInstaller) installDMGToAppBundle(dmgPath, currentExecPath string) error {
	// Mount the DMG
	log.Println("Mounting DMG...")
	mountPoint, err := ui.mountDMG(dmgPath)
	if err != nil {
		return fmt.Errorf("failed to mount DMG: %w", err)
	}
	defer ui.unmountDMG(mountPoint)

	// Find .app bundle in mount point
	appBundleName := "Stellar Siege.app"
	srcAppBundle := filepath.Join(mountPoint, appBundleName)

	if _, err := os.Stat(srcAppBundle); os.IsNotExist(err) {
		return fmt.Errorf(".app bundle not found in DMG: %s", srcAppBundle)
	}

	// Get current .app bundle path (up from Contents/MacOS/stellar-siege)
	// currentExecPath is like: /Applications/Stellar Siege.app/Contents/MacOS/stellar-siege
	currentAppBundle := filepath.Dir(filepath.Dir(filepath.Dir(currentExecPath)))

	log.Printf("Current app bundle: %s", currentAppBundle)
	log.Printf("New app bundle: %s", srcAppBundle)

	// Backup current .app by renaming
	backupPath := currentAppBundle + ".old"
	log.Printf("Backing up to: %s", backupPath)

	// Remove old backup if exists
	os.RemoveAll(backupPath)

	// Rename current app to .old
	if err := os.Rename(currentAppBundle, backupPath); err != nil {
		return fmt.Errorf("failed to backup current app: %w", err)
	}

	// Copy new app bundle
	log.Println("Copying new app bundle...")
	if err := ui.copyDir(srcAppBundle, currentAppBundle); err != nil {
		// Restore backup on failure
		log.Println("Copy failed, restoring backup...")
		os.Rename(backupPath, currentAppBundle)
		return fmt.Errorf("failed to copy new app bundle: %w", err)
	}

	// Remove quarantine attribute
	log.Println("Removing quarantine attribute...")
	exec.Command("xattr", "-cr", currentAppBundle).Run()

	// Launch new app and exit
	log.Println("Launching new version...")
	cmd := exec.Command("open", currentAppBundle)
	if err := cmd.Start(); err != nil {
		log.Printf("Warning: Failed to launch new app: %v", err)
	}

	// Give it a moment to start
	time.Sleep(500 * time.Millisecond)

	// Exit current process
	log.Println("Update complete, exiting...")
	os.Exit(0)

	return nil
}

// mountDMG mounts a DMG file and returns the mount point
func (ui *UpdateInstaller) mountDMG(dmgPath string) (string, error) {
	cmd := exec.Command("hdiutil", "attach", dmgPath, "-nobrowse", "-noverify")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("hdiutil attach failed: %w\nOutput: %s", err, string(output))
	}

	// Parse output to find mount point
	// Output format: /dev/disk4s2       	Apple_HFS                      	/Volumes/Stellar Siege
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "/Volumes/") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				mountPoint := parts[len(parts)-1]
				return mountPoint, nil
			}
		}
	}

	return "", fmt.Errorf("failed to parse mount point from hdiutil output")
}

// unmountDMG unmounts a DMG
func (ui *UpdateInstaller) unmountDMG(mountPoint string) error {
	cmd := exec.Command("hdiutil", "detach", mountPoint, "-quiet")
	return cmd.Run()
}

// installUpdateWindows handles Windows ZIP installation
func (ui *UpdateInstaller) installUpdateWindows(zipPath string) error {
	log.Println("Installing Windows update...")

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Extract ZIP to temp directory
	tempDir := filepath.Join(os.TempDir(), "stellar-siege-update-extracted")
	os.RemoveAll(tempDir) // Clean old extractions
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	log.Println("Extracting ZIP...")
	if err := ui.extractZip(zipPath, tempDir); err != nil {
		return fmt.Errorf("failed to extract ZIP: %w", err)
	}

	// Find new executable in extracted files
	newExePath := filepath.Join(tempDir, "stellar-siege-windows", "stellar-siege.exe")
	if _, err := os.Stat(newExePath); os.IsNotExist(err) {
		return fmt.Errorf("stellar-siege.exe not found in ZIP")
	}

	// Create batch script for update
	batchScript := fmt.Sprintf(`@echo off
echo Updating Stellar Siege...
timeout /t 2 /nobreak >nul
move /Y "%s" "%s.old"
move /Y "%s" "%s"
echo Update complete, starting game...
start "" "%s"
del "%%~f0"
`, currentExe, currentExe, newExePath, currentExe, currentExe)

	scriptPath := filepath.Join(os.TempDir(), "stellar-siege-update.bat")
	if err := os.WriteFile(scriptPath, []byte(batchScript), 0755); err != nil {
		return fmt.Errorf("failed to create update script: %w", err)
	}

	// Launch batch script and exit
	log.Println("Launching update script...")
	cmd := exec.Command("cmd", "/C", scriptPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch update script: %w", err)
	}

	time.Sleep(500 * time.Millisecond)
	log.Println("Update initiated, exiting...")
	os.Exit(0)

	return nil
}

// installUpdateLinux handles Linux tar.gz installation
func (ui *UpdateInstaller) installUpdateLinux(tarPath string) error {
	log.Println("Installing Linux update...")

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Extract tar.gz to temp directory
	tempDir := filepath.Join(os.TempDir(), "stellar-siege-update-extracted")
	os.RemoveAll(tempDir)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	log.Println("Extracting tar.gz...")
	if err := ui.extractTarGz(tarPath, tempDir); err != nil {
		return fmt.Errorf("failed to extract tar.gz: %w", err)
	}

	// Find new binary
	newBinaryPath := filepath.Join(tempDir, "stellar-siege-linux", "stellar-siege")
	if _, err := os.Stat(newBinaryPath); os.IsNotExist(err) {
		return fmt.Errorf("stellar-siege binary not found in archive")
	}

	// Create shell script for update
	shellScript := fmt.Sprintf(`#!/bin/bash
echo "Updating Stellar Siege..."
sleep 1
mv "%s" "%s.old"
mv "%s" "%s"
chmod +x "%s"
echo "Update complete, starting game..."
"%s" &
rm -- "$0"
`, currentExe, currentExe, newBinaryPath, currentExe, currentExe, currentExe)

	scriptPath := filepath.Join(os.TempDir(), "stellar-siege-update.sh")
	if err := os.WriteFile(scriptPath, []byte(shellScript), 0755); err != nil {
		return fmt.Errorf("failed to create update script: %w", err)
	}

	// Launch shell script and exit
	log.Println("Launching update script...")
	cmd := exec.Command("sh", scriptPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch update script: %w", err)
	}

	time.Sleep(500 * time.Millisecond)
	log.Println("Update initiated, exiting...")
	os.Exit(0)

	return nil
}

// extractZip extracts a ZIP file to destination directory
func (ui *UpdateInstaller) extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// extractTarGz extracts a tar.gz file to destination directory
func (ui *UpdateInstaller) extractTarGz(tarPath, destDir string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, header.Name)

		// Check for ZipSlip vulnerability
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", target)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// copyDir recursively copies a directory
func (ui *UpdateInstaller) copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		return ui.copyFile(path, destPath, info.Mode())
	})
}

// copyFile copies a single file
func (ui *UpdateInstaller) copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

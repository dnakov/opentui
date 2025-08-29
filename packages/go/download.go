package opentui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	downloadOnce sync.Once
	downloadErr  error
)

func init() {
	downloadOnce.Do(func() {
		downloadErr = ensureAssets()
	})
	if downloadErr != nil {
		panic(fmt.Sprintf("Failed to download OpenTUI assets: %v", downloadErr))
	}
}

func ensureAssets() error {
	// Get the directory where this package is located
	_, filename, _, _ := runtime.Caller(0)
	packageDir := filepath.Dir(filename)
	
	headerPath := filepath.Join(packageDir, "opentui.h")
	if _, err := os.Stat(headerPath); err == nil {
		// Assets already exist
		return nil
	}

	// Download assets
	const repoURL = "https://github.com/sst/opentui/releases/download"
	
	// Create lib directory structure
	platforms := []string{
		"aarch64-linux", "aarch64-macos", "aarch64-windows",
		"x86_64-linux", "x86_64-macos", "x86_64-windows",
	}

	for _, platform := range platforms {
		libDir := filepath.Join(packageDir, "lib", platform)
		if err := os.MkdirAll(libDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", libDir, err)
		}
	}

	// Download header file
	headerURL := fmt.Sprintf("%s/latest/download/opentui.h", repoURL)
	if err := downloadFile(headerPath, headerURL); err != nil {
		return fmt.Errorf("failed to download header: %v", err)
	}

	// Download current platform's library
	platform := getCurrentPlatform()
	libURL := fmt.Sprintf("%s/latest/download/%s", repoURL, getLibraryAssetName(platform))
	libPath := filepath.Join(packageDir, "lib", platform, getLibraryFileName(platform))
	
	if err := downloadFile(libPath, libURL); err != nil {
		return fmt.Errorf("failed to download library for %s: %v", platform, err)
	}

	return nil
}

func downloadFile(filePath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d downloading %s", resp.StatusCode, url)
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getCurrentPlatform() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	
	if goarch == "arm64" {
		goarch = "aarch64"
	} else if goarch == "amd64" {
		goarch = "x86_64"
	}

	return fmt.Sprintf("%s-%s", goarch, goos)
}

func getLibraryAssetName(platform string) string {
	parts := strings.Split(platform, "-")
	arch, os := parts[0], parts[1]
	
	switch os {
	case "windows":
		return fmt.Sprintf("opentui-%s-windows.dll", arch)
	case "macos":
		return fmt.Sprintf("libopentui-%s-macos.dylib", arch)
	case "linux":
		return fmt.Sprintf("libopentui-%s-linux.so", arch)
	default:
		return ""
	}
}

func getLibraryFileName(platform string) string {
	parts := strings.Split(platform, "-")
	os := parts[1]
	
	switch os {
	case "windows":
		return "opentui.dll"
	case "macos":
		return "libopentui.dylib"
	case "linux":
		return "libopentui.so"
	default:
		return ""
	}
}
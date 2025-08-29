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
		downloadErr = ensureLibraries()
	})
	if downloadErr != nil {
		panic(fmt.Sprintf("Failed to download OpenTUI libraries: %v", downloadErr))
	}
}

func ensureLibraries() error {
	_, filename, _, _ := runtime.Caller(0)
	packageDir := filepath.Dir(filename)
	
	platform := getCurrentPlatform()
	libPath := filepath.Join(packageDir, "lib", platform, getLibraryFileName(platform))
	
	if _, err := os.Stat(libPath); err == nil {
		return nil
	}

	const repoURL = "https://github.com/dnakov/opentui/releases/download"
	
	if err := os.MkdirAll(filepath.Dir(libPath), 0755); err != nil {
		return fmt.Errorf("failed to create lib directory: %v", err)
	}

	libURL := fmt.Sprintf("%s/latest/download/%s", repoURL, getLibraryAssetName(platform))
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
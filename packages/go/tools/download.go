package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	releaseURL = "https://api.github.com/repos/sst/opentui/releases/latest"
	repoURL    = "https://github.com/sst/opentui/releases/download"
)

func main() {
	if err := downloadAssets(); err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading assets: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Assets downloaded successfully")
}

func downloadAssets() error {
	// Create lib directory structure
	platforms := []string{
		"aarch64-linux",
		"aarch64-macos", 
		"aarch64-windows",
		"x86_64-linux",
		"x86_64-macos",
		"x86_64-windows",
	}

	for _, platform := range platforms {
		if err := os.MkdirAll(filepath.Join("lib", platform), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// Download header file from latest release
	headerURL := fmt.Sprintf("%s/latest/download/opentui.h", repoURL)
	if err := downloadFile("opentui.h", headerURL); err != nil {
		return fmt.Errorf("failed to download header: %v", err)
	}

	// Download appropriate library for current platform
	platform := getPlatform()
	libURL := fmt.Sprintf("%s/latest/download/%s", repoURL, getLibraryAssetName(platform))
	libPath := filepath.Join("lib", platform, getLibraryFileName(platform))
	
	if err := downloadFile(libPath, libURL); err != nil {
		return fmt.Errorf("failed to download library for %s: %v", platform, err)
	}

	return nil
}

func downloadFile(filepath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download %s: status %d", url, resp.StatusCode)
	}

	if err := os.MkdirAll(filepath[:strings.LastIndex(filepath, "/")], 0755); err != nil {
		return err
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getPlatform() string {
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
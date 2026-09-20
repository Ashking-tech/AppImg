package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func runPipeline(path string, u *ui) string {

	//validate
	u.step("Validate", func() (string, error) {
		regular, err := isRegular(path)
		if err != nil {
			return "", err
		}
		if !regular {
			return "", fmt.Errorf("not a regular file: %s", path)
		}
		return path, nil
	})

	//create /Applications
	u.step("Create Applications dir", func() (string, error) {
		if err := CreateApplicationsDir(); err != nil {
			return "", err
		}
		return "", nil
	})

	//get destination
	home, err := os.UserHomeDir()
	if err != nil {
		fail("%v", err)
	}

	destination := filepath.Join(
		home,
		"Applications",
		filepath.Base(path),
	)

	//Moving the file
	u.step("Move AppImage", func() (string, error) {
		if err := MoveAppImage(path, destination); err != nil {
			return "", err
		}
		return destination, nil
	})

	//extracting app image
	u.step("Extract AppImage", func() (string, error) {
		if err := ExtractAppImage(destination); err != nil {
			return "", err
		}
		return "", nil
	})

	//finding desktopfile
	extractRoot := filepath.Join(filepath.Dir(destination), "squashfs-root")
	var desktop string
	u.step("Find desktop file", func() (string, error) {
		d, err := FindDesktopFile(extractRoot)
		if err != nil {
			return "", err
		}
		desktop = d
		return d, nil
	})

	//finding icon
	var icon string
	u.step("Find icon", func() (string, error) {
		ic, err := FindIcon(extractRoot)
		if err != nil {
			return "", err
		}
		icon = ic
		return ic, nil
	})

	//installing icon outside squashfs-root so the entry survives
	appName := strings.TrimSuffix(filepath.Base(destination), filepath.Ext(destination))
	var installedIcon string
	u.step("Install icon", func() (string, error) {
		dst, err := InstallIcon(icon, appName)
		if err != nil {
			return "", err
		}
		installedIcon = dst
		return dst, nil
	})

	//patching extracted desktop file to point at installed paths
	u.step("Fix desktop file", func() (string, error) {
		if err := FixDesktopFile(desktop, destination, installedIcon); err != nil {
			return "", err
		}
		return "", nil
	})

	//generating + installing desktop entry
	u.step("Install desktop entry", func() (string, error) {
		desktopPath, err := GenerateDesktopFile(appName, destination, installedIcon)
		if err != nil {
			return "", err
		}
		if err := InstallDesktopFile(desktopPath); err != nil {
			return "", err
		}
		return desktopPath, nil
	})

	return destination
}

// is the file a regular or normal file
func isRegular(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return info.Mode().IsRegular(), nil
}

func CreateApplicationsDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	newPath := filepath.Join(home, "Applications")
	err = os.MkdirAll(newPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func MoveAppImage(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	if err := os.Chmod(dst, info.Mode()); err != nil {
		return err
	}
	return os.Remove(src)
}

func ExtractAppImage(path string) error {
	if err := os.Chmod(path, 0755); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, path, "--appimage-extract")
	cmd.Dir = filepath.Dir(path)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func FindDesktopFile(root string) (string, error) {
	var found string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".desktop" {
			found = path
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if found == "" {
		return "", fmt.Errorf("no .desktop file found in %s", root)
	}

	return found, nil
}

func FindIcon(root string) (string, error) {
	var found string
	best := 0
	priority := map[string]int{
		".ico": 1,
		".png": 2,
		".svg": 3,
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if p, ok := priority[filepath.Ext(path)]; ok && p > best {
			found = path
			best = p
			if best == 3 {
				return filepath.SkipAll
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if found == "" {
		return "", fmt.Errorf("no icon found")
	}

	return found, nil
}

func InstallIcon(iconSrc, appName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	iconsDir := filepath.Join(home, ".local", "share", "icons")
	if err := os.MkdirAll(iconsDir, 0755); err != nil {
		return "", err
	}

	dst := filepath.Join(iconsDir, appName+filepath.Ext(iconSrc))

	in, err := os.Open(iconSrc)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return "", err
	}
	if err := out.Close(); err != nil {
		return "", err
	}

	return dst, nil
}

func FixDesktopFile(desktopSrc, appPath, iconPath string) error {
	data, err := os.ReadFile(desktopSrc)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Exec=") {
			lines[i] = "Exec=" + appPath
		} else if strings.HasPrefix(trimmed, "Icon=") {
			lines[i] = "Icon=" + iconPath
		}
	}

	return os.WriteFile(desktopSrc, []byte(strings.Join(lines, "\n")), 0644)
}

func InstallDesktopFile(desktopPath string) error {
	if err := os.Chmod(desktopPath, 0755); err != nil {
		return err
	}

	if _, err := exec.LookPath("update-desktop-database"); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		dir := filepath.Dir(desktopPath)
		if err := exec.CommandContext(ctx, "update-desktop-database", dir).Run(); err != nil {
			return err
		}
	}

	return nil
}

func GenerateDesktopFile(name, appPath, iconPath string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	applicationsDir := filepath.Join(
		home,
		".local",
		"share",
		"applications",
	)

	// Make sure the applications directory exists
	err = os.MkdirAll(applicationsDir, 0755)
	if err != nil {
		return "", err
	}

	desktopPath := filepath.Join(
		home,
		".local",
		"share",
		"applications",
		name+".desktop",
	)

	content := fmt.Sprintf(`[Desktop Entry]
Name=%s
Exec=%s
Icon=%s
Type=Application
Terminal=false
Categories=Utility;
`, name, appPath, iconPath)

	if err := os.WriteFile(desktopPath, []byte(content), 0644); err != nil {
		return "", err
	}

	return desktopPath, nil
}

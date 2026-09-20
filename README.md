# 📦 appimg

Install an `.AppImage` like a native Linux app — moves it into `~/Applications`,
extracts its desktop entry + icon, and registers it on your system. Uninstalls cleanly too.

<p>
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go" alt="Go 1.26" />
  <img src="https://img.shields.io/badge/platform-linux-orange?style=flat-square" alt="Linux" />
  <img src="https://img.shields.io/badge/TUI-Bubble_Tea-FF75B7?style=flat-square" alt="Bubble Tea TUI" />
</p>

## ✨ Features

- 🖥️ **Boxed TUI** — pick an AppImage, watch each install step tick by in a bordered interface
- ⌨️   **Simple CLI** — scriptable flags for install / uninstall
- 🧹 **Clean uninstall** — removes the binary, `.desktop` entry, and icon
- 🔁 **Cross-device safe** — falls back to copy when `~/Applications` lives on another filesystem

## 📋 Requirements

- Linux with `~/Applications` writable
- Go 1.26+ (to build)
- Optional: `update-desktop-database` (refreshed automatically when present)

## 🔧 Install

One-liner (builds from source, needs Go + git):

```sh
curl -fsSL https://raw.githubusercontent.com/Ashking-tech/AppImg/main/install.sh | bash
```

Custom directory:

```sh
curl -fsSL https://raw.githubusercontent.com/Ashking-tech/AppImg/main/install.sh | bash -s -- --dir /usr/local/bin
```

Or build manually:

```sh
git clone https://github.com/Ashking-tech/AppImg.git
cd AppImg
go build -o appimg .
```

## 🚀 Usage

**Interactive TUI** (no args, in a terminal):

```sh
./appimg
```

> ↑↓ navigate · / filter · enter select · **u** uninstall · q quit

**Install via CLI:**

```sh
./appimg --filename ~/Downloads/MyApp.AppImage
./appimg -f ~/Downloads/MyApp.AppImage
./appimg ~/Downloads/MyApp.AppImage
```

**Uninstall:**

```sh
./appimg --uninstall MyApp
./appimg -u MyApp
```

Non-interactive shells (pipes, scripts) get plain `[1/9] … ok` output instead of colors.

## ⚙️ How it works

```
Validate → Move to ~/Applications → Extract (--appimage-extract)
  → Find .desktop + icon in squashfs-root
  → Install icon to ~/.local/share/icons
  → Patch desktop Exec/Icon paths
  → Write ~/.local/share/applications/<app>.desktop
  → update-desktop-database
```

Uninstall reverses it: binary + desktop entry + icon, then refreshes the database.

## 🗂️ Layout

| File          | What lives there                        |
|---------------|-----------------------------------------|
| `main.go`     | Core install pipeline + helpers         |
| `cli.go`      | Flag parsing, step output               |
| `tui.go`      | Boxed Bubble Tea interface              |
| `uninstall.go`| `RemoveApp` + target resolution         |

## 🙏 Credits

- **Core pipeline** (`main.go`: validate → move → extract → desktop/icon handling) — **written by me**
- **TUI** (Bubble Tea picker, progress boxes, uninstall screen) — **built with AI assistance**

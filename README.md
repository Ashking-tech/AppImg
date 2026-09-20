# appimg

Install an `.AppImage` like a native Linux app — moves it into `~/Applications`,
extracts its desktop entry + icon, and registers it on your system. Uninstalls cleanly too.

<p>
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go" alt="Go 1.26" />
  <img src="https://img.shields.io/badge/platform-linux-orange?style=flat-square" alt="Linux" />
  <img src="https://img.shields.io/badge/TUI-Bubble_Tea-FF75B7?style=flat-square" alt="Bubble Tea TUI" />
</p>

## Features

- **Clean uninstall** — removes the binary, `.desktop` entry, and icon
- **Cross-device safe** — falls back to copy when `~/Applications` lives on another filesystem

## Requirements

- Linux (any) 
- Go 1.26+ 
- Optional: `update-desktop-database` (refreshed automatically when present)

## Install

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

## Usage

**Interactive TUI** (no args, in a terminal):

```sh
./appimg
```

> Up/Down navigate · / filter · enter select · **u** uninstall · q quit

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

## Architecture

```mermaid
flowchart TD
    A["appimg (no args)\nappimg --filename X"] --> B{Mode?}
    B -->|no args + TTY| C["TUI (tui.go)\nBubble Tea"]
    B -->|flags / piped| D["CLI (cli.go)\nflag parsing"]

    C --> E["Pick: ./ + ~/Downloads\n*.AppImage"]
    E -->|"u key"| F["Uninstall screen\n~/Applications list"]
    E -->|enter| G[Confirm]
    G -->|enter| H[Install pipeline]
    F -->|enter + confirm| I["RemoveApp\n(uninstall.go)"]

    D --> J{--uninstall?}
    J -->|yes| I
    J -->|no| H

    H --> H1["Validate (isRegular)"]
    H1 --> H2["Create ~/Applications"]
    H2 --> H3["Move (rename → copy fallback)"]
    H3 --> H4["Extract (--appimage-extract)"]
    H4 --> H5["Find .desktop + icon\n(squashfs-root)"]
    H5 --> H6["Install icon\n~/.local/share/icons"]
    H6 --> H7["Fix desktop Exec/Icon"]
    H7 --> H8["Write + install .desktop\n~/.local/share/applications"]

    I --> I1["Remove binary\n~/Applications/<app>"]
    I --> I2["Remove .desktop entry"]
    I --> I3["Remove installed icon"]
    H8 --> K["update-desktop-database"]
    I1 & I2 & I3 --> K
```

## Layout

| File          | What lives there                        |
|---------------|-----------------------------------------|
| `main.go`     | Core install pipeline + helpers         |
| `cli.go`      | Flag parsing, step output               |
| `tui.go`      | Boxed Bubble Tea interface              |
| `uninstall.go`| `RemoveApp` + target resolution         |

## Credits

- **Core pipeline** (`main.go`: validate → move → extract → desktop/icon handling) — **written by me**
- **TUI** (Bubble Tea picker, progress boxes, uninstall screen) — **built with AI assistance**

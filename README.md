# fhist - Elaborate Shell History Search with Context Window

`fhist` is an advanced interactive shell history viewer and search tool built in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss). It reads directly from standard shell history logs (`.bash_history`, `.zsh_history`, fish history) and offers a **two-tiered TUI** to view both matching commands and the **N commands before and after** each match.

---

## Key Features

- **Two-Tiered TUI Layout**:
  - **Top Pane**: Real-time search filter and hit list (Fuzzy, Contains, Prefix, Regex matching).
  - **Bottom Pane (Context Preview)**: Displays **N commands before** and **N commands after** the selected history hit in chronological order.
- **Dynamic Context Scaling**: Press `+` or `-` (or `[` / `]`) to dynamically expand or contract the context radius $N$ live inside the TUI.
- **Standard Shell History Support**:
  - Directly loads `~/.bash_history`, `~/.zsh_history`, and `fish_history`.
  - Press `Ctrl+S` to cycle history sources on the fly (`bash`, `zsh`, `fish`, `all`).
- **Interactive Controls**:
  - `Ctrl+R`: Cycle search algorithms (Fuzzy, Contains, Prefix, Regex).
  - `Ctrl+E`: Copy the entire $2N+1$ context window to your clipboard.
  - `Enter`: Select command, copy to clipboard, and exit for immediate shell buffer execution.
  - `Ctrl+O`: Output the selected command AND context window to terminal stdout.

---

## Installation & Setup

`fhist` is compiled and installed at `/home/frans/bin/fhist` and `/home/frans/Utils/fhist`.

### Shell Integration (`Ctrl+R` Keybinding)

To replace standard `Ctrl+R` with `fhist`:

#### Bash (`~/.bashrc`)
```bash
eval "$(fhist --init bash)"
```

#### Zsh (`~/.zshrc`)
```zsh
eval "$(fhist --init zsh)"
```

#### Fish (`~/.config/fish/config.fish`)
```fish
fhist --init fish | source
```

---

## Command Line Options

```bash
fhist [OPTIONS] [QUERY]
```

- `-n, --context N`: Number of commands before/after hit in context view (default: `5`).
- `-q, --query TEXT`: Initial search query filter.
- `-s, --source SOURCE`: History source (`auto`, `bash`, `zsh`, `fish`, `all`).
- `-f, --file PATH`: Read from custom history file path.
- `--print-context`: Output selected command along with surrounding context window on exit.
- `--init SHELL`: Generate shell integration script (`bash`, `zsh`, `fish`).

---

## Keyboard Shortcuts inside TUI

| Key Binding | Action |
| :--- | :--- |
| `↑` / `↓` or `Ctrl+P` / `Ctrl+N` | Navigate search results up / down |
| `PgUp` / `PgDn` | Scroll search list page up / page down |
| `+` / `-` or `[` / `]` | Increase / decrease context radius $N$ |
| `Ctrl+R` | Toggle search mode (Fuzzy $\rightarrow$ Contains $\rightarrow$ Prefix $\rightarrow$ Regex) |
| `Ctrl+S` | Cycle history source (`bash` $\rightarrow$ `zsh` $\rightarrow$ `fish` $\rightarrow$ `all`) |
| `Ctrl+E` | Copy entire context window block to clipboard |
| `Enter` | Select command (copies to clipboard, prints to stdout, and exits) |
| `Ctrl+O` | Select command & print full context window on exit |
| `Esc` / `Ctrl+C` | Exit without selecting |

---

## License

This project is open-source software licensed under the [MIT License](LICENSE).

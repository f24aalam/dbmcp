# dbmcp (godbmcp)

`godbmcp` is a command-line tool to manage database connections and run an **MCP (Model Context Protocol) server** for databases.

It allows you to:

* Add and manage database connections
* List configured databases
* Start an MCP server backed by those connections

---

## 🚀 Installation

### 📦 Prebuilt binaries (Recommended)

You **do not need Go installed**. Just run the installer for your platform.

---

### 🪟 Windows (PowerShell)

Run PowerShell **as your normal user**:

```powershell
iwr https://raw.githubusercontent.com/f24aalam/dbmcp/main/install.ps1 | iex
```

After installation, restart your terminal and verify:

```powershell
godbmcp --help
```

---

### 🐧 Linux / 🍎 macOS

```bash
curl -fsSL https://raw.githubusercontent.com/f24aalam/dbmcp/master/install.sh | bash
```

Then verify:

```bash
godbmcp --help
```

---

### 📍 Installation Location

| OS            | Path                             |
| ------------- | -------------------------------- |
| Windows       | `C:\Users\<you>\bin\godbmcp.exe` |
| Linux / macOS | `/usr/local/bin/godbmcp`         |

---

## 🧰 Usage

```bash
godbmcp [command]
```

### Available Commands

| Command      | Description                   |
| ------------ | ----------------------------- |
| `add`        | Add a new database connection |
| `list`       | List all saved connections    |
| `mcp`        | Start the MCP server          |
| `completion` | Generate shell autocompletion |
| `help`       | Show help for any command     |

---

### Example: Add a database connection

```bash
godbmcp add
```

---

### Example: List connections

```bash
godbmcp list
```

---

### Example: Start MCP server

```bash
godbmcp mcp
```

---

## 🧠 Shell Autocompletion (Optional)

### Bash

```bash
godbmcp completion bash > /etc/bash_completion.d/godbmcp
```

### Zsh

```bash
godbmcp completion zsh > "${fpath[1]}/_godbmcp"
```

Restart your shell to activate completions.

---

## 🛠️ Build from Source (Optional)

### Requirements

* Go 1.20+

### Build

```bash
git clone https://github.com/f24aalam/dbmcp.git
cd dbmcp
go build -o godbmcp ./cmd/godbmcp
```

---

## 📄 License

MIT License. See [LICENSE](./LICENSE) for details.

---

## ✅ Status

* ✔ Cross-platform binaries (Windows, Linux, macOS)
* ✔ amd64 + arm64 supported
* ✔ One-line install scripts
* ✔ No runtime dependencies

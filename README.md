```
       ____                       
  ____/ / /_  ____ ___  _________ 
 / __  / __ \/ __ `__ \/ ___/ __ \
/ /_/ / /_/ / / / / / / /__/ /_/ /
\__,_/_.___/_/ /_/ /_/\___/ .___/ 
                         /_/      
```

**Database MCP server and CLI**

## Why dbmcp?

Other database MCP tools often require:

- Cloning a whole repo
- Node.js / npm
- Credentials only via environment variables

**dbmcp** keeps it simple:

- Single static binary — no extra runtime
- Credentials stored locally (OS keychain where available)
- **`dbmcp init`** for a guided first-time setup in your project
- **`dbmcp add`** when you just want to register a connection
- Use a **connection ID** in MCP clients, or let **init** write agent config files under your repo

---

## Easy setup

Install **dbmcp**, then run **init** from your application’s project root.

### 1. Install

**Linux / macOS**

```bash
curl -fsSL https://raw.githubusercontent.com/f24aalam/dbmcp/master/install.sh | bash
```

Add `~/.local/bin` to `PATH` if the installer puts the binary there.

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/f24aalam/dbmcp/master/install.ps1 | iex
```

### 2. Run init

```bash
cd /path/to/your/project
dbmcp init
```

You’ll see the **dbmcp** banner, then a step-by-step wizard:

1. **How to connect** — use a **saved connection**, **scan the project** for DB hints, or **add a new** URL/path.
2. **Review** connection details, then **dbmcp** opens a real connection (**ping**). If that fails, you’ll see an error and nothing is saved.
3. **Save** — for new or scanned connections, you’re asked **“Save this connection to dbmcp?”** (default **Yes**). Choose **No** to verify only; the flow ends without MCP install (no stored connection ID). **Existing** connections are only validated, not re-saved.
4. Optionally **install MCP config** for agents in **this repo**. If you choose **No**, the flow ends there (no agent picker).
5. If **Yes**, pick agents; configs are **merged** into existing files without wiping other servers.

**URL normalization (init / scan / pasted URLs in init):** **`mysql://...`** values are converted to the DSN form the Go MySQL driver expects. **PostgreSQL** URLs from Prisma and similar tools often include **`?schema=public`** (or `connection_limit`, `pool_timeout`); those keys are **removed** before connecting, because they are not valid **lib/pq** connection options. The database name in the path is unchanged.

Then use **`dbmcp list`** for connection IDs, or open the files **init** created (see table below).

> **Command in generated configs:** MCP snippets use **`godbmcp`** as the executable name. The install script typically installs **`dbmcp`**. Symlink, rename, or edit generated JSON/TOML if your binary name differs.

---

## `dbmcp init`

| | |
|--|--|
| **What it does** | Connects this directory to a database (discover, reuse, or enter), **pings** to validate, optionally **saves**, then optionally adds **project-local** MCP entries for coding agents. |
| **Where to run** | Repository / app root — paths like `.cursor/mcp.json` are written relative to `cwd`. |

**Connection paths**

| Mode | Behavior |
|------|----------|
| **Use existing connections** | Pick from **dbmcp**’s saved list; **no** re-save; ping only. |
| **Scan this project** | Search for DB URLs / config-like keys; pick a candidate or enter details manually. |
| **Add new connection** | Prompts for name, type (MySQL / PostgreSQL / SQLite), and URL or file path. |

**Optional MCP install — files touched**

| Agent | File (under project root) |
|-------|---------------------------|
| Claude Code | `.mcp.json` |
| Cursor | `.cursor/mcp.json` |
| Copilot (VS Code) | `.vscode/mcp.json` |
| Copilot CLI | `.copilot/mcp-config.json` |
| Codex | `.codex/config.toml` |
| OpenCode | `opencode.json` |

---

## Installation (reference)

Same as [Easy setup → Install](#1-install).

### Uninstall

**Linux / macOS**

```bash
rm ~/.local/bin/dbmcp
# or
sudo rm /usr/local/bin/dbmcp
```

**Windows**

```powershell
Remove-Item "$HOME\bin\dbmcp.exe"
```

---

## Quick start (without init)

### 1. Add a connection

```bash
dbmcp add
```

**Examples**

| Engine | Example |
|--------|---------|
| MySQL | `user:password@tcp(localhost:3306)/mydb` |
| PostgreSQL | `postgres://user:password@localhost:5432/mydb` |
| SQLite | `/home/user/data/mydb.sqlite` |

Prisma-style Postgres URLs such as `.../mydb?schema=public` work in **`dbmcp init`** (and any path that uses the same URL resolver) because **`schema=`** is stripped before connect. The plain **`dbmcp add`** form currently saves the string you type as-is—use a URL without ORM-only query params there, or run **init** instead.

### 2. List connections

```bash
dbmcp list
```

### 3. Start the MCP server

```bash
dbmcp mcp --connection-id <CONNECTION_ID>
```

---

## Usage in MCP clients

### Cursor (global example)

`~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "dbmcp": {
      "command": "dbmcp",
      "args": ["mcp", "--connection-id", "<YOUR_CONNECTION_ID>"]
    }
  }
}
```

Use the same **`command`** / **`args`** shape in other clients; **init** writes the equivalent under your project using **`godbmcp`** as noted above.

### OpenCode

Point OpenCode at the same server definition (stdio, `mcp` subcommand + `--connection-id`).

---

## Commands

| Command | Description |
|---------|-------------|
| `dbmcp init` | Guided setup: connection (existing / scan / new), validate, optional per-agent MCP config in this repo |
| `dbmcp add` | Add a new database connection |
| `dbmcp list` | List saved connections |
| `dbmcp mcp` | Run the MCP server |
| `dbmcp completion` | Shell completion scripts |

---

## Supported databases

- MySQL  
- PostgreSQL  
- SQLite  

---

## License

MIT

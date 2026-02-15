# godbmcp

Database MCP server with easy credential management.

## Why godbmcp?

Other database MCP tools require:
- Cloning the entire repository
- Node.js/npm installed
- Setting database credentials in environment variables

**godbmcp** makes it simple:
- Prebuilt binary - no dependencies needed
- Store credentials securely locally
- Just run `godbmcp add` to add a connection
- Copy the generated connection ID and use it in your MCP client

## Installation

### Linux/macOS
```bash
curl -fsSL https://raw.githubusercontent.com/f24aalam/dbmcp/master/install.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/f24aalam/dbmcp/master/install.ps1 | iex
```

## Quick Start

### 1. Add a database connection
```bash
godbmcp add
```

**MySQL example:**
```
user:password@tcp(localhost:3306)/mydb
```

**PostgreSQL example:**
```
postgres://user:password@localhost:5432/mydb
```

**SQLite example:**
```
/home/user/data/mydb.sqlite
```

Follow the prompts to enter your database details (MySQL or SQLite).

### 2. List connections
```bash
godbmcp list
```

### 3. Start MCP server
```bash
godbmcp mcp --connection-id <CONNECTION_ID>
```

## Usage in MCP Clients

### Cursor
Add to `~/.cursor/mcp.json`:
```json
{
  "mcpServers": {
    "godbmcp": {
      "command": "godbmcp",
      "args": ["mcp", "--connection-id", "<YOUR_CONNECTION_ID>"]
    }
  }
}
```

### OpenCode
Add to your OpenCode MCP settings with the same command format.

## Commands

| Command | Description |
|---------|-------------|
| `godbmcp add` | Add a new database connection |
| `godbmcp list` | List all saved connections |
| `godbmcp mcp` | Start the MCP server |
| `godbmcp completion` | Generate shell autocompletion |

## Supported Databases

- MySQL
- PostgreSQL
- SQLite

## License

MIT

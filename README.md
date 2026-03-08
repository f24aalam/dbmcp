# dbmcp

Database MCP server with easy credential management.

## Why dbmcp?

Other database MCP tools require:
- Cloning the entire repository
- Node.js/npm installed
- Setting database credentials in environment variables

**dbmcp** makes it simple:
- Prebuilt binary - no dependencies needed
- Store credentials securely locally
- Just run `dbmcp add` to add a connection
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
dbmcp add
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
dbmcp list
```

### 3. Start MCP server
```bash
dbmcp mcp --connection-id <CONNECTION_ID>
```

## Usage in MCP Clients

### Cursor
Add to `~/.cursor/mcp.json`:
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

### OpenCode
Add to your OpenCode MCP settings with the same command format.

## Commands

| Command | Description |
|---------|-------------|
| `dbmcp add` | Add a new database connection |
| `dbmcp list` | List all saved connections |
| `dbmcp mcp` | Start the MCP server |
| `dbmcp completion` | Generate shell autocompletion |

## Supported Databases

- MySQL
- PostgreSQL
- SQLite

## License

MIT

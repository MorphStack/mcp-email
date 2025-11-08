# MCP Email Server

A Model Context Protocol (MCP) server for handling SMTP/IMAP email operations. This server allows you to search past emails, send new emails, and manage multiple email accounts.

## Features

- **Search Emails**: Flexible search across all email fields (sender, recipient, subject, body, date range)
- **Send Emails**: Send emails with support for text, HTML, attachments, CC, and BCC
- **Multi-Account Support**: Manage multiple email accounts simultaneously
- **Local Caching**: SQLite-based cache for fast email searches
- **Full-Text Search**: Fast full-text search using SQLite FTS5
- **Generic IMAP/SMTP**: Works with any email provider that supports IMAP/SMTP

## Requirements

- Go 1.23 or higher
- Docker (for containerized deployment)
- Email account with IMAP/SMTP access (app passwords recommended)

## Configuration

Configuration is done via environment variables. You can configure either a single account or multiple accounts.

**See `.env.example` for a complete example configuration file.**

### Single Account Configuration

```bash
IMAP_HOST=imap.example.com
IMAP_PORT=993
IMAP_USERNAME=user@example.com
IMAP_PASSWORD=your_app_password
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=user@example.com
SMTP_PASSWORD=your_app_password
ACCOUNT_NAME=default
```

### Multiple Accounts Configuration

```bash
# Account 1
ACCOUNT_1_NAME=work
ACCOUNT_1_IMAP_HOST=imap.work.com
ACCOUNT_1_IMAP_PORT=993
ACCOUNT_1_IMAP_USERNAME=user@work.com
ACCOUNT_1_IMAP_PASSWORD=your_app_password
ACCOUNT_1_SMTP_HOST=smtp.work.com
ACCOUNT_1_SMTP_PORT=587
ACCOUNT_1_SMTP_USERNAME=user@work.com
ACCOUNT_1_SMTP_PASSWORD=your_app_password

# Account 2
ACCOUNT_2_NAME=personal
ACCOUNT_2_IMAP_HOST=imap.gmail.com
ACCOUNT_2_IMAP_PORT=993
ACCOUNT_2_IMAP_USERNAME=user@gmail.com
ACCOUNT_2_IMAP_PASSWORD=your_app_password
ACCOUNT_2_SMTP_HOST=smtp.gmail.com
ACCOUNT_2_SMTP_PORT=587
ACCOUNT_2_SMTP_USERNAME=user@gmail.com
ACCOUNT_2_SMTP_PASSWORD=your_app_password
```

### Optional Configuration

```bash
CACHE_PATH=/data/email_cache.db
SEARCH_RESULT_LIMIT=100
LOG_LEVEL=info
```

### Common Email Provider Settings

#### Gmail
- IMAP: `imap.gmail.com:993`
- SMTP: `smtp.gmail.com:587`
- **Note**: Requires App Password (not regular password)

#### Outlook/Office365
- IMAP: `outlook.office365.com:993`
- SMTP: `smtp.office365.com:587`

#### Yahoo
- IMAP: `imap.mail.yahoo.com:993`
- SMTP: `smtp.mail.yahoo.com:587`

## MCP Tools

### `list_folders`
List available mailboxes/folders for configured email accounts.

**Parameters:**
- `account_name` (optional): Specific account name, or all accounts if omitted

### `search_emails`
Search cached emails with flexible filters.

**Parameters:**
- `account_name` (optional): Filter by specific account
- `folder` (optional): Filter by folder/mailbox
- `sender` (optional): Filter by sender email/name
- `recipient` (optional): Filter by recipient email
- `subject` (optional): Filter by subject (substring match)
- `body` (optional): Filter by body content (full-text search)
- `date_from` (optional): Start date (ISO 8601 format)
- `date_to` (optional): End date (ISO 8601 format)
- `limit` (optional): Result limit (default: 100, max: 1000)

### `get_email`
Retrieve full email by ID from cache or IMAP.

**Parameters:**
- `email_id` (required): Email ID (from search results)
- `account_name` (optional): Account name if needed

### `send_email`
Send a new email with support for text, HTML, attachments, CC, BCC.

**Parameters:**
- `account_name` (required): Account to send from
- `to` (required): Recipient email address(es) (comma-separated)
- `cc` (optional): CC recipients (comma-separated)
- `bcc` (optional): BCC recipients (comma-separated)
- `subject` (required): Email subject
- `body_text` (optional): Plain text body
- `body_html` (optional): HTML body
- `attachments` (optional): Array of attachment paths/URLs
- `reply_to` (optional): Reply-To header
- `in_reply_to` (optional): In-Reply-To header (for replies)

## Building

### Local Build

```bash
go mod download
go build -o mcp-email-server ./cmd/server
```

### Docker Build

```bash
docker build -t mcp-email-server .
```

## Running

### Local Run

```bash
./mcp-email-server
```

### Docker Run

```bash
docker run -it --rm \
  -e IMAP_HOST=imap.example.com \
  -e IMAP_PORT=993 \
  -e IMAP_USERNAME=user@example.com \
  -e IMAP_PASSWORD=your_app_password \
  -e SMTP_HOST=smtp.example.com \
  -e SMTP_PORT=587 \
  -e SMTP_USERNAME=user@example.com \
  -e SMTP_PASSWORD=your_app_password \
  -v $(pwd)/data:/data \
  mcp-email-server
```

Alternatively, you can use an environment file:

```bash
# Copy the example file
cp .env.example .env
# Edit .env with your credentials
# Then run:
docker run -it --rm --env-file .env -v $(pwd)/data:/data mcp-email-server
```

## MCP Configuration

To use this server with Claude Desktop or VS Code, you need to configure it in your `mcp.json` file.

### Claude Desktop Configuration

**Location:** 
- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`
- Linux: `~/.config/Claude/claude_desktop_config.json`

### VS Code Configuration

**Location:** `.vscode/mcp.json` in your project directory

### Configuration Examples

#### 1. Local Execution (Single Account)

Use this if you've built the server locally:

```json
{
  "mcpServers": {
    "mcp-email": {
      "command": "/absolute/path/to/mcp-email-server",
      "env": {
        "IMAP_HOST": "imap.gmail.com",
        "IMAP_PORT": "993",
        "IMAP_USERNAME": "your-email@gmail.com",
        "IMAP_PASSWORD": "your-app-password",
        "SMTP_HOST": "smtp.gmail.com",
        "SMTP_PORT": "587",
        "SMTP_USERNAME": "your-email@gmail.com",
        "SMTP_PASSWORD": "your-app-password",
        "CACHE_PATH": "/tmp/email_cache.db",
        "SEARCH_RESULT_LIMIT": "100",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

**See `mcp.json.example.local` for a complete example.**

#### 2. Docker Execution (Single Account)

Use this to run the server in a Docker container:

```json
{
  "mcpServers": {
    "mcp-email": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "-v", "/tmp:/data",
        "-e", "IMAP_HOST=imap.gmail.com",
        "-e", "IMAP_PORT=993",
        "-e", "IMAP_USERNAME=your-email@gmail.com",
        "-e", "IMAP_PASSWORD=your-app-password",
        "-e", "SMTP_HOST=smtp.gmail.com",
        "-e", "SMTP_PORT=587",
        "-e", "SMTP_USERNAME=your-email@gmail.com",
        "-e", "SMTP_PASSWORD=your-app-password",
        "-e", "CACHE_PATH=/data/email_cache.db",
        "-e", "SEARCH_RESULT_LIMIT=100",
        "-e", "LOG_LEVEL=info",
        "mcp-email-server:latest"
      ]
    }
  }
}
```

**See `mcp.json.example.docker` for a complete example.**

**Note:** Make sure you've built the Docker image first:
```bash
docker build -t mcp-email-server:latest .
```

#### 3. Multi-Account Configuration

For multiple email accounts, use the `ACCOUNT_N_*` environment variables:

```json
{
  "mcpServers": {
    "mcp-email": {
      "command": "/absolute/path/to/mcp-email-server",
      "env": {
        "ACCOUNT_1_NAME": "work",
        "ACCOUNT_1_IMAP_HOST": "imap.work.com",
        "ACCOUNT_1_IMAP_PORT": "993",
        "ACCOUNT_1_IMAP_USERNAME": "user@work.com",
        "ACCOUNT_1_IMAP_PASSWORD": "your-app-password",
        "ACCOUNT_1_SMTP_HOST": "smtp.work.com",
        "ACCOUNT_1_SMTP_PORT": "587",
        "ACCOUNT_1_SMTP_USERNAME": "user@work.com",
        "ACCOUNT_1_SMTP_PASSWORD": "your-app-password",
        "ACCOUNT_2_NAME": "personal",
        "ACCOUNT_2_IMAP_HOST": "imap.gmail.com",
        "ACCOUNT_2_IMAP_PORT": "993",
        "ACCOUNT_2_IMAP_USERNAME": "user@gmail.com",
        "ACCOUNT_2_IMAP_PASSWORD": "your-app-password",
        "ACCOUNT_2_SMTP_HOST": "smtp.gmail.com",
        "ACCOUNT_2_SMTP_PORT": "587",
        "ACCOUNT_2_SMTP_USERNAME": "user@gmail.com",
        "ACCOUNT_2_SMTP_PASSWORD": "your-app-password",
        "CACHE_PATH": "/tmp/email_cache.db",
        "SEARCH_RESULT_LIMIT": "100",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

**See `mcp.json.example.multi-account` for a complete example.**

### Example Files

- `mcp.json.example.local` - Local execution with single account
- `mcp.json.example.docker` - Docker execution with single account (env vars in args)
- `mcp.json.example.docker-envfile` - Docker execution using an env file
- `mcp.json.example.multi-account` - Local execution with multiple accounts

#### Using Docker with Environment File

For a cleaner Docker configuration, you can use an environment file:

```json
{
  "mcpServers": {
    "mcp-email": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "--env-file", "/absolute/path/to/.env",
        "-v", "/tmp:/data",
        "mcp-email-server:latest"
      ]
    }
  }
}
```

**See `mcp.json.example.docker-envfile` for a complete example.**

This approach keeps your credentials in a separate `.env` file (which should not be committed to version control).

### Setup Instructions

1. **Choose your configuration** (local or Docker)
2. **Copy the appropriate example file** to your MCP config location
3. **Update the paths and credentials** with your actual values
4. **Restart Claude Desktop or VS Code** to load the configuration

### Security Notes

- Never commit your actual `mcp.json` with real passwords to version control
- Use app passwords instead of your regular email password
- Consider using environment variables or secrets management for production

## Development

### Project Structure

```
mcp-email/
├── cmd/server/          # Main application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── email/           # IMAP/SMTP client implementations
│   ├── cache/           # SQLite cache layer
│   ├── mcp/             # MCP server implementation
│   └── tools/           # MCP tool implementations
├── pkg/types/           # Shared data types
└── Dockerfile           # Docker container definition
```

### Testing

```bash
go test ./...
```

### Linting

```bash
golangci-lint run --timeout=5m
```

## CI/CD

This project uses GitHub Actions for continuous integration and releases.

### CI Workflow

- Runs on every push and pull request
- Tests on multiple Go versions (1.21, 1.22, 1.23)
- Runs linters and code quality checks
- Builds and validates on multiple platforms
- Builds and tests Docker image

### Release Workflow

- Automatically triggered on version tags (e.g., `v1.0.0`)
- Validates semantic versioning (must be greater than previous tag)
- Builds binaries for all platforms (Linux, macOS, Windows - amd64, arm64)
- Publishes Docker images to GitHub Container Registry (and optionally Docker Hub)
- Creates GitHub releases with all artifacts

### Creating a Release

1. Create and push a version tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. The release workflow will automatically:
   - Validate the version
   - Build all binaries
   - Create Docker images
   - Publish to registries
   - Create GitHub release

See [RELEASE.md](RELEASE.md) and [.github/SETUP.md](.github/SETUP.md) for detailed information.

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.


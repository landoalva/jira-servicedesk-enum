# jira-servicedesk-enum

A Go tool for enumerating Atlassian Jira Service Desk.
Our user permissions and fetching users from service desk portals using session cookie authentication.
Useful for password spraying and penetration testing. Brought to you by the [RasterSec](https://www.rastersec.com) team 🙌.

## Installation

```bash
go install github.com/RasterSec/jira-servicedesk-enum@latest
```

Or build from source:

```bash
go build
```

## Authentication

This tool uses the `customer.account.session.token` cookie for authentication.

## Usage

### Signup

Trigger service desk signup:

```bash
./jira-servicedesk-enum signup --url https://example.atlassian.net --email user@example.com
```

### Check Permissions

Check what permissions we have:

```bash
./jira-servicedesk-enum permissions \
  --url https://example.atlassian.net \
  --cookie "secret..."
```

### Enumerate Users

List all users across accessible service desks:

```bash
./jira-servicedesk-enum users \
  --url https://example.atlassian.net \
  --cookie "secret..."
```

Target a specific service desk by ID:

```bash
./jira-servicedesk-enum users \
  --url https://example.atlassian.net \
  --cookie "secret..." \
  --desk 123
```

Limit results per service desk:

```bash
./jira-servicedesk-enum users \
  --url https://example.atlassian.net \
  --cookie "secret..." \
  --max 1000
```

Search with a custom query:

```bash
./jira-servicedesk-enum users \
  --url https://example.atlassian.net \
  --cookie "secret..." \
  --query "john"
```

## License

Licensed under the Apache License, Version 2.0.

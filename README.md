# Blacklight

A lightweight API probing tool built in Go.  
It parses an OpenAPI specification, enumerates endpoints, and runs three types of probes:  

- **Unauthenticated requests**  
- **Authenticated requests** (using a supplied JWT)  
- **Auth bypass attempts** (using common header tricks and fake tokens)  

It also flags potential **IDOR (Insecure Direct Object Reference)** candidates based on path patterns.  

---

## Features
- Parses endpoints from OpenAPI spec (`paths` section).  
- Sends requests with and without authentication.  
- Tests bypass headers such as `X-Forwarded-For: 127.0.0.1`.  
- Logs potential IDOR candidates in yellow.  
- Generates reports:  
  - `unauth_report.txt` — unauthenticated probe results  
  - `auth_report.txt` — authenticated probe results  
  - `bypass_report.txt` — auth bypass probe results  
---

## Usage
```bash
go run main.go \
  --spec ./openapi.json \
  --base-url https://api.example.com \
  --token <valid_jwt_token> \
  --outdir ./reports
```

Or using a cookie instead of an Authorization header:
```bash
go run main.go \
  --spec ./openapi.json \
  --base-url https://api.example.com \
  --cookie "session=abcd1234" \
  --outdir ./reports
```

## Arguments
--spec — Path to OpenAPI spec JSON file
--base-url — Base URL of the target API
--token — Valid JWT token for authenticated probes
--cookie — Session cookie string for authenticated probes
--outdir — Output directory for reports (default: reports)


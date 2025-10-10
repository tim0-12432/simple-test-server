<div align="center">
    <img width="35%" alt="logo" src=""/>
</div>
<div align="center">
    <h1>Simple Test Server</h1>
    <span>This is the second version of simple test servers. This time written in Go and with a tiny cute UI.</span>
</div>

---

## Contents

1. [About](#about)
2. [Usage](#usage)
3. [Development](#development)
4. [Contributing](#contributing)
5. [License](#license)

---

## About

## Usage

```bash
docker run -it --rm --name simple-test-server -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock ghcr.io/tim0-12432/simple-test-server-2:latest
```

## Features

### SMB File Server Tab
The SMB tab allows you to quickly spin up a Samba file server for testing file sharing protocols. It uses the `ghcr.io/servercontainers/samba:smbd-only-latest` Docker image and exposes ports 139 and 445 for SMB/CIFS connectivity. Create an SMB container through the web interface to test file sharing scenarios, monitor connection logs, and manage the server lifecycle.

## Development

During frontend development the Vite dev server may run on a different port than the backend. You can override the backend base URL used by the frontend by setting the environment variable `VITE_BACKEND_URL` before starting the dev server. Example:

```bash
# point frontend dev server to a backend running on port 8080
VITE_BACKEND_URL=http://localhost:8080 npm run dev
```

If `VITE_BACKEND_URL` is not set, the frontend defaults to `http://localhost:8000` in development mode and `window.location.origin` in production builds.

### Running Tests

**Backend Tests:**
```bash
go test ./...
go test -race ./...  # with race detection
```

**Frontend Tests:**
```bash
cd frontend
npm test           # or bun test
npm run test       # watch mode
```

### Environment Variables
Configure the application using an `app.env` file or environment variables:
- `HOST` - Server host (default: 0.0.0.0)
- `PORT` - Server port (default: 8000)  
- `ENV` - Environment mode (default: PROD)
- `ADMIN_USER` - Database admin email (default: admin@hosting.test)
- `ADMIN_PASS` - Database admin password (default: pleaseChange123!)
- `UPLOAD_MAX_BYTES` - Max file upload size in bytes (default: 10MB)

**Required:** Docker socket access (`/var/run/docker.sock`) for container management.


## Contributing

## License

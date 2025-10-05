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

## Development

During frontend development the Vite dev server may run on a different port than the backend. You can override the backend base URL used by the frontend by setting the environment variable `VITE_BACKEND_URL` before starting the dev server. Example:

```bash
# point frontend dev server to a backend running on port 8080
VITE_BACKEND_URL=http://localhost:8080 npm run dev
```

If `VITE_BACKEND_URL` is not set, the frontend defaults to `http://localhost:8080` in development mode and `window.location.origin` in production builds.


## Contributing

## License

Watcher for development

This `dev` helper watches Go files, rebuilds and restarts the server when source files change, and then executes a list of HTTP requests (useful for quick integration testing).

Files
- `dev/watch-dev.js` - Node script that watches `**/*.go`, builds `main.go` into `dev_server`, runs it, and triggers requests.
- `dev/requests.json` - Array of request objects that are executed after the server starts.

Usage
1. Ensure Node.js is installed.
2. From project root run: `node dev/watch-dev.js`

Notes
- The script builds `main.go` and places the binary next to the project root as `dev_server` (`dev_server.exe` on Windows).
- Customize `dev/requests.json` to add/remove requests. Each entry supports:
  - `method` (GET, POST, etc.)
  - `url` (absolute URL)
  - `headers` (object)
  - `body` (object or string)
  - `logBody` (boolean) - log response body
- The watcher ignores `vendor`, `pb_data`, `dev`, and `node_modules` directories.
- On Windows the script uses `taskkill` to ensure the process tree is terminated.

Limitations and improvements
- This is a minimal script. For more features consider using `nodemon`, `reflex`, or `air` (Go live-reload tools).
- You may prefer to add an npm script in `package.json` to run this with `npm run dev:go`.

#!/usr/bin/env node
// Watch Go files, restart server, then run HTTP requests
// Usage: node dev/watch-dev.js

const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');
const http = require('http');

const projectRoot = path.resolve(__dirname, '..');
const serverMain = path.join(projectRoot, 'main.go');
const goCmd = process.platform === 'win32' ? 'go.exe' : 'go';

let serverProc = null;
let restarting = false;

function log(...args) {
  console.log(new Date().toISOString(), ...args);
}

function startServer() {
  if (serverProc) return;
  log('Building and starting Go server...');
  // Build first to ensure compile errors are shown before running
  const exe = process.platform === 'win32' ? path.join(projectRoot, 'dev_server.exe') : path.join(projectRoot, 'dev_server');
  const build = spawn(goCmd, ['build', '-o', exe, serverMain], { cwd: projectRoot, stdio: 'inherit' });
  build.on('exit', (code) => {
    if (code !== 0) {
      log('Build failed, not starting server (exit', code + ')');
      return;
    }
    serverProc = spawn(exe, [], { cwd: projectRoot, stdio: 'inherit' });
    serverProc.on('exit', (code, signal) => {
      log('Server exited', code, signal);
      serverProc = null;
    });
    serverProc.on('error', (err) => {
      log('Server process error', err);
      serverProc = null;
    });
    // Give server some time to start, then run requests
    setTimeout(runRequests, 3000);
  });
}

function stopServer(cb) {
  if (!serverProc) return cb && cb();
  log('Stopping server...');
  if (process.platform === 'win32') {
    // on windows, child.kill may not kill process tree; use taskkill
    const killer = spawn('taskkill', ['/PID', serverProc.pid.toString(), '/T', '/F']);
    killer.on('exit', () => {
      serverProc = null;
      cb && cb();
    });
  } else {
    serverProc.kill('SIGINT');
    serverProc.on('exit', () => {
      serverProc = null;
      cb && cb();
    });
  }
}

function restartServer() {
  if (restarting) return;
  restarting = true;
  stopServer(() => {
    // small delay to release ports
    setTimeout(() => {
      startServer();
      restarting = false;
    }, 200);
  });
}

function runRequests() {
  const requestsFile = path.join(projectRoot, 'dev', 'requests.json');
  if (!fs.existsSync(requestsFile)) {
    log('No requests.json found at', requestsFile);
    return;
  }
  let data;
  try {
    data = JSON.parse(fs.readFileSync(requestsFile, 'utf8'));
  } catch (e) {
    log('Failed to parse requests.json:', e.message);
    return;
  }
  if (!Array.isArray(data)) return;

  (async () => {
    for (const req of data) {
      try {
        await doRequest(req);
      } catch (e) {
        log('Request failed', e.message || e);
      }
    }
  })();
}

function doRequest(req) {
  return new Promise((resolve, reject) => {
    const url = new URL(req.url);
    const opts = {
      method: req.method || 'GET',
      hostname: url.hostname,
      port: url.port || (url.protocol === 'https:' ? 443 : 80),
      path: url.pathname + url.search,
      headers: req.headers || { 'Content-Type': 'application/json' },
    };

    const proto = url.protocol === 'https:' ? require('https') : http;
    const r = proto.request(opts, (res) => {
      let body = '';
      res.setEncoding('utf8');
      res.on('data', (chunk) => (body += chunk));
      res.on('end', () => {
        log('=>', req.method || 'GET', req.url, '->', res.statusCode);
        if (req.logBody) {
          console.log(body);
        }
        resolve({ status: res.statusCode, body });
      });
    });
    r.on('error', reject);
    if (req.body) {
      const payload = typeof req.body === 'string' ? req.body : JSON.stringify(req.body);
      r.write(payload);
    }
    r.end();
  });
}

// Watch for .go file changes
const chokidar = require('chokidar');
const watcher = chokidar.watch('**/*.go', { cwd: projectRoot, ignored: ['**/vendor/**', '**/pb_data/**', '**/dev/**', '**/node_modules/**'] });

watcher.on('ready', () => {
  log('Watching .go files for changes...');
  startServer();
});

watcher.on('change', (p) => {
  log('File changed:', p);
  restartServer();
});

process.on('SIGINT', () => {
  log('Received SIGINT, shutting down...');
  watcher.close();
  stopServer(() => process.exit(0));
});

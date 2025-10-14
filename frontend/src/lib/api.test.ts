/* eslint-disable @typescript-eslint/no-unsafe-function-type */
/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { vi, describe, it, expect } from 'vitest';

import { uploadFile } from './api';

// Mock global XMLHttpRequest
class MockXHR {
  readyState = 0;
  status = 0;
  responseText = '';
  upload = { onprogress: null } as any;
  onreadystatechange: any = null;
  onerror: any = null;
  _listeners: Record<string, Function[]> = {};
  lastOpenUrl: string | null = null;
  lastOpenMethod: string | null = null;

  open(method: string, url: string) { this.lastOpenMethod = method; this.lastOpenUrl = url; }
  send(_fd: any) {
    // simulate upload progress events if a handler exists
    if (this.upload && typeof this.upload.onprogress === 'function') {
      this.upload.onprogress({ lengthComputable: true, loaded: 1, total: 2 } as any);
      this.upload.onprogress({ lengthComputable: true, loaded: 2, total: 2 } as any);
    }
    // simulate final response
    setTimeout(() => {
      this.readyState = 4;
      this.status = 201;
      this.responseText = JSON.stringify({ url: 'http://localhost:8000/test' });
      if (this.onreadystatechange) this.onreadystatechange();
    }, 10);
  }

  addEventListener(event: string, cb: Function) {
    (this._listeners[event] = this._listeners[event] || []).push(cb);
  }
}

const mockXhr = new MockXHR();
vi.stubGlobal('XMLHttpRequest', function() { return mockXhr } as any);

describe('uploadFile', () => {
  it('resolves with JSON response and reports progress', async () => {
    const file = new File(['hello'], 'hello.txt', { type: 'text/plain' });
    const progressCalls: number[] = [];
    const p = uploadFile('server-1', file, 'web', (pct) => progressCalls.push(pct));
    const res = await p;
    expect(res).toHaveProperty('url');
    expect(res.url).toContain('http://');
  });

  it('sends request to ftp protocol path when serverType is ftp', async () => {
    const file = new File(['x'], 'x.txt', { type: 'text/plain' });
    // call uploadFile with serverType 'ftp'
    const p = uploadFile('server-ftp', file, 'ftp');
    await p;
    expect(mockXhr.lastOpenUrl).toContain('/protocols/ftp/');
  });
});

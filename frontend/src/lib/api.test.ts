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

  open(_method: string, _url: string) {}
  send(_fd: any) {
    // simulate progress
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

vi.stubGlobal('XMLHttpRequest', MockXHR as any);

describe('uploadFile', () => {
  it('resolves with JSON response and reports progress', async () => {
    const file = new File(['hello'], 'hello.txt', { type: 'text/plain' });
    const progressCalls: number[] = [];
    const p = uploadFile('server-1', file, 'web', (pct) => progressCalls.push(pct));
    const res = await p;
    expect(res).toHaveProperty('url');
    expect(res.url).toContain('http://');
  });
});

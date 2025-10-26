import '@testing-library/jest-dom';
import { beforeAll, afterAll, beforeEach, vi } from 'vitest';

let originalFetch: typeof fetch | undefined;

beforeAll(() => {
  originalFetch = (globalThis as any).fetch;
  (globalThis as any).fetch = vi.fn().mockImplementation(async (input: RequestInfo | URL, init?: RequestInit) => {
    const url = typeof input === 'string' ? input : String((input as URL).toString());

    // Helper to build a Response with JSON
    const jsonResponse = (obj: unknown) => new Response(JSON.stringify(obj), {
      status: 200,
      headers: { 'Content-Type': 'application/json' },
    });

    // Mail messages endpoint -> return { emails: [] }
    if (url.includes('/protocols/mail/') && url.endsWith('/messages')) {
      return jsonResponse({ emails: [] });
    }

    // Filetree endpoints -> return empty entries
    if (url.includes('/filetree')) {
      return jsonResponse({ entries: [], truncated: false });
    }

    // Web logs endpoint -> return empty logs array
    if (url.includes('/logs')) {
      return jsonResponse({ logs: [] });
    }

    // Upload endpoints or generic protocol endpoints -> return minimal container data
    if (url.includes('/protocols/') || url.includes('/api/v1/containers')) {
      return jsonResponse({ id: 'test-server', name: 'test-server', status: 'stopped' });
    }

    // Default: return empty JSON object
    return jsonResponse({});
  });
});

beforeEach(() => {
  const f = (globalThis as any).fetch;
  if (f && typeof f.mockClear === 'function') f.mockClear();
});

afterAll(() => {
  (globalThis as any).fetch = originalFetch;
});

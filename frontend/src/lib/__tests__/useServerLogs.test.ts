import { vi, describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useServerLogs } from '../useServerLogs';

vi.mock('../api', () => ({
  API_URL: 'http://localhost',
  request: vi.fn(),
  websocketConnect: vi.fn(),
  uploadFile: vi.fn(),
  fetchFileTree: vi.fn(),
  fetchWebFileTree: vi.fn(),
  fetchMailMessages: vi.fn(),
  fetchServerLogs: vi.fn(),
}));

import { fetchServerLogs } from '../api';

describe('useServerLogs', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('fetches mail logs on mount', async () => {
    const now = new Date().toISOString();
    const mockLogs = {
      lines: [
        { ts: now, line: 'SMTP server started on port 1025' },
        { ts: now, line: 'Received mail from sender@test.com' },
      ],
      truncated: false,
      container_running: true,
    };

    (fetchServerLogs as any).mockResolvedValue(mockLogs);

    const { result } = renderHook(() => useServerLogs('mail-server-1', 'mail', 500));

    expect(result.current.loading).toBe(true);
    expect(result.current.lines).toEqual([]);

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(fetchServerLogs).toHaveBeenCalledWith('mail-server-1', 'mail', 500);
    expect(result.current.lines).toHaveLength(2);
    expect(result.current.lines[0].line).toBe('SMTP server started on port 1025');
    expect(result.current.truncated).toBe(false);
    expect(result.current.error).toBeNull();
  });

  it('handles errors when fetching logs', async () => {
    (fetchServerLogs as any).mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useServerLogs('mail-server-1', 'mail'));

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.error).toBe('Network error');
    expect(result.current.lines).toEqual([]);
  });

  it('sets truncated flag when logs are truncated', async () => {
    const now = new Date().toISOString();
    (fetchServerLogs as any).mockResolvedValue({
      lines: [{ ts: now, line: 'log line' }],
      truncated: true,
      container_running: true,
    });

    const { result } = renderHook(() => useServerLogs('mail-server-1', 'mail', 50));

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.truncated).toBe(true);
  });

  it('refetches logs when refreshSignal changes', async () => {
    const now = new Date().toISOString();
    (fetchServerLogs as any).mockResolvedValue({
      lines: [{ ts: now, line: 'initial log' }],
      truncated: false,
      container_running: true,
    });

    const { result, rerender } = renderHook(
      ({ signal }) => useServerLogs('mail-server-1', 'mail', 500, signal),
      { initialProps: { signal: 0 } }
    );

    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(fetchServerLogs).toHaveBeenCalledTimes(1);

    (fetchServerLogs as any).mockResolvedValue({
      lines: [{ ts: now, line: 'refreshed log' }],
      truncated: false,
      container_running: true,
    });

    rerender({ signal: 1 });

    await waitFor(() => expect(fetchServerLogs).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(result.current.lines[0].line).toBe('refreshed log'));
  });
});

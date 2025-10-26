import { vi, describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useMailLogs } from '../useMailLogs';

vi.mock('../api', () => ({
  API_URL: 'http://localhost',
  request: vi.fn(),
  websocketConnect: vi.fn(),
  uploadFile: vi.fn(),
  fetchFileTree: vi.fn(),
  fetchWebLogs: vi.fn(),
  fetchWebFileTree: vi.fn(),
  fetchMailMessages: vi.fn(),
  fetchMailLogs: vi.fn(),
}));

import { fetchMailLogs } from '../api';

describe('useMailLogs', () => {
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

    (fetchMailLogs as any).mockResolvedValue(mockLogs);

    const { result } = renderHook(() => useMailLogs('mail-server-1', 500));

    expect(result.current.loading).toBe(true);
    expect(result.current.lines).toEqual([]);

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(fetchMailLogs).toHaveBeenCalledWith('mail-server-1', 500);
    expect(result.current.lines).toHaveLength(2);
    expect(result.current.lines[0].line).toBe('SMTP server started on port 1025');
    expect(result.current.truncated).toBe(false);
    expect(result.current.error).toBeNull();
  });

  it('handles errors when fetching logs', async () => {
    (fetchMailLogs as any).mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useMailLogs('mail-server-1'));

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.error).toBe('Network error');
    expect(result.current.lines).toEqual([]);
  });

  it('sets truncated flag when logs are truncated', async () => {
    const now = new Date().toISOString();
    (fetchMailLogs as any).mockResolvedValue({
      lines: [{ ts: now, line: 'log line' }],
      truncated: true,
      container_running: true,
    });

    const { result } = renderHook(() => useMailLogs('mail-server-1', 50));

    await waitFor(() => expect(result.current.loading).toBe(false));

    expect(result.current.truncated).toBe(true);
  });

  it('refetches logs when refreshSignal changes', async () => {
    const now = new Date().toISOString();
    (fetchMailLogs as any).mockResolvedValue({
      lines: [{ ts: now, line: 'initial log' }],
      truncated: false,
      container_running: true,
    });

    const { result, rerender } = renderHook(
      ({ signal }) => useMailLogs('mail-server-1', 500, signal),
      { initialProps: { signal: 0 } }
    );

    await waitFor(() => expect(result.current.loading).toBe(false));
    expect(fetchMailLogs).toHaveBeenCalledTimes(1);

    (fetchMailLogs as any).mockResolvedValue({
      lines: [{ ts: now, line: 'refreshed log' }],
      truncated: false,
      container_running: true,
    });

    rerender({ signal: 1 });

    await waitFor(() => expect(fetchMailLogs).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(result.current.lines[0].line).toBe('refreshed log'));
  });
});

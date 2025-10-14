import '@testing-library/jest-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';

vi.mock('@/lib/api', () => ({
  fetchWebLogs: vi.fn(),
}));

import { fetchWebLogs } from '@/lib/api';
import { LogsPanel } from '../LogsPanel';

describe('LogsPanel', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('renders lines from API', async () => {
    const now = new Date().toISOString();
    (fetchWebLogs as any).mockResolvedValue({ lines: [ { ts: now, line: '127.0.0.1 - - "GET /index.html HTTP/1.1" 200 123' } ], truncated: false, container_running: true });

    render(<LogsPanel serverId="server-1" />);

    await waitFor(() => expect(screen.getByText(/127.0.0.1/)).toBeInTheDocument());
  });

  it('shows truncated indicator', async () => {
    const now = new Date().toISOString();
    (fetchWebLogs as any).mockResolvedValue({ lines: [ { ts: now, line: 'line1' } ], truncated: true, container_running: true });

    render(<LogsPanel serverId="server-1" />);

    await waitFor(() => expect(screen.getByText(/truncated/i)).toBeInTheDocument());
  });
});

import '@testing-library/jest-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MailTab } from '../MailTab';

vi.mock('@/lib/api', () => ({
  default: vi.fn(),
  request: vi.fn(),
  fetchServerLogs: vi.fn(),
  fetchMailMessages: vi.fn(),
}));

import { fetchServerLogs, fetchMailMessages } from '@/lib/api';

describe('MailTab', () => {
  beforeEach(() => {
    vi.resetAllMocks();
    (fetchMailMessages as any).mockResolvedValue([]);
    (fetchServerLogs as any).mockResolvedValue({
      lines: [],
      truncated: false,
      container_running: true,
    });
  });

  it('renders MailTab with container information', () => {
    render(<MailTab id="test" reloadTabs={() => {}} />);
    expect(screen.getByText(/Container Information/i)).toBeInTheDocument();
  });

  it('renders Server Logs accordion', () => {
    render(<MailTab id="test" reloadTabs={() => {}} />);
    expect(screen.getByText(/Server Logs/i)).toBeInTheDocument();
  });

  it('displays mail logs when backend returns logs', async () => {
    const user = userEvent.setup();
    const now = new Date().toISOString();
    (fetchServerLogs as any).mockResolvedValue({
      lines: [
        { ts: now, line: 'SMTP server started on port 1025' },
        { ts: now, line: 'Received mail from sender@test.com' },
      ],
      truncated: false,
      container_running: true,
    });

    render(<MailTab id="mail-server-1" reloadTabs={() => {}} />);

    // Click to open the Server Logs accordion
    const serverLogsButton = screen.getByRole('button', { name: /Server Logs/i });
    await user.click(serverLogsButton);

    await waitFor(() => {
      expect(screen.getByText(/SMTP server started on port 1025/)).toBeInTheDocument();
    });

    expect(screen.getByText(/Received mail from sender@test.com/)).toBeInTheDocument();
  });

  it('displays E-Mails accordion', () => {
    render(<MailTab id="test" reloadTabs={() => {}} />);
    expect(screen.getByText(/E-Mails/i)).toBeInTheDocument();
  });
});


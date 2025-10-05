import '@testing-library/jest-dom';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

vi.mock('@/lib/api', () => ({
  fetchWebFileTree: vi.fn(),
}));

import { fetchWebFileTree } from '@/lib/api';
import FileTreeView from '../FileTreeView';

describe('FileTreeView', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('renders root entries and expands directory to load children', async () => {
    const now = new Date().toISOString();
    // root contains one dir and one file
    (fetchWebFileTree as any).mockImplementation(async (_serverId: string, path: string | null) => {
      if (!path) {
        return { entries: [ { name: 'assets', path: 'assets', type: 'dir', size: 4096, modifiedAt: now }, { name: 'index.html', path: 'index.html', type: 'file', size: 128, modifiedAt: now } ], truncated: false } as unknown as { entries: import('@/types/FileTree').FileTreeEntry[]; truncated: boolean };
      }
      if (path === 'assets') {
        return { entries: [ { name: 'app.js', path: 'assets/app.js', type: 'file', size: 1024, modifiedAt: now } ], truncated: false } as unknown as { entries: import('@/types/FileTree').FileTreeEntry[]; truncated: boolean };
      }
      return { entries: [], truncated: false };
    });

    render(<FileTreeView serverId="server-1" baseUrl="http://localhost:8080" />);

    // root items should appear
    await waitFor(() => expect(screen.getByText('assets')).toBeInTheDocument());
    expect(screen.getByText('index.html')).toBeInTheDocument();

    // expand the assets directory by clicking its expander (Chevron)
    const expanders = screen.queryAllByRole('button');
    if (expanders.length > 0) {
      // click first expander if it exists
      await userEvent.click(expanders[0]);
    } else {
      // fallback: click the assets label
      await userEvent.click(screen.getByText('assets'));
    }

    // child should be loaded and visible
    await waitFor(() => expect(screen.getByText('app.js')).toBeInTheDocument());
  });
});

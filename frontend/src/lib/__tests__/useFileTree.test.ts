import { vi, describe, it, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useFileTree } from '../useFileTree';

vi.mock('../api', () => ({
  API_URL: 'http://localhost',
  request: vi.fn(),
  websocketConnect: vi.fn(),
  uploadFile: vi.fn(),
  fetchFileTree: vi.fn(),
  fetchWebLogs: vi.fn(),
  fetchWebFileTree: vi.fn(),
}));

import { fetchWebFileTree } from '../api';

describe('useFileTree', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('fetches root children and caches result', async () => {
    
    (fetchWebFileTree as any).mockResolvedValueOnce({ entries: [{ name: 'index.html', path: 'index.html', type: 'file', size: 10, modifiedAt: new Date().toISOString() }], truncated: false } as unknown as { entries: import('@/types/FileTree').FileTreeEntry[]; truncated: boolean });

    const { result } = renderHook(() => useFileTree('server-1'));

    const data = await act(async () => {
      return await result.current.getChildren(null);
    });

    expect((fetchWebFileTree as any)).toHaveBeenCalledWith('server-1', null);
    expect(data.entries).toHaveLength(1);

    // subsequent call should return cached value and not call fetch again
    (fetchWebFileTree as any).mockClear();
    const cached = await act(async () => result.current.getChildren(null));
    expect((fetchWebFileTree as any)).not.toHaveBeenCalled();
    expect(cached.entries).toHaveLength(1);
  });
});

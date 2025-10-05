import { useRef, useState } from 'react';
import type { FileTreeEntry } from '../types/FileTree';
import { fetchWebFileTree } from './api';

export function useFileTree(serverId: string) {
    const cache = useRef(new Map<string, { entries: FileTreeEntry[]; truncated: boolean }>());
    const [loadingPaths, setLoadingPaths] = useState<Record<string, boolean>>({});

    function setLoading(path: string, v: boolean) {
        setLoadingPaths((s) => ({ ...s, [path]: v }));
    }

    async function getChildren(path: string | null = null): Promise<{ entries: FileTreeEntry[]; truncated: boolean }> {
        const key = path || '';
        const existing = cache.current.get(key);
        if (existing) {
            return existing;
        }
        setLoading(key, true);
        try {
            const res = await fetchWebFileTree(serverId, path);
            cache.current.set(key, { entries: res.entries, truncated: res.truncated });
            return { entries: res.entries, truncated: res.truncated };
        } finally {
            setLoading(key, false);
        }
    }

    async function refresh(path: string | null = null): Promise<{ entries: FileTreeEntry[]; truncated: boolean }> {
        const key = path || '';
        setLoading(key, true);
        try {
            const res = await fetchWebFileTree(serverId, path);
            cache.current.set(key, { entries: res.entries, truncated: res.truncated });
            return { entries: res.entries, truncated: res.truncated };
        } finally {
            setLoading(key, false);
        }
    }

    function getCached(path: string | null = null) {
        return cache.current.get(path || '') ?? null;
    }

    return { getChildren, refresh, getCached, loadingPaths } as const;
}

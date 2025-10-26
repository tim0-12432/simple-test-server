import { useEffect, useState, useRef } from 'react';
import type { LogLine } from '@/types/Log';
import { fetchServerLogs } from './api';

export function useServerLogs(serverId: string, type: 'web' | 'mail', tail: number = 500, refreshSignal?: number) {
    const [lines, setLines] = useState<LogLine[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [truncated, setTruncated] = useState(false);
    const latestSignal = useRef<number | undefined>(refreshSignal);

    useEffect(() => {
        latestSignal.current = refreshSignal;
        let cancelled = false;
        setLoading(true);
        setError(null);
        (async () => {
            try {
                const res = await fetchServerLogs(serverId, type, tail);
                if (cancelled) return;
                setLines(res.lines);
                setTruncated(res.truncated);
            } catch (e: any) {
                if (cancelled) return;
                setError(e?.message ?? 'Failed to load logs');
            } finally {
                if (!cancelled) setLoading(false);
            }
        })();
        return () => { cancelled = true; };
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [serverId, tail, refreshSignal]);

    return { lines, loading, error, truncated };
}

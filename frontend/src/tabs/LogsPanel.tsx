import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/progress';
import { useWebLogs } from '@/lib/useWebLogs';
import type { LogLine } from '@/types/Log';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

export function LogsPanel({ serverId, refreshSignal }: { serverId: string; refreshSignal?: number }) {
    const [tail, setTail] = useState<number>(500);
    const [localRefresh, setLocalRefresh] = useState<number>(0);
    const { lines, loading, error, truncated } = useWebLogs(serverId, tail, (refreshSignal ?? 0) + localRefresh);

    const handleLocalRefresh = (e?: React.MouseEvent) => {
        if (e) e.stopPropagation();
        setLocalRefresh((s) => s + 1);
    };

    return (
        <div className="w-full">
            <div className="flex items-center gap-2 mb-2">
                <label className="text-sm">Lines:</label>
                <Select onValueChange={(value) => setTail(Number(value))} value={tail.toString()}>
                    <SelectTrigger>
                        <SelectValue placeholder="Lines" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value={'50'}>50</SelectItem>
                        <SelectItem value={'200'}>200</SelectItem>
                        <SelectItem value={'500'}>500</SelectItem>
                        <SelectItem value={'2000'}>2000</SelectItem>
                    </SelectContent>
                </Select>
                <Button variant="ghost" onClick={handleLocalRefresh}>Refresh</Button>
                {truncated && <span className="text-xs text-muted-foreground">(truncated)</span>}
            </div>
            <Progress active={loading} className="w-full mb-2 h-2" />
            {error ? <div className="text-destructive">{error}</div> : (
                <pre className="bg-surface/40 p-2 rounded text-sm overflow-auto max-h-96">{lines.map((l: LogLine) => `${l.ts} ${l.line}`).join('\n')}</pre>
            )}
        </div>
    );
}

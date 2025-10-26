import { useEffect, useRef, useState } from 'react';
import { websocketConnect } from './api';
import type { OtelTelemetry } from '@/types/OtelData';

export function useOtelTelemetry(serverId: string | null, maxMessages = 1000) {
    const [messages, setMessages] = useState<OtelTelemetry[]>([]);
    const [connected, setConnected] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const wsRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        setError(null);
        if (wsRef.current) {
            // eslint-disable-next-line @typescript-eslint/no-unused-vars
            try { wsRef.current.close(); } catch (e) { /* ignore */ }
            wsRef.current = null;
            setConnected(false);
        }

        if (serverId) {
            const ws = websocketConnect<string>(`/protocols/otel/${serverId}/telemetry`, messageHandler, errorHandler);
            wsRef.current = ws;
            ws.onopen = () => setConnected(true);
            ws.onclose = () => setConnected(false);
            ws.onerror = () => setConnected(false);
        }

        return () => {
            if (wsRef.current) {
                // eslint-disable-next-line @typescript-eslint/no-unused-vars
                try { wsRef.current.close(); } catch (e) { /* ignore */ }
                wsRef.current = null;
            }
        };
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [serverId]);

    function messageHandler(msg: OtelTelemetry) {
        console.log("Received OTEL telemetry message:", msg);
        setMessages(prev => {
            const next = [...prev, msg];
            if (next.length > maxMessages) {
                return next.slice(next.length - maxMessages);
            }
            return next;
        });
    }

    function errorHandler(err: Event) {
        setError(`WebSocket error: ${err instanceof Error ? err.message : 'Unknown error'}`);
    }

    function clear() {
        setMessages([]);
    }

    return { messages, connected, error, clear };
}

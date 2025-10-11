import React from 'react';
import type MqttData from '@/types/MqttData';

type MessageLogProps = {
    messages: MqttData[];
};

export default function MessageLog({ messages }: MessageLogProps) {
    if (!messages || messages.length === 0) {
        return <div className="text-sm text-muted-foreground">No messages received yet</div>;
    }

    return (
        <div className="p-2">
            <ul>
                {messages.map((m, idx) => {
                    let formattedPayload = m.payload;
                    try {
                        formattedPayload = JSON.stringify(JSON.parse(m.payload), null, 2);
                    } catch {
                        // If payload is not valid JSON, keep original
                    }
                    return (
                        <li key={idx} className="py-2 pb-3 border-b last:border-0 border-border">
                            <div className="font-medium">{m.topic}</div>
                            <pre className="text-sm px-4 py-2">{formattedPayload}</pre>
                            <div className="text-xs text-muted-foreground">{m.timestamp ?? new Date().toISOString()}</div>
                        </li>
                    );
                })}
            </ul>
        </div>
    );
}

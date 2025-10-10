import React from 'react';
import type MqttData from '@/types/MqttData';

type MessageLogProps = {
    messages: MqttData[];
};

export default function MessageLog({ messages }: MessageLogProps) {
    if (!messages || messages.length === 0) {
        return <div className="p-2">No messages received yet</div>;
    }

    return (
        <div className="p-2">
            <ul>
                {messages.map((m, idx) => (
                    <li key={idx} className="py-1">
                        <div className="font-medium">{m.topic}</div>
                        <div className="text-sm text-muted">{m.payload}</div>
                        <div className="text-xs text-muted-foreground">{m.timestamp ?? new Date().toISOString()}</div>
                    </li>
                ))}
            </ul>
        </div>
    );
}

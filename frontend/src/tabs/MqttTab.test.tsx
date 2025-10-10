import React from 'react';
import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';

vi.mock('@/lib/api', () => ({
    websocketConnect: vi.fn(() => ({
        onopen: null,
        onclose: null,
        onerror: null,
        close: vi.fn(),
    })),
}));

import MqttTab from './MqttTab';

const baseProps = {
    type: 'MQTT',
    id: 'mqtt-server-1',
    reloadTabs: () => {}
};

test('renders topic tree and message log accordions', () => {
    render(<MqttTab {...(baseProps as any)} />);
    expect(screen.getByText('Topic Tree')).toBeInTheDocument();
    expect(screen.getByText('Message Log')).toBeInTheDocument();
});

test('shows no messages text when none', () => {
    render(<MqttTab {...(baseProps as any)} />);
    expect(screen.getByText(/No messages received yet/i)).toBeInTheDocument();
});

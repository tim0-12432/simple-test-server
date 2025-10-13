import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
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

test('renders topic tree and message log accordions', async () => {
    render(<MqttTab {...(baseProps as any)} />);
    expect(await screen.findByText('Topic Tree')).toBeInTheDocument();
    expect(await screen.findByText('Message Log')).toBeInTheDocument();
});

test('shows no messages text when none', async () => {
    const user = userEvent.setup();
    render(<MqttTab {...(baseProps as any)} />);
    // open the Message Log accordion before asserting its content
    await user.click(await screen.findByText('Message Log'));
    expect(await screen.findByText(/No messages received yet/i)).toBeInTheDocument();
});



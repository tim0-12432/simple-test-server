/* eslint-disable @typescript-eslint/no-explicit-any */
import React from 'react';
// React import removed; JSX runtime handles React
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi, describe, it, expect } from 'vitest';

// Mock UI modules and API BEFORE importing the tested module so imports don't resolve real aliased files
vi.mock('../lib/api', () => ({
  uploadFile: vi.fn(( _serverId: string, _file: File, _serverType: string, onProgress?: (pct: number) => void) => {
    return new Promise((resolve) => {
      if (onProgress) {
        onProgress(20);
        onProgress(60);
        onProgress(100);
      }
      resolve({ url: 'http://localhost:8000/test' });
    });
  })
}));

// Simple component stubs that render children so tests can assert visible text
vi.mock('../components/ui/alert', () => ({
  Alert: (props: any) => React.createElement('div', {}, props.children),
  AlertDescription: (props: any) => React.createElement('div', {}, props.children),
  AlertTitle: (props: any) => React.createElement('div', {}, props.children),
}));
vi.mock('../components/ui/accordion', () => ({
  Accordion: (props: any) => React.createElement('div', {}, props.children),
}));
vi.mock('../components/tab-accordion', () => ({
  __esModule: true,
  default: (props: any) => React.createElement('div', {}, props.children),
}));
vi.mock('../components/server-information', () => ({
  __esModule: true,
  default: (props: any) => React.createElement('div', {}, props.children),
}));
vi.mock('../components/ui/button', () => ({
  Button: (props: any) => React.createElement('button', { type: props.type ?? 'button', onClick: props.onClick, disabled: props.disabled }, props.children),
}));
vi.mock('../components/ui/kibo-ui/dropzone', () => ({
  Dropzone: (props: any) => {
    const Simulate = () => React.createElement('button', { onClick: () => props.onDrop && props.onDrop([new File(['content'], 'test.txt', { type: 'text/plain' })]) }, 'Simulate Drop');
    return React.createElement('div', {}, React.createElement(Simulate), props.children);
  },
  DropzoneContent: (props: any) => React.createElement('div', {}, props.children),
  DropzoneEmptyState: (props: any) => React.createElement('div', {}, props.children),
}));

import { WebTab } from './WebTab';

// Minimal GeneralTabInformation mock
const baseProps = {
  type: 'WEB' as const,
  id: 'server-1',
  reloadTabs: () => {}
};

describe('WebTab upload flow', () => {
  it('uploads a dropped file and displays uploaded URL', async () => {
    render(React.createElement(WebTab, baseProps));

    // Simulate dropping a file using the Dropzone mock
    const simulateDropButton = screen.getByText('Simulate Drop');
    fireEvent.click(simulateDropButton);

    const uploadButton = screen.getByRole('button', { name: /upload/i });
    expect((uploadButton as HTMLButtonElement).disabled).toBe(false);
    fireEvent.click(uploadButton);

    // Should show the uploaded resource link
    await waitFor(() => expect(screen.getByText(/open uploaded resource/i)).toBeTruthy());
  });
});

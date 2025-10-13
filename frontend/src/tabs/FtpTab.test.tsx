/* eslint-disable @typescript-eslint/no-explicit-any */
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { vi, describe, it, expect } from 'vitest';

// Mock API before importing module
vi.mock('../lib/api', () => ({
  uploadFile: vi.fn((_serverId: string, _file: File, _serverType: string, onProgress?: (pct: number) => void) => {
    return new Promise((resolve) => {
      if (onProgress) {
        onProgress(20);
        onProgress(60);
        onProgress(100);
      }
      resolve({ url: 'ftp://localhost/test' });
    });
  })
}));

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

import FtpTab from './FtpTab';

const baseProps = {
  type: 'FTP' as const,
  id: 'server-ftp',
  reloadTabs: () => {}
};

describe('FtpTab upload flow', () => {
  it('uploads a dropped file and displays uploaded URL', async () => {
    render(React.createElement(FtpTab, baseProps));

    const simulateDropButton = screen.getByText('Simulate Drop');
    fireEvent.click(simulateDropButton);

    const uploadButton = await screen.findByRole('button', { name: /upload/i });
    await waitFor(() => expect((uploadButton as HTMLButtonElement).disabled).toBe(false));
    fireEvent.click(uploadButton);

    await waitFor(() => expect(screen.getByText(/open uploaded resource/i)).toBeTruthy());
  });
});

import { render, screen } from '@testing-library/react';
import { SmbTab } from '../SmbTab';

const noop = () => {};

test('renders SmbTab without crashing', () => {
  render(<SmbTab id="test" type={'SMB'} reloadTabs={noop} /> as any);
  expect(screen.getByText(/Folder Tree/i)).toBeInTheDocument();
});

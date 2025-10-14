import { render, screen } from '@testing-library/react';
import { MailTab } from '../MailTab';

test('renders MailTab with container information', () => {
  render(<MailTab id="test" reloadTabs={() => {}} />);
  expect(screen.getByText(/Container Information/i)).toBeInTheDocument();
});

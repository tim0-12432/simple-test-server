import { render, screen } from '@testing-library/react';
import MailTab from '../MailTab';

test('renders MailTab with info alert', () => {
  render(<MailTab id="test" reloadTabs={() => {}} />);
  expect(screen.getByText(/Mail tab is not yet implemented/i)).toBeInTheDocument();
});

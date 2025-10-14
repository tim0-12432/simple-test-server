import { render, screen } from '@testing-library/react';
import TabFactory from '../TabFactory';

test('TabFactory returns MailTab for MAIL type', () => {
  const component = TabFactory('MAIL', { id: 'abc', type: 'MAIL', reloadTabs: () => {} });
  // Render the returned component
  render(component as any);
  expect(screen.getByText(/Mail tab is not yet implemented/i)).toBeInTheDocument();
});

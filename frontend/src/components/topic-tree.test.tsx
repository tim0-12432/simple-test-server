import { render, screen } from '@testing-library/react'
import TopicTree from './topic-tree'

const messages = [
  { topic: 'sensors/temp', payload: '23.5' },
  { topic: 'sensors/humidity', payload: '65' },
]

test('renders topic tree with messages', () => {
  render(<TopicTree messages={messages as any} />)
  expect(screen.getByText('sensors')).toBeInTheDocument()
  expect(screen.getByText('temp')).toBeInTheDocument()
  expect(screen.getByText('23.5')).toBeInTheDocument()
})

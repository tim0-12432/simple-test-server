import { render, screen } from '@testing-library/react'
import TopicTree from './topic-tree'

const messages = [
  { topic: 'sensors/temp', payload: '23.5' },
  { topic: 'sensors/humidity', payload: '65' },
]

test('renders topic tree with messages', async () => {
  render(<TopicTree messages={messages as any} />)
  expect(await screen.findByText('sensors')).toBeInTheDocument()
  expect(await screen.findByText('temp')).toBeInTheDocument()
  expect(await screen.findByText('23.5')).toBeInTheDocument()
})


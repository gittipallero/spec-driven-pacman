import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { HUD } from './HUD';

describe('HUD', () => {
  it('renders HUD container', () => {
    render(<HUD score={0} lives={3} level={1} />);
    
    expect(screen.getByTestId('hud')).toBeInTheDocument();
  });

  it('displays score with leading zeros', () => {
    render(<HUD score={1234} lives={3} level={1} />);
    
    expect(screen.getByTestId('score')).toHaveTextContent('001234');
  });

  it('displays score of zero correctly', () => {
    render(<HUD score={0} lives={3} level={1} />);
    
    expect(screen.getByTestId('score')).toHaveTextContent('000000');
  });

  it('displays large scores correctly', () => {
    render(<HUD score={999999} lives={3} level={1} />);
    
    expect(screen.getByTestId('score')).toHaveTextContent('999999');
  });

  it('displays correct number of lives', () => {
    render(<HUD score={0} lives={3} level={1} />);
    
    const lives = screen.getByTestId('lives');
    expect(lives.children).toHaveLength(3);
  });

  it('updates lives count', () => {
    const { rerender } = render(<HUD score={0} lives={3} level={1} />);
    expect(screen.getByTestId('lives').children).toHaveLength(3);

    rerender(<HUD score={0} lives={2} level={1} />);
    expect(screen.getByTestId('lives').children).toHaveLength(2);

    rerender(<HUD score={0} lives={1} level={1} />);
    expect(screen.getByTestId('lives').children).toHaveLength(1);
  });

  it('displays zero lives', () => {
    render(<HUD score={0} lives={0} level={1} />);
    
    expect(screen.getByTestId('lives').children).toHaveLength(0);
  });

  it('displays level number', () => {
    render(<HUD score={0} lives={3} level={5} />);
    
    expect(screen.getByTestId('level')).toHaveTextContent('5');
  });

  it('updates when props change', () => {
    const { rerender } = render(<HUD score={100} lives={3} level={1} />);
    
    expect(screen.getByTestId('score')).toHaveTextContent('000100');
    
    rerender(<HUD score={200} lives={2} level={2} />);
    
    expect(screen.getByTestId('score')).toHaveTextContent('000200');
    expect(screen.getByTestId('lives').children).toHaveLength(2);
    expect(screen.getByTestId('level')).toHaveTextContent('2');
  });
});


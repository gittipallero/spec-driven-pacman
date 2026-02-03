import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { GameOver } from './GameOver';

describe('GameOver', () => {
  const mockOnRestart = vi.fn();
  const mockOnMainMenu = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('does not render when status is playing', () => {
    render(
      <GameOver
        status="playing"
        score={1000}
        level={1}
        onRestart={mockOnRestart}
        onMainMenu={mockOnMainMenu}
      />
    );

    expect(screen.queryByTestId('game-over-overlay')).not.toBeInTheDocument();
  });

  it('does not render when status is paused', () => {
    render(
      <GameOver
        status="paused"
        score={1000}
        level={1}
        onRestart={mockOnRestart}
        onMainMenu={mockOnMainMenu}
      />
    );

    expect(screen.queryByTestId('game-over-overlay')).not.toBeInTheDocument();
  });

  describe('game over state', () => {
    it('renders game over overlay', () => {
      render(
        <GameOver
          status="gameOver"
          score={5000}
          level={3}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      expect(screen.getByTestId('game-over-overlay')).toBeInTheDocument();
    });

    it('shows GAME OVER title', () => {
      render(
        <GameOver
          status="gameOver"
          score={5000}
          level={3}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      expect(screen.getByTestId('game-over-title')).toHaveTextContent('GAME OVER');
    });

    it('displays final score', () => {
      render(
        <GameOver
          status="gameOver"
          score={12340}
          level={3}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      expect(screen.getByTestId('final-score')).toHaveTextContent('12,340');
    });

    it('shows PLAY AGAIN button for game over', () => {
      render(
        <GameOver
          status="gameOver"
          score={5000}
          level={3}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      expect(screen.getByTestId('restart-button')).toHaveTextContent('PLAY AGAIN');
    });

    it('calls onRestart when restart button clicked', () => {
      render(
        <GameOver
          status="gameOver"
          score={5000}
          level={3}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      fireEvent.click(screen.getByTestId('restart-button'));
      expect(mockOnRestart).toHaveBeenCalledTimes(1);
    });
  });

  describe('level complete state', () => {
    it('renders level complete overlay', () => {
      render(
        <GameOver
          status="levelComplete"
          score={8000}
          level={4}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      expect(screen.getByTestId('game-over-overlay')).toBeInTheDocument();
    });

    it('shows LEVEL COMPLETE title', () => {
      render(
        <GameOver
          status="levelComplete"
          score={8000}
          level={4}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      expect(screen.getByTestId('game-over-title')).toHaveTextContent('LEVEL COMPLETE!');
    });

    it('shows next level number', () => {
      render(
        <GameOver
          status="levelComplete"
          score={8000}
          level={4}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      expect(screen.getByTestId('next-level')).toHaveTextContent('4');
    });

    it('shows CONTINUE button for level complete', () => {
      render(
        <GameOver
          status="levelComplete"
          score={8000}
          level={4}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      expect(screen.getByTestId('restart-button')).toHaveTextContent('CONTINUE');
    });
  });

  describe('main menu button', () => {
    it('shows main menu button', () => {
      render(
        <GameOver
          status="gameOver"
          score={5000}
          level={3}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      expect(screen.getByTestId('main-menu-button')).toHaveTextContent('MAIN MENU');
    });

    it('calls onMainMenu when main menu button clicked', () => {
      render(
        <GameOver
          status="gameOver"
          score={5000}
          level={3}
          onRestart={mockOnRestart}
          onMainMenu={mockOnMainMenu}
        />
      );

      fireEvent.click(screen.getByTestId('main-menu-button'));
      expect(mockOnMainMenu).toHaveBeenCalledTimes(1);
    });
  });
});


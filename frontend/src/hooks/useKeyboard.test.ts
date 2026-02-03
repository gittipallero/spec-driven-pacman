import { renderHook } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useKeyboard } from './useKeyboard';
import type { Direction } from '../stores/gameStore';

describe('useKeyboard', () => {
  let mockOnDirection: ReturnType<typeof vi.fn>;
  let mockOnPause: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    mockOnDirection = vi.fn();
    mockOnPause = vi.fn();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  const fireKeyDown = (key: string, target?: EventTarget | null) => {
    const event = new KeyboardEvent('keydown', { key, bubbles: true });
    if (target) {
      Object.defineProperty(event, 'target', { value: target });
    }
    window.dispatchEvent(event);
  };

  describe('direction keys', () => {
    it('calls onDirection with "up" for ArrowUp', () => {
      renderHook(() =>
        useKeyboard({ onDirection: mockOnDirection, enabled: true })
      );

      fireKeyDown('ArrowUp');

      expect(mockOnDirection).toHaveBeenCalledWith('up');
    });

    it('calls onDirection with "down" for ArrowDown', () => {
      renderHook(() =>
        useKeyboard({ onDirection: mockOnDirection, enabled: true })
      );

      fireKeyDown('ArrowDown');

      expect(mockOnDirection).toHaveBeenCalledWith('down');
    });

    it('calls onDirection with "left" for ArrowLeft', () => {
      renderHook(() =>
        useKeyboard({ onDirection: mockOnDirection, enabled: true })
      );

      fireKeyDown('ArrowLeft');

      expect(mockOnDirection).toHaveBeenCalledWith('left');
    });

    it('calls onDirection with "right" for ArrowRight', () => {
      renderHook(() =>
        useKeyboard({ onDirection: mockOnDirection, enabled: true })
      );

      fireKeyDown('ArrowRight');

      expect(mockOnDirection).toHaveBeenCalledWith('right');
    });

    it('supports WASD keys', () => {
      renderHook(() =>
        useKeyboard({ onDirection: mockOnDirection, enabled: true })
      );

      fireKeyDown('w');
      expect(mockOnDirection).toHaveBeenCalledWith('up');

      fireKeyDown('a');
      expect(mockOnDirection).toHaveBeenCalledWith('left');

      fireKeyDown('s');
      expect(mockOnDirection).toHaveBeenCalledWith('down');

      fireKeyDown('d');
      expect(mockOnDirection).toHaveBeenCalledWith('right');
    });
  });

  describe('Escape key', () => {
    it('calls onPause when Escape is pressed', () => {
      renderHook(() =>
        useKeyboard({ onPause: mockOnPause, enabled: true })
      );

      fireKeyDown('Escape');

      expect(mockOnPause).toHaveBeenCalledTimes(1);
    });

    it('does not call onDirection when Escape is pressed', () => {
      renderHook(() =>
        useKeyboard({
          onDirection: mockOnDirection,
          onPause: mockOnPause,
          enabled: true,
        })
      );

      fireKeyDown('Escape');

      expect(mockOnDirection).not.toHaveBeenCalled();
      expect(mockOnPause).toHaveBeenCalledTimes(1);
    });
  });

  describe('enabled state', () => {
    it('does not respond to keys when disabled', () => {
      renderHook(() =>
        useKeyboard({
          onDirection: mockOnDirection,
          onPause: mockOnPause,
          enabled: false,
        })
      );

      fireKeyDown('ArrowUp');
      fireKeyDown('Escape');

      expect(mockOnDirection).not.toHaveBeenCalled();
      expect(mockOnPause).not.toHaveBeenCalled();
    });

    it('responds to keys when enabled (default)', () => {
      renderHook(() =>
        useKeyboard({ onDirection: mockOnDirection })
      );

      fireKeyDown('ArrowUp');

      expect(mockOnDirection).toHaveBeenCalledWith('up');
    });
  });

  describe('cleanup', () => {
    it('removes event listener on unmount', () => {
      const { unmount } = renderHook(() =>
        useKeyboard({ onDirection: mockOnDirection, enabled: true })
      );

      unmount();

      fireKeyDown('ArrowUp');

      expect(mockOnDirection).not.toHaveBeenCalled();
    });
  });

  describe('input field filtering', () => {
    it('ignores keys when typing in input field', () => {
      renderHook(() =>
        useKeyboard({ onDirection: mockOnDirection, enabled: true })
      );

      const input = document.createElement('input');
      fireKeyDown('ArrowUp', input);

      expect(mockOnDirection).not.toHaveBeenCalled();
    });

    it('ignores keys when typing in textarea', () => {
      renderHook(() =>
        useKeyboard({ onDirection: mockOnDirection, enabled: true })
      );

      const textarea = document.createElement('textarea');
      fireKeyDown('ArrowUp', textarea);

      expect(mockOnDirection).not.toHaveBeenCalled();
    });
  });

  describe('unhandled keys', () => {
    it('ignores unrecognized keys', () => {
      renderHook(() =>
        useKeyboard({
          onDirection: mockOnDirection,
          onPause: mockOnPause,
          enabled: true,
        })
      );

      fireKeyDown('Enter');
      fireKeyDown('Space');
      fireKeyDown('x');

      expect(mockOnDirection).not.toHaveBeenCalled();
      expect(mockOnPause).not.toHaveBeenCalled();
    });
  });
});


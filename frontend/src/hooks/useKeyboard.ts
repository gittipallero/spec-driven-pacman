import { useEffect, useCallback } from 'react';
import type { Direction } from '../stores/gameStore';

interface UseKeyboardOptions {
  onDirection?: (direction: Direction) => void;
  onPause?: () => void;
  enabled?: boolean;
}

const KEY_TO_DIRECTION: Record<string, Direction> = {
  ArrowUp: 'up',
  ArrowDown: 'down',
  ArrowLeft: 'left',
  ArrowRight: 'right',
  w: 'up',
  s: 'down',
  a: 'left',
  d: 'right',
  W: 'up',
  S: 'down',
  A: 'left',
  D: 'right',
};

/**
 * Hook for handling keyboard input in the game
 * @param options.onDirection - Callback for direction keys (arrows or WASD)
 * @param options.onPause - Callback for Escape key (pause/resume)
 * @param options.enabled - Whether to listen for keyboard events
 */
export function useKeyboard({
  onDirection,
  onPause,
  enabled = true,
}: UseKeyboardOptions = {}): void {
  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      // Ignore if typing in an input field
      if (
        event.target instanceof HTMLInputElement ||
        event.target instanceof HTMLTextAreaElement
      ) {
        return;
      }

      // Handle Escape for pause
      if (event.key === 'Escape') {
        event.preventDefault();
        onPause?.();
        return;
      }

      // Handle direction keys
      const direction = KEY_TO_DIRECTION[event.key];
      if (direction) {
        event.preventDefault();
        onDirection?.(direction);
      }
    },
    [onDirection, onPause]
  );

  useEffect(() => {
    if (!enabled) {
      return;
    }

    window.addEventListener('keydown', handleKeyDown);

    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [enabled, handleKeyDown]);
}

export default useKeyboard;


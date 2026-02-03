import { memo, useMemo } from 'react';
import type { Direction, PacManState } from '../../stores/gameStore';
import styles from './PacMan.module.css';

interface PacManProps {
  x: number;
  y: number;
  direction: Direction;
  state: PacManState;
  cellSize?: number;
}

const DIRECTION_ROTATION: Record<Direction, number> = {
  right: 0,
  down: 90,
  left: 180,
  up: 270,
  none: 0,
};

function PacManComponent({
  x,
  y,
  direction,
  state,
  cellSize = 16,
}: PacManProps) {
  const style = useMemo(() => ({
    left: `${x * cellSize}px`,
    top: `${y * cellSize}px`,
    width: `${cellSize}px`,
    height: `${cellSize}px`,
    transform: `rotate(${DIRECTION_ROTATION[direction]}deg)`,
  }), [x, y, direction, cellSize]);

  const className = useMemo(() => {
    const classes = [styles.pacman];
    if (state === 'moving') classes.push(styles.moving);
    if (state === 'dying') classes.push(styles.dying);
    return classes.join(' ');
  }, [state]);

  return (
    <div
      className={className}
      style={style}
      data-testid="pacman"
      data-direction={direction}
      data-state={state}
    >
      <div className={styles.body}>
        <div className={styles.mouthTop} />
        <div className={styles.mouthBottom} />
      </div>
    </div>
  );
}

export const PacMan = memo(PacManComponent);
export default PacMan;


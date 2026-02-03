import { memo } from 'react';
import styles from './HUD.module.css';

interface HUDProps {
  score: number;
  lives: number;
  level: number;
}

function HUDComponent({ score, lives, level }: HUDProps) {
  return (
    <div className={styles.hud} data-testid="hud">
      <div className={styles.section}>
        <span className={styles.label}>SCORE</span>
        <span className={styles.value} data-testid="score">
          {score.toString().padStart(6, '0')}
        </span>
      </div>

      <div className={styles.section}>
        <span className={styles.label}>LEVEL</span>
        <span className={styles.value} data-testid="level">
          {level}
        </span>
      </div>

      <div className={styles.section}>
        <span className={styles.label}>LIVES</span>
        <div className={styles.lives} data-testid="lives">
          {Array.from({ length: lives }).map((_, i) => (
            <span key={i} className={styles.lifeIcon} aria-label="life">
              ●
            </span>
          ))}
        </div>
      </div>
    </div>
  );
}

export const HUD = memo(HUDComponent);
export default HUD;


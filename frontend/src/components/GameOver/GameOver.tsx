import { memo } from 'react';
import type { GameStatus } from '../../stores/gameStore';
import styles from './GameOver.module.css';

interface GameOverProps {
  status: GameStatus;
  score: number;
  level: number;
  onRestart: () => void;
  onMainMenu: () => void;
}

function GameOverComponent({
  status,
  score,
  level,
  onRestart,
  onMainMenu,
}: GameOverProps) {
  const isGameOver = status === 'gameOver';
  const isLevelComplete = status === 'levelComplete';

  if (!isGameOver && !isLevelComplete) {
    return null;
  }

  return (
    <div className={styles.overlay} data-testid="game-over-overlay">
      <div className={styles.content}>
        <h1 
          className={`${styles.title} ${isGameOver ? styles.gameOver : styles.levelComplete}`}
          data-testid="game-over-title"
        >
          {isGameOver ? 'GAME OVER' : 'LEVEL COMPLETE!'}
        </h1>

        <div className={styles.stats}>
          <div className={styles.stat}>
            <span className={styles.statLabel}>FINAL SCORE</span>
            <span className={styles.statValue} data-testid="final-score">
              {score.toLocaleString()}
            </span>
          </div>
          
          {!isGameOver && (
            <div className={styles.stat}>
              <span className={styles.statLabel}>NEXT LEVEL</span>
              <span className={styles.statValue} data-testid="next-level">
                {level}
              </span>
            </div>
          )}
        </div>

        <div className={styles.buttons}>
          <button
            className={`${styles.button} ${styles.primary}`}
            onClick={onRestart}
            data-testid="restart-button"
          >
            {isGameOver ? 'PLAY AGAIN' : 'CONTINUE'}
          </button>
          
          <button
            className={`${styles.button} ${styles.secondary}`}
            onClick={onMainMenu}
            data-testid="main-menu-button"
          >
            MAIN MENU
          </button>
        </div>
      </div>
    </div>
  );
}

export const GameOver = memo(GameOverComponent);
export default GameOver;


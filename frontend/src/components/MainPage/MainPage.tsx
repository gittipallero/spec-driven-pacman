import { CRTEffect } from "../CRTEffect";
import { ArcadeTitle } from "../ArcadeTitle";
import { StartButton } from "../StartButton";
import styles from "./MainPage.module.css";

interface MainPageProps {
  /** Callback when start game is clicked */
  onStartGame?: () => void;
}

export function MainPage({ onStartGame }: MainPageProps) {
  const handleStartClick = () => {
    if (onStartGame) {
      onStartGame();
    } else {
      // Default behavior: log to console (placeholder for future navigation)
      console.log("Start game clicked - navigation will be added later");
    }
  };

  return (
    <CRTEffect>
      <main className={styles.mainPage}>
        <div className={styles.content}>
          <ArcadeTitle />
          <div className={styles.buttonContainer}>
            <StartButton onClick={handleStartClick} />
          </div>
        </div>
      </main>
    </CRTEffect>
  );
}


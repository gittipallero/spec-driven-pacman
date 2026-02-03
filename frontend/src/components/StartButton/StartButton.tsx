import styles from "./StartButton.module.css";

interface StartButtonProps {
  /** Custom button text. Defaults to "START GAME" */
  label?: string;
  /** Click handler */
  onClick?: () => void;
}

export function StartButton({
  label = "START GAME",
  onClick,
}: StartButtonProps) {
  return (
    <button className={styles.button} onClick={onClick} type="button">
      {label}
    </button>
  );
}


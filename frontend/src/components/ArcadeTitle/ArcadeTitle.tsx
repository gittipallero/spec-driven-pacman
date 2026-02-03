import styles from "./ArcadeTitle.module.css";

interface ArcadeTitleProps {
  /** Custom title text. Defaults to "Spec-Driven-Pacman" */
  title?: string;
}

export function ArcadeTitle({
  title = "Spec-Driven-Pacman",
}: ArcadeTitleProps) {
  return <h1 className={styles.title}>{title}</h1>;
}


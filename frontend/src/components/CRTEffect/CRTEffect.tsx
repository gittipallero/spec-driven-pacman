import type { ReactNode } from "react";
import styles from "./CRTEffect.module.css";

interface CRTEffectProps {
  /** Content to wrap with CRT effect */
  children: ReactNode;
}

export function CRTEffect({ children }: CRTEffectProps) {
  return (
    <div className={styles.crtContainer}>
      <div className={styles.scanLines} aria-hidden="true" />
      <div className={styles.content}>{children}</div>
    </div>
  );
}


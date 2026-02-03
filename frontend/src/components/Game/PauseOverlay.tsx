import { memo } from 'react';
import styles from './PauseOverlay.module.css';

interface PauseOverlayProps {
    visible: boolean;
}

function PauseOverlayComponent({ visible }: PauseOverlayProps) {
    if (!visible) {
        return null;
    }

    return (
        <div className={styles.overlay} data-testid="pause-overlay">
            <div className={styles.content}>
                <h2 className={styles.title}>PAUSED</h2>
                <p className={styles.hint}>Press ESC to resume</p>
            </div>
        </div>
    );
}

export const PauseOverlay = memo(PauseOverlayComponent);
export default PauseOverlay;


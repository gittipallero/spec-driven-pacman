import { memo, useMemo } from 'react';
import type { Direction, GhostState, GhostName } from '../../stores/gameStore';
import styles from './Ghost.module.css';

interface GhostProps {
    name: GhostName;
    x: number;
    y: number;
    direction: Direction;
    state: GhostState;
    cellSize?: number;
}

const GHOST_COLORS: Record<GhostName, string> = {
    blinky: '#ff0000', // Red
    pinky: '#ffb8ff',  // Pink
    inky: '#00ffff',   // Cyan
    clyde: '#ffb852',  // Orange
};

function GhostComponent({
    name,
    x,
    y,
    direction,
    state,
    cellSize = 16,
}: GhostProps) {
    const style = useMemo(() => ({
        left: `${x * cellSize}px`,
        top: `${y * cellSize}px`,
        width: `${cellSize}px`,
        height: `${cellSize}px`,
        '--ghost-color': state === 'vulnerable' ? '#2121de' : GHOST_COLORS[name],
    } as React.CSSProperties), [x, y, cellSize, name, state]);

    const className = useMemo(() => {
        const classes = [styles.ghost];
        if (state === 'vulnerable') classes.push(styles.vulnerable);
        if (state === 'eaten') classes.push(styles.eaten);
        if (state === 'respawning') classes.push(styles.respawning);
        return classes.join(' ');
    }, [state]);

    // Eaten ghosts only show eyes
    if (state === 'eaten') {
        return (
            <div
                className={className}
                style={style}
                data-testid={`ghost-${name}`}
                data-state={state}
            >
                <div className={styles.eyesOnly}>
                    <div className={styles.eye}>
                        <div className={styles.pupil} data-direction={direction} />
                    </div>
                    <div className={styles.eye}>
                        <div className={styles.pupil} data-direction={direction} />
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div
            className={className}
            style={style}
            data-testid={`ghost-${name}`}
            data-state={state}
            data-name={name}
        >
            <div className={styles.body}>
                <div className={styles.head}>
                    <div className={styles.eyes}>
                        <div className={styles.eye}>
                            <div className={styles.pupil} data-direction={direction} />
                        </div>
                        <div className={styles.eye}>
                            <div className={styles.pupil} data-direction={direction} />
                        </div>
                    </div>
                </div>
                <div className={styles.skirt}>
                    <div className={styles.wave} />
                    <div className={styles.wave} />
                    <div className={styles.wave} />
                </div>
            </div>
        </div>
    );
}

export const Ghost = memo(GhostComponent);
export default Ghost;


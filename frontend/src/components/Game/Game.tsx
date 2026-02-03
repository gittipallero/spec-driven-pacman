import { useEffect, useCallback, useState } from 'react';
import { useGameStore } from '../../stores/gameStore';
import { api } from '../../services/api';
import { websocket } from '../../services/websocket';
import { useKeyboard } from '../../hooks/useKeyboard';
import { Maze } from '../Maze';
import { PacMan } from '../PacMan';
import { Ghost } from '../Ghost';
import { HUD } from '../HUD';
import { GameOver } from '../GameOver';
import { PauseOverlay } from './PauseOverlay';
import styles from './Game.module.css';
import type { Direction, GameState, Ghost as GhostType } from '../../stores/gameStore';

interface GameProps {
    onExit?: () => void;
}

export function Game({ onExit }: GameProps) {
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const gameStatus = useGameStore((state: GameState) => state.gameStatus);
    const maze = useGameStore((state: GameState) => state.maze);
    const pacman = useGameStore((state: GameState) => state.pacman);
    const ghosts = useGameStore((state: GameState) => state.ghosts);
    const score = useGameStore((state: GameState) => state.score);
    const lives = useGameStore((state: GameState) => state.lives);
    const level = useGameStore((state: GameState) => state.level);
    const isConnected = useGameStore((state: GameState) => state.isConnected);
    const connectionError = useGameStore((state: GameState) => state.connectionError);

    const { setSessionId, setMaze, updateFromServer, setConnected, setConnectionError, reset } = useGameStore();

    // Start game on mount
    useEffect(() => {
        let mounted = true;

        const startGame = async () => {
            try {
                setIsLoading(true);
                setError(null);

                // Start game via REST API
                const response = await api.startGame();

                if (!mounted) return;

                setSessionId(response.sessionId);
                setMaze(response.maze);
                updateFromServer(response.state);

// Set up WebSocket handlers BEFORE connecting
        websocket.onMessage((state) => {
          if (mounted) {
            updateFromServer(state);
          }
        });

        websocket.onConnection((connected) => {
          if (mounted) {
            setConnected(connected);
          }
        });

        websocket.onErrorMessage((code, message) => {
          if (mounted) {
            setConnectionError(`${code}: ${message}`);
          }
        });

        // Now connect to WebSocket
        await websocket.connect(response.sessionId);

        if (!mounted) return;

        setIsLoading(false);
            } catch (err) {
                if (mounted) {
                    setError(err instanceof Error ? err.message : 'Failed to start game');
                    setIsLoading(false);
                }
            }
        };

        startGame();

        // Cleanup on unmount
        return () => {
            mounted = false;
            websocket.disconnect();
            reset();
        };
    }, [setSessionId, setMaze, updateFromServer, setConnected, setConnectionError, reset]);

    // Handle direction input
    const handleDirection = useCallback(
        (direction: Direction) => {
            if (gameStatus === 'playing') {
                websocket.sendInput(direction);
            }
        },
        [gameStatus]
    );

    // Handle pause toggle
    const handlePause = useCallback(() => {
        if (gameStatus === 'playing') {
            websocket.sendPause();
        } else if (gameStatus === 'paused') {
            websocket.sendResume();
        }
    }, [gameStatus]);

    // Handle restart (play again or continue)
    const handleRestart = useCallback(async () => {
        try {
            setIsLoading(true);
            websocket.disconnect();
            reset();

            const response = await api.startGame();
            setSessionId(response.sessionId);
            setMaze(response.maze);
            updateFromServer(response.state);

            await websocket.connect(response.sessionId);
            setIsLoading(false);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to restart game');
            setIsLoading(false);
        }
    }, [reset, setSessionId, setMaze, updateFromServer]);

    // Handle main menu
    const handleMainMenu = useCallback(() => {
        websocket.disconnect();
        reset();
        onExit?.();
    }, [reset, onExit]);

    // Keyboard controls
    useKeyboard({
        onDirection: handleDirection,
        onPause: handlePause,
        enabled: gameStatus === 'playing' || gameStatus === 'paused',
    });

    // Loading state
    if (isLoading) {
        return (
            <div className={styles.container}>
                <div className={styles.loading}>
                    <p className={styles.loadingText}>LOADING...</p>
                </div>
            </div>
        );
    }

    // Error state
    if (error) {
        return (
            <div className={styles.container}>
                <div className={styles.error}>
                    <p className={styles.errorTitle}>ERROR</p>
                    <p className={styles.errorMessage}>{error}</p>
                    <button className={styles.retryButton} onClick={() => window.location.reload()}>
                        RETRY
                    </button>
                    {onExit && (
                        <button className={styles.exitButton} onClick={onExit}>
                            MAIN MENU
                        </button>
                    )}
                </div>
            </div>
        );
    }

    // No maze loaded
    if (!maze) {
        return (
            <div className={styles.container}>
                <div className={styles.error}>
                    <p className={styles.errorMessage}>No maze data</p>
                </div>
            </div>
        );
    }

    const isPaused = gameStatus === 'paused';

    return (
        <div className={styles.container}>
            <div className={styles.gameArea}>
                {/* HUD */}
                <HUD score={score} lives={lives} level={level} />

                <div className={styles.mazeContainer}>
                    <Maze maze={maze} cellSize={16} />

                    {/* Render Pac-Man */}
                    {pacman && (
                        <PacMan
                            x={pacman.x}
                            y={pacman.y}
                            direction={pacman.direction}
                            state={pacman.state}
                            cellSize={16}
                        />
                    )}

                    {/* Render Ghosts */}
                    {ghosts.map((ghost: GhostType) => (
                        <Ghost
                            key={ghost.name}
                            name={ghost.name}
                            x={ghost.x}
                            y={ghost.y}
                            direction={ghost.direction}
                            state={ghost.state}
                            cellSize={16}
                        />
                    ))}

                    {/* Pause overlay */}
                    <PauseOverlay visible={isPaused} />

                    {/* Game over / Level complete overlay */}
                    <GameOver
                        status={gameStatus}
                        score={score}
                        level={level}
                        onRestart={handleRestart}
                        onMainMenu={handleMainMenu}
                    />
                </div>

                {/* Connection status */}
                {!isConnected && !connectionError && (
                    <div className={styles.connectionStatus}>
                        Reconnecting...
                    </div>
                )}

                {/* Connection lost */}
                {connectionError && (
                    <div className={styles.connectionLost}>
                        <p className={styles.connectionLostTitle}>CONNECTION LOST</p>
                        <p className={styles.connectionLostMessage}>
                            Please check your connection and try again
                        </p>
                    </div>
                )}
            </div>
        </div>
    );
}

export default Game;


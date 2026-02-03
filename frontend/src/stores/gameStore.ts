import { create } from 'zustand';

// Types matching backend
export type Direction = 'up' | 'down' | 'left' | 'right' | 'none';
export type GameStatus = 'waiting' | 'playing' | 'paused' | 'gameOver' | 'levelComplete';
export type PacManState = 'idle' | 'moving' | 'dying';
export type GhostState = 'normal' | 'vulnerable' | 'eaten' | 'respawning';
export type GhostMode = 'chase' | 'scatter' | 'frightened';
export type GhostName = 'blinky' | 'pinky' | 'inky' | 'clyde';

export interface Position {
  x: number;
  y: number;
}

export interface PacMan {
  x: number;
  y: number;
  direction: Direction;
  state: PacManState;
}

export interface Ghost {
  name: GhostName;
  x: number;
  y: number;
  direction: Direction;
  state: GhostState;
  mode: GhostMode;
}

export interface Maze {
  width: number;
  height: number;
  tiles: number[][];
}

export interface GameState {
  sessionId: string | null;
  gameStatus: GameStatus;
  score: number;
  lives: number;
  level: number;
  pacman: PacMan | null;
  ghosts: Ghost[];
  maze: Maze | null;
  dotsRemaining: number;
  vulnerabilityTimer: number;
  tick: number;
  isConnected: boolean;
  connectionError: string | null;
}

export interface GameActions {
  setSessionId: (id: string) => void;
  setMaze: (maze: Maze) => void;
  updateFromServer: (state: ServerState) => void;
  reset: () => void;
  setConnected: (connected: boolean) => void;
  setConnectionError: (error: string | null) => void;
}

export interface ServerState {
  status: GameStatus;
  score: number;
  lives: number;
  level: number;
  pacman: PacMan;
  ghosts: Ghost[];
  dotsRemaining: number;
  vulnerabilityTimer: number;
  tick?: number;
}

const initialState: GameState = {
  sessionId: null,
  gameStatus: 'waiting',
  score: 0,
  lives: 3,
  level: 1,
  pacman: null,
  ghosts: [],
  maze: null,
  dotsRemaining: 0,
  vulnerabilityTimer: 0,
  tick: 0,
  isConnected: false,
  connectionError: null,
};

export const useGameStore = create<GameState & GameActions>((set) => ({
  ...initialState,

  setSessionId: (id: string) => set({ sessionId: id }),

  setMaze: (maze: Maze) => set({ maze }),

  updateFromServer: (state: ServerState) =>
    set({
      gameStatus: state.status,
      score: state.score,
      lives: state.lives,
      level: state.level,
      pacman: state.pacman,
      ghosts: state.ghosts,
      dotsRemaining: state.dotsRemaining,
      vulnerabilityTimer: state.vulnerabilityTimer,
      tick: state.tick ?? 0,
    }),

  reset: () => set(initialState),

  setConnected: (connected: boolean) => set({ isConnected: connected }),

  setConnectionError: (error: string | null) => set({ connectionError: error }),
}));

// Selectors for optimized re-renders
export const selectSessionId = (state: GameState) => state.sessionId;
export const selectGameStatus = (state: GameState) => state.gameStatus;
export const selectScore = (state: GameState) => state.score;
export const selectLives = (state: GameState) => state.lives;
export const selectLevel = (state: GameState) => state.level;
export const selectPacman = (state: GameState) => state.pacman;
export const selectGhosts = (state: GameState) => state.ghosts;
export const selectMaze = (state: GameState) => state.maze;
export const selectDotsRemaining = (state: GameState) => state.dotsRemaining;
export const selectIsConnected = (state: GameState) => state.isConnected;
export const selectIsPaused = (state: GameState) => state.gameStatus === 'paused';
export const selectIsGameOver = (state: GameState) => state.gameStatus === 'gameOver';
export const selectIsLevelComplete = (state: GameState) => state.gameStatus === 'levelComplete';


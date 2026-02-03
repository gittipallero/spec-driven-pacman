import { useRef, useEffect, useCallback, memo } from 'react';
import type { Maze as MazeType } from '../../stores/gameStore';
import styles from './Maze.module.css';

// Tile types matching backend
const TILE_EMPTY = 0;
const TILE_WALL = 1;
const TILE_DOT = 2;
const TILE_POWER = 3;
const TILE_GHOST_HOUSE = 4;
const TILE_TUNNEL = 5;

// Colors
const COLORS = {
  background: '#000000',
  wall: '#2121de',
  wallBorder: '#0000aa',
  dot: '#ffb897',
  powerPellet: '#ffb897',
  corridor: '#000000',
  ghostHouse: '#111111',
  tunnel: '#000000',
};

interface MazeProps {
  maze: MazeType;
  cellSize?: number;
}

function MazeComponent({ maze, cellSize = 16 }: MazeProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const staticLayerRef = useRef<HTMLCanvasElement | null>(null);
  const animationFrameRef = useRef<number>(0);
  const powerPelletPhaseRef = useRef<number>(0);

  const width = maze.width * cellSize;
  const height = maze.height * cellSize;

  // Render static maze layer (walls)
  const renderStaticLayer = useCallback(() => {
    if (!staticLayerRef.current) {
      staticLayerRef.current = document.createElement('canvas');
      staticLayerRef.current.width = width;
      staticLayerRef.current.height = height;
    }

    const ctx = staticLayerRef.current.getContext('2d');
    if (!ctx) return;

    ctx.fillStyle = COLORS.background;
    ctx.fillRect(0, 0, width, height);

    for (let y = 0; y < maze.height; y++) {
      for (let x = 0; x < maze.width; x++) {
        const tile = maze.tiles[y][x];
        const px = x * cellSize;
        const py = y * cellSize;

        switch (tile) {
          case TILE_WALL:
            ctx.fillStyle = COLORS.wall;
            ctx.fillRect(px, py, cellSize, cellSize);
            // Add subtle border effect
            ctx.strokeStyle = COLORS.wallBorder;
            ctx.lineWidth = 1;
            ctx.strokeRect(px + 0.5, py + 0.5, cellSize - 1, cellSize - 1);
            break;
          case TILE_GHOST_HOUSE:
            ctx.fillStyle = COLORS.ghostHouse;
            ctx.fillRect(px, py, cellSize, cellSize);
            break;
          case TILE_TUNNEL:
          case TILE_EMPTY:
          default:
            ctx.fillStyle = COLORS.corridor;
            ctx.fillRect(px, py, cellSize, cellSize);
            break;
        }
      }
    }
  }, [maze, cellSize, width, height]);

  // Render dynamic layer (dots, power pellets)
  const renderDynamicLayer = useCallback(
    (ctx: CanvasRenderingContext2D) => {
      // Draw dots and power pellets
      for (let y = 0; y < maze.height; y++) {
        for (let x = 0; x < maze.width; x++) {
          const tile = maze.tiles[y][x];
          const px = x * cellSize + cellSize / 2;
          const py = y * cellSize + cellSize / 2;

          if (tile === TILE_DOT) {
            ctx.fillStyle = COLORS.dot;
            ctx.beginPath();
            ctx.arc(px, py, cellSize / 8, 0, Math.PI * 2);
            ctx.fill();
          } else if (tile === TILE_POWER) {
            // Flashing power pellet
            const alpha = 0.5 + 0.5 * Math.sin(powerPelletPhaseRef.current);
            ctx.fillStyle = COLORS.powerPellet;
            ctx.globalAlpha = alpha;
            ctx.beginPath();
            ctx.arc(px, py, cellSize / 3, 0, Math.PI * 2);
            ctx.fill();
            ctx.globalAlpha = 1;
          }
        }
      }
    },
    [maze, cellSize]
  );

  // Animation loop
  const animate = useCallback(() => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext('2d');
    if (!ctx || !staticLayerRef.current) return;

    // Update power pellet animation phase
    powerPelletPhaseRef.current += 0.1;

    // Draw static layer
    ctx.drawImage(staticLayerRef.current, 0, 0);

    // Draw dynamic layer
    renderDynamicLayer(ctx);

    animationFrameRef.current = requestAnimationFrame(animate);
  }, [renderDynamicLayer]);

  // Initialize static layer and start animation
  useEffect(() => {
    renderStaticLayer();
    animationFrameRef.current = requestAnimationFrame(animate);

    return () => {
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }
    };
  }, [renderStaticLayer, animate]);

  return (
    <canvas
      ref={canvasRef}
      width={width}
      height={height}
      className={styles.maze}
      data-testid="maze-canvas"
    />
  );
}

export const Maze = memo(MazeComponent);
export default Maze;


import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Maze } from './Maze';
import type { Maze as MazeType } from '../../stores/gameStore';

// Mock canvas context
const mockContext = {
  fillStyle: '',
  strokeStyle: '',
  lineWidth: 0,
  globalAlpha: 1,
  fillRect: vi.fn(),
  strokeRect: vi.fn(),
  beginPath: vi.fn(),
  arc: vi.fn(),
  fill: vi.fn(),
  drawImage: vi.fn(),
  getImageData: vi.fn(),
  putImageData: vi.fn(),
};

vi.stubGlobal('HTMLCanvasElement', class {
  width = 0;
  height = 0;
  getContext() {
    return mockContext;
  }
});

describe('Maze', () => {
  const createTestMaze = (width: number, height: number): MazeType => {
    const tiles: number[][] = [];
    for (let y = 0; y < height; y++) {
      tiles[y] = [];
      for (let x = 0; x < width; x++) {
        // Create a simple maze pattern
        if (x === 0 || x === width - 1 || y === 0 || y === height - 1) {
          tiles[y][x] = 1; // Wall
        } else if (x === 1 && y === 1) {
          tiles[y][x] = 3; // Power pellet
        } else {
          tiles[y][x] = 2; // Dot
        }
      }
    }
    return { width, height, tiles };
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders canvas element', () => {
    const maze = createTestMaze(5, 5);
    render(<Maze maze={maze} />);

    const canvas = screen.getByTestId('maze-canvas');
    expect(canvas).toBeInTheDocument();
    expect(canvas.tagName).toBe('CANVAS');
  });

  it('sets correct canvas dimensions', () => {
    const maze = createTestMaze(28, 31);
    const cellSize = 16;
    render(<Maze maze={maze} cellSize={cellSize} />);

    const canvas = screen.getByTestId('maze-canvas');
    expect(canvas).toHaveAttribute('width', String(28 * cellSize));
    expect(canvas).toHaveAttribute('height', String(31 * cellSize));
  });

  it('uses default cell size of 16', () => {
    const maze = createTestMaze(10, 10);
    render(<Maze maze={maze} />);

    const canvas = screen.getByTestId('maze-canvas');
    expect(canvas).toHaveAttribute('width', '160');
    expect(canvas).toHaveAttribute('height', '160');
  });

  it('accepts custom cell size', () => {
    const maze = createTestMaze(10, 10);
    render(<Maze maze={maze} cellSize={20} />);

    const canvas = screen.getByTestId('maze-canvas');
    expect(canvas).toHaveAttribute('width', '200');
    expect(canvas).toHaveAttribute('height', '200');
  });
});


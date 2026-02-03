import type { Maze, ServerState } from '../stores/gameStore';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface GameStartResponse {
  sessionId: string;
  websocketUrl: string;
  state: ServerState;
  maze: Maze;
}

export interface ApiError {
  code: string;
  message: string;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error: ApiError = await response.json().catch(() => ({
        code: 'UNKNOWN_ERROR',
        message: response.statusText,
      }));
      throw new Error(error.message);
    }

    // Handle 204 No Content
    if (response.status === 204) {
      return undefined as T;
    }

    return response.json();
  }

  /**
   * Start a new game session
   */
  async startGame(): Promise<GameStartResponse> {
    return this.request<GameStartResponse>('/api/game/start', {
      method: 'POST',
    });
  }

  /**
   * Get the current state of a game session
   */
  async getGameState(sessionId: string): Promise<ServerState> {
    return this.request<ServerState>(`/api/game/${sessionId}`);
  }

  /**
   * Pause a game session
   */
  async pauseGame(sessionId: string): Promise<ServerState> {
    return this.request<ServerState>(`/api/game/${sessionId}/pause`, {
      method: 'POST',
    });
  }

  /**
   * Resume a paused game session
   */
  async resumeGame(sessionId: string): Promise<ServerState> {
    return this.request<ServerState>(`/api/game/${sessionId}/resume`, {
      method: 'POST',
    });
  }

  /**
   * End a game session
   */
  async endGame(sessionId: string): Promise<void> {
    return this.request<void>(`/api/game/${sessionId}`, {
      method: 'DELETE',
    });
  }
}

// Export singleton instance
export const api = new ApiClient(API_URL);

// Export class for testing
export { ApiClient };


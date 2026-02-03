import type { Direction, ServerState } from '../stores/gameStore';

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080';

export type MessageHandler = (state: ServerState & { tick: number }) => void;
export type ErrorHandler = (code: string, message: string) => void;
export type ConnectionHandler = (connected: boolean) => void;

interface ClientMessage {
  type: 'input' | 'pause' | 'resume';
  direction?: Direction;
  timestamp?: number;
}

interface ServerMessage {
  type: 'state' | 'event' | 'error';
  tick?: number;
  status?: string;
  score?: number;
  lives?: number;
  level?: number;
  pacman?: unknown;
  ghosts?: unknown;
  dotsRemaining?: number;
  vulnerabilityTimer?: number;
  code?: string;
  message?: string;
}

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private sessionId: string | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private inputQueue: ClientMessage[] = [];
  private isReconnecting = false;

  private onStateUpdate: MessageHandler | null = null;
  private onError: ErrorHandler | null = null;
  private onConnectionChange: ConnectionHandler | null = null;

  /**
   * Connect to the game WebSocket
   */
  connect(sessionId: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this.sessionId = sessionId;
      const url = `${WS_URL}/ws/game/${sessionId}`;

      try {
        this.ws = new WebSocket(url);

        this.ws.onopen = () => {
          console.log('WebSocket connected');
          this.reconnectAttempts = 0;
          this.isReconnecting = false;
          this.onConnectionChange?.(true);
          
          // Send queued inputs
          this.flushInputQueue();
          resolve();
        };

        this.ws.onmessage = (event) => {
          this.handleMessage(event.data);
        };

        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          reject(error);
        };

        this.ws.onclose = (event) => {
          console.log('WebSocket closed:', event.code, event.reason);
          this.onConnectionChange?.(false);
          
          // Attempt reconnection if not intentionally closed
          if (!event.wasClean && this.sessionId) {
            this.attemptReconnect();
          }
        };
      } catch (error) {
        reject(error);
      }
    });
  }

  /**
   * Disconnect from the WebSocket
   */
  disconnect(): void {
    this.sessionId = null;
    if (this.ws) {
      this.ws.close(1000, 'Client disconnect');
      this.ws = null;
    }
  }

  /**
   * Send a direction input to the server
   */
  sendInput(direction: Direction): void {
    const message: ClientMessage = {
      type: 'input',
      direction,
      timestamp: Date.now(),
    };

    if (this.isConnected()) {
      this.send(message);
    } else {
      // Queue input for when connection restores
      this.inputQueue.push(message);
    }
  }

  /**
   * Send pause command
   */
  sendPause(): void {
    this.send({ type: 'pause' });
  }

  /**
   * Send resume command
   */
  sendResume(): void {
    this.send({ type: 'resume' });
  }

  /**
   * Set the state update handler
   */
  onMessage(handler: MessageHandler): void {
    this.onStateUpdate = handler;
  }

  /**
   * Set the error handler
   */
  onErrorMessage(handler: ErrorHandler): void {
    this.onError = handler;
  }

  /**
   * Set the connection change handler
   */
  onConnection(handler: ConnectionHandler): void {
    this.onConnectionChange = handler;
  }

  /**
   * Check if connected
   */
  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  private send(message: ClientMessage): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    }
  }

  private handleMessage(data: string): void {
    try {
      const message: ServerMessage = JSON.parse(data);

      switch (message.type) {
        case 'state':
          if (this.onStateUpdate) {
            this.onStateUpdate({
              tick: message.tick ?? 0,
              status: message.status as ServerState['status'],
              score: message.score ?? 0,
              lives: message.lives ?? 0,
              level: message.level ?? 1,
              pacman: message.pacman as ServerState['pacman'],
              ghosts: message.ghosts as ServerState['ghosts'],
              dotsRemaining: message.dotsRemaining ?? 0,
              vulnerabilityTimer: message.vulnerabilityTimer ?? 0,
            });
          }
          break;

        case 'error':
          if (this.onError && message.code && message.message) {
            this.onError(message.code, message.message);
          }
          break;

        case 'event':
          // Handle game events (future enhancement)
          console.log('Game event:', message);
          break;
      }
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error);
    }
  }

  private attemptReconnect(): void {
    if (this.isReconnecting || this.reconnectAttempts >= this.maxReconnectAttempts) {
      if (this.reconnectAttempts >= this.maxReconnectAttempts) {
        this.onError?.('CONNECTION_LOST', 'Failed to reconnect after multiple attempts');
      }
      return;
    }

    this.isReconnecting = true;
    this.reconnectAttempts++;

    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
    console.log(`Attempting reconnect ${this.reconnectAttempts}/${this.maxReconnectAttempts} in ${delay}ms`);

    setTimeout(() => {
      if (this.sessionId) {
        this.connect(this.sessionId).catch(() => {
          this.isReconnecting = false;
          this.attemptReconnect();
        });
      }
    }, delay);
  }

  private flushInputQueue(): void {
    while (this.inputQueue.length > 0) {
      const message = this.inputQueue.shift();
      if (message) {
        this.send(message);
      }
    }
  }
}

// Export singleton instance
export const websocket = new WebSocketClient();


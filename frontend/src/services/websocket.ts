import type { WSMessage } from '../types/ws';

type MessageHandler = (msg: WSMessage) => void;

export class WebSocketService {
  private ws: WebSocket | null = null;
  private url: string;
  private handlers: Map<string, Set<MessageHandler>> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 10;
  private reconnectDelay = 1000;
  private heartbeatInterval: ReturnType<typeof setInterval> | null = null;
  private onStatusChange?: (status: 'Connecting' | 'Connected' | 'Disconnected') => void;

  constructor(url: string, onStatusChange?: typeof this.onStatusChange) {
    this.url = url;
    this.onStatusChange = onStatusChange;
  }

  connect(): void {
    this.onStatusChange?.('Connecting');
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
      this.reconnectDelay = 1000;
      this.onStatusChange?.('Connected');
      this.startHeartbeat();
    };

    this.ws.onmessage = (event: MessageEvent) => {
      try {
        const msg: WSMessage = JSON.parse(event.data);
        this.dispatch(msg);
      } catch (error) {
        console.error('WS parse error:', error);
      }
    };

    this.ws.onclose = () => {
      this.onStatusChange?.('Disconnected');
      this.stopHeartbeat();
      this.scheduleReconnect();
    };

    this.ws.onerror = (err) => {
      console.error('WS error:', err);
      this.ws?.close();
    };
  }

  send<T>(type: string, payload: T): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, payload }));
    }
  }

  on<T>(type: string, handler: (msg: WSMessage<T>) => void): () => void {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, new Set());
    }
    this.handlers.get(type)?.add(handler as MessageHandler);
    return () => this.handlers.get(type)?.delete(handler as MessageHandler);
  }

  private dispatch(msg: WSMessage): void {
    this.handlers.get(msg.type)?.forEach((handler) => handler(msg));
    this.handlers.get('*')?.forEach((handler) => handler(msg));
  }

  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnect attempts reached');
      return;
    }
    this.reconnectAttempts += 1;
    const delay = Math.min(this.reconnectDelay * Math.pow(1.5, this.reconnectAttempts), 30000);
    setTimeout(() => this.connect(), delay);
  }

  private startHeartbeat(): void {
    this.heartbeatInterval = setInterval(() => {
      this.send('PING', {});
    }, 30000);
  }

  private stopHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }

  disconnect(): void {
    this.maxReconnectAttempts = 0;
    this.ws?.close();
  }
}

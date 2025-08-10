import { WebSocketMessage } from '../types';

type MessageHandler = (message: WebSocketMessage) => void;
type ConnectionHandler = () => void;
type ErrorHandler = (error: Event) => void;

class WebSocketService {
  private socket: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private messageHandlers: MessageHandler[] = [];
  private connectionHandlers: ConnectionHandler[] = [];
  private errorHandlers: ErrorHandler[] = [];
  private currentRoomId: string | null = null;
  private isConnecting = false;

  connect(roomId: string): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        // Disconnect existing connection if any
        if (this.socket) {
          this.socket.close();
          this.socket = null;
        }

        if (this.isConnecting) {
          reject(new Error('Connection already in progress'));
          return;
        }

        this.isConnecting = true;
        
        // Construct WebSocket URL - use the same host as the current page but with the backend port
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.hostname;
        const port = '5000'; // Backend port is always 5000
        const wsUrl = `${protocol}//${host}:${port}/ws?room=${roomId}`;
        
        console.log('🔌 Connecting to WebSocket:', wsUrl);
        this.socket = new WebSocket(wsUrl);
        this.currentRoomId = roomId;

        this.socket.onopen = () => {
          console.log('✅ WebSocket connected to room:', roomId);
          this.reconnectAttempts = 0;
          this.isConnecting = false;
          this.connectionHandlers.forEach(handler => handler());
          resolve();
        };

        this.socket.onmessage = (event) => {
          try {
            const message: WebSocketMessage = JSON.parse(event.data);
            console.log('📨 Received WebSocket message:', message);
            this.messageHandlers.forEach(handler => handler(message));
          } catch (error) {
            console.error('❌ Failed to parse WebSocket message:', error);
          }
        };

        this.socket.onclose = (event) => {
          console.log('❌ WebSocket closed:', event.code, event.reason);
          this.isConnecting = false;
          if (this.currentRoomId === roomId) {
            this.handleReconnect(roomId);
          }
        };

        this.socket.onerror = (error) => {
          console.error('⚠️ WebSocket error:', error);
          this.isConnecting = false;
          this.errorHandlers.forEach(handler => handler(error));
          reject(error);
        };
      } catch (error) {
        this.isConnecting = false;
        reject(error);
      }
    });
  }

  private handleReconnect(roomId: string) {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(`🔄 Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);
      
      setTimeout(() => {
        this.connect(roomId).catch(error => {
          console.error('❌ Reconnection failed:', error);
        });
      }, this.reconnectDelay * this.reconnectAttempts);
    } else {
      console.error('❌ Max reconnection attempts reached');
    }
  }

  disconnect() {
    if (this.socket) {
      console.log('🔌 Disconnecting WebSocket');
      this.socket.close();
      this.socket = null;
      this.currentRoomId = null;
      this.isConnecting = false;
    }
  }

  send(message: WebSocketMessage) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      console.log('📤 Sending WebSocket message:', message);
      this.socket.send(JSON.stringify(message));
    } else {
      console.warn('⚠️ WebSocket is not connected');
    }
  }

  onMessage(handler: MessageHandler) {
    this.messageHandlers.push(handler);
  }

  onConnect(handler: ConnectionHandler) {
    this.connectionHandlers.push(handler);
  }

  onError(handler: ErrorHandler) {
    this.errorHandlers.push(handler);
  }

  isConnected(): boolean {
    return this.socket?.readyState === WebSocket.OPEN;
  }

  // Clear all handlers (useful for cleanup)
  clearHandlers() {
    this.messageHandlers = [];
    this.connectionHandlers = [];
    this.errorHandlers = [];
  }
}

const wsService = new WebSocketService();
export default wsService;


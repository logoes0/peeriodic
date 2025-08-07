// API Response types
export interface ApiResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
}

// Room types
export interface Room {
  id: string;
  title: string;
  content?: string;
  created_at?: string;
}

export interface CreateRoomRequest {
  title: string;
  uid: string;
}

export interface CreateRoomResponse {
  id: string;
  title: string;
}

// Document types
export interface SaveDocumentRequest {
  content: string;
}

export interface SaveDocumentResponse {
  status: string;
  roomId: string;
  contentLength: number;
}

// WebSocket message types
export interface WebSocketMessage {
  type: 'init' | 'update';
  data: string;
}

// Component props types
export interface RoomListProps {
  rooms: Room[];
  onRoomClick: (roomId: string) => void;
  onRoomRename: (roomId: string, newName: string) => void;
  onRoomDelete: (roomId: string) => void;
  onRoomShare: (roomId: string) => void;
}

export interface EditorProps {
  roomId: string;
  roomName: string;
  onBack: () => void;
  onShare: () => void;
}

export interface HomeProps {
  onCreateRoom: () => void;
  onRoomClick: (roomId: string) => void;
}


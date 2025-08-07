import { 
  ApiResponse, 
  Room, 
  CreateRoomRequest, 
  CreateRoomResponse,
  SaveDocumentRequest,
  SaveDocumentResponse 
} from '../types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000';

class ApiService {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string, 
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const defaultOptions: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
      },
      ...options,
    };

    try {
      const response = await fetch(url, defaultOptions);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Room API methods
  async getRooms(uid: string): Promise<Room[]> {
    const response = await this.request<ApiResponse<Room[]>>(`/api/rooms?uid=${uid}`);
    return response.data || [];
  }

  async createRoom(request: CreateRoomRequest): Promise<CreateRoomResponse> {
    const response = await this.request<ApiResponse<CreateRoomResponse>>('/api/rooms', {
      method: 'POST',
      body: JSON.stringify(request),
    });
    return response.data!;
  }

  async getRoom(roomId: string): Promise<Room> {
    const response = await this.request<ApiResponse<Room>>(`/api/rooms/${roomId}`);
    return response.data!;
  }

  async deleteRoom(roomId: string): Promise<void> {
    await this.request(`/api/rooms/${roomId}`, {
      method: 'DELETE',
    });
  }

  // Document API methods
  async saveDocument(roomId: string, content: string): Promise<SaveDocumentResponse> {
    const request: SaveDocumentRequest = { content };
    const response = await this.request<ApiResponse<SaveDocumentResponse>>(`/api/save?room=${roomId}`, {
      method: 'POST',
      body: JSON.stringify(request),
    });
    return response.data!;
  }

  async getDocument(roomId: string): Promise<{ id: string; content: string }> {
    const response = await this.request<ApiResponse<{ id: string; content: string }>>(`/api/rooms/${roomId}`);
    return response.data!;
  }
}

// Export singleton instance
export const apiService = new ApiService();
export default apiService;


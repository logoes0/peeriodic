import { Room } from '../types';

const STORAGE_KEYS = {
  ROOMS: 'myRooms',
  USER_ID: 'userId',
} as const;

export class StorageService {
  // Room management
  static getRooms(): Room[] {
    try {
      const rooms = localStorage.getItem(STORAGE_KEYS.ROOMS);
      return rooms ? JSON.parse(rooms) : [];
    } catch (error) {
      console.error('Failed to parse rooms from localStorage:', error);
      return [];
    }
  }

  static saveRooms(rooms: Room[]): void {
    try {
      localStorage.setItem(STORAGE_KEYS.ROOMS, JSON.stringify(rooms));
    } catch (error) {
      console.error('Failed to save rooms to localStorage:', error);
    }
  }

  static addRoom(room: Room): void {
    const rooms = this.getRooms();
    rooms.push(room);
    this.saveRooms(rooms);
  }

  static updateRoom(roomId: string, updates: Partial<Room>): void {
    const rooms = this.getRooms();
    const roomIndex = rooms.findIndex(room => room.id === roomId);
    
    if (roomIndex !== -1) {
      rooms[roomIndex] = { ...rooms[roomIndex], ...updates };
      this.saveRooms(rooms);
    }
  }

  static deleteRoom(roomId: string): void {
    const rooms = this.getRooms();
    const filteredRooms = rooms.filter(room => room.id !== roomId);
    this.saveRooms(filteredRooms);
  }

  static getRoom(roomId: string): Room | undefined {
    const rooms = this.getRooms();
    return rooms.find(room => room.id === roomId);
  }

  // User management
  static getUserId(): string {
    let userId = localStorage.getItem(STORAGE_KEYS.USER_ID);
    if (!userId) {
      userId = this.generateUserId();
      localStorage.setItem(STORAGE_KEYS.USER_ID, userId);
    }
    return userId;
  }

  private static generateUserId(): string {
    return `user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  // Clear all data
  static clearAll(): void {
    localStorage.removeItem(STORAGE_KEYS.ROOMS);
    localStorage.removeItem(STORAGE_KEYS.USER_ID);
  }
}


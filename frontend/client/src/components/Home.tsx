import React, { useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import { Room } from '../types';
import { StorageService } from '../utils/storage';
import { UrlUtils } from '../utils/url';
import RoomList from './RoomList';
import './Home.css';

interface HomeProps {
  onCreateRoom: () => void;
  onRoomClick: (roomId: string) => void;
}

const Home: React.FC<HomeProps> = ({ onCreateRoom, onRoomClick }) => {
  const [rooms, setRooms] = useState<Room[]>([]);

  useEffect(() => {
    const savedRooms = StorageService.getRooms();
    setRooms(savedRooms);
  }, []);

  const handleCreateRoom = () => {
    const roomId = uuidv4();
    const newRoom: Room = { 
      id: roomId, 
      title: 'New Room' 
    };
    
    StorageService.addRoom(newRoom);
    setRooms(prev => [...prev, newRoom]);
    onCreateRoom();
  };

  const handleRoomClick = (roomId: string) => {
    onRoomClick(roomId);
  };

  const handleRoomRename = (roomId: string, newName: string) => {
    StorageService.updateRoom(roomId, { title: newName });
    setRooms(prev => 
      prev.map(room => 
        room.id === roomId ? { ...room, title: newName } : room
      )
    );
  };

  const handleRoomDelete = (roomId: string) => {
    StorageService.deleteRoom(roomId);
    setRooms(prev => prev.filter(room => room.id !== roomId));
  };

  const handleRoomShare = async (roomId: string) => {
    const shareUrl = UrlUtils.getShareUrl(roomId);
    try {
      await UrlUtils.copyToClipboard(shareUrl);
      // Could add a toast notification here
      console.log('Share link copied to clipboard');
    } catch (error) {
      console.error('Failed to copy share link:', error);
    }
  };

  return (
    <div className="home">
      <header className="home-header">
        <h1>Peeriodic</h1>
        <p>Create and share documents in real-time</p>
      </header>

      <main className="home-main">
        <button 
          className="create-room-btn"
          onClick={handleCreateRoom}
        >
          Create New Room
        </button>

        {rooms.length > 0 && (
          <section className="rooms-section">
            <h2>Your Rooms</h2>
            <RoomList
              rooms={rooms}
              onRoomClick={handleRoomClick}
              onRoomRename={handleRoomRename}
              onRoomDelete={handleRoomDelete}
              onRoomShare={handleRoomShare}
            />
          </section>
        )}
      </main>
    </div>
  );
};

export default Home;


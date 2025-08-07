import React, { useState, useEffect } from 'react';
import { StorageService } from './utils/storage';
import { UrlUtils } from './utils/url';
import Home from './components/Home';
import Editor from './components/Editor';
import './App.css';

const App: React.FC = () => {
  const [currentRoomId, setCurrentRoomId] = useState<string | null>(null);
  const [roomName, setRoomName] = useState('Shared Room');

  useEffect(() => {
    // Check if we're in a room from URL
    const roomId = UrlUtils.getRoomIdFromUrl();
    if (roomId) {
      setCurrentRoomId(roomId);
      
      // Get room name from local storage
      const savedRooms = StorageService.getRooms();
      const matchedRoom = savedRooms.find(room => room.id === roomId);
      setRoomName(matchedRoom?.title || 'Shared Room');
    }
  }, []);

  const handleCreateRoom = () => {
    const roomId = UrlUtils.generateRoomId();
    UrlUtils.setRoomIdInUrl(roomId);
    setCurrentRoomId(roomId);
    setRoomName('New Room');
  };

  const handleRoomClick = (roomId: string) => {
    UrlUtils.setRoomIdInUrl(roomId);
    setCurrentRoomId(roomId);
    
    // Get room name from local storage
    const savedRooms = StorageService.getRooms();
    const matchedRoom = savedRooms.find(room => room.id === roomId);
    setRoomName(matchedRoom?.title || 'Shared Room');
  };

  const handleBackToHome = () => {
    UrlUtils.removeRoomIdFromUrl();
    setCurrentRoomId(null);
    setRoomName('Shared Room');
  };

  const handleShare = () => {
    // This could trigger a toast notification or other UI feedback
    console.log('Share action triggered');
  };

  return (
    <div className="App">
      {!currentRoomId ? (
        <Home
          onCreateRoom={handleCreateRoom}
          onRoomClick={handleRoomClick}
        />
      ) : (
        <Editor
          roomId={currentRoomId}
          roomName={roomName}
          onBack={handleBackToHome}
          onShare={handleShare}
        />
      )}
    </div>
  );
};

export default App;


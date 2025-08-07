import React, { useState } from 'react';
import { Room, RoomListProps } from '../types';
import './RoomList.css';

const RoomList: React.FC<RoomListProps> = ({
  rooms,
  onRoomClick,
  onRoomRename,
  onRoomDelete,
  onRoomShare,
}) => {
  const [editingRoomId, setEditingRoomId] = useState<string | null>(null);
  const [editedName, setEditedName] = useState('');

  const handleRenameClick = (room: Room) => {
    setEditingRoomId(room.id);
    setEditedName(room.title);
  };

  const handleRenameSave = (roomId: string) => {
    if (editedName.trim()) {
      onRoomRename(roomId, editedName.trim());
    }
    setEditingRoomId(null);
    setEditedName('');
  };

  const handleRenameCancel = () => {
    setEditingRoomId(null);
    setEditedName('');
  };

  const handleKeyDown = (e: React.KeyboardEvent, roomId: string) => {
    if (e.key === 'Enter') {
      handleRenameSave(roomId);
    } else if (e.key === 'Escape') {
      handleRenameCancel();
    }
  };

  return (
    <ul className="room-list">
      {rooms.map((room) => (
        <li key={room.id} className="room-item">
          <div className="room-content">
            {editingRoomId === room.id ? (
              <input
                type="text"
                value={editedName}
                onChange={(e) => setEditedName(e.target.value)}
                onBlur={() => handleRenameSave(room.id)}
                onKeyDown={(e) => handleKeyDown(e, room.id)}
                className="room-name-input"
                autoFocus
              />
            ) : (
              <a
                href={`/editor?room=${room.id}`}
                className="room-link"
                onClick={(e) => {
                  e.preventDefault();
                  onRoomClick(room.id);
                }}
              >
                {room.title}
              </a>
            )}
          </div>

          <div className="room-actions">
            {editingRoomId !== room.id && (
              <>
                <button
                  onClick={() => handleRenameClick(room)}
                  className="btn btn-rename"
                  title="Rename room"
                >
                  ✏️
                </button>
                <button
                  onClick={() => onRoomShare(room.id)}
                  className="btn btn-share"
                  title="Share room"
                >
                  🔗
                </button>
                <button
                  onClick={() => onRoomDelete(room.id)}
                  className="btn btn-delete"
                  title="Delete room"
                >
                  🗑️
                </button>
              </>
            )}
          </div>
        </li>
      ))}
    </ul>
  );
};

export default RoomList;


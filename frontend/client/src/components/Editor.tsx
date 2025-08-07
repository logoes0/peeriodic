import React, { useState, useEffect, useCallback, useRef } from 'react';
import { WebSocketMessage } from '../types';
import { StorageService } from '../utils/storage';
import { UrlUtils } from '../utils/url';
import apiService from '../services/api';
import wsService from '../services/websocket';
import './Editor.css';

interface EditorProps {
  roomId: string;
  roomName: string;
  onBack: () => void;
  onShare: () => void;
}

const Editor: React.FC<EditorProps> = ({ roomId, roomName, onBack, onShare }) => {
  const [document, setDocument] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isConnected, setIsConnected] = useState(false);
  const [copiedRoomId, setCopiedRoomId] = useState<string | null>(null);
  const autoSaveIntervalRef = useRef<NodeJS.Timeout | null>(null);

  // Load initial document content
  useEffect(() => {
    const loadDocument = async () => {
      try {
        setIsLoading(true);
        const response = await apiService.getDocument(roomId);
        setDocument(response.content || '');
      } catch (error) {
        console.error('Failed to load document:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadDocument();
  }, [roomId]);

  // Setup WebSocket connection
  useEffect(() => {
    const setupWebSocket = async () => {
      try {
        await wsService.connect(roomId);
        setIsConnected(true);

        wsService.onMessage((message: WebSocketMessage) => {
          if (message.type === 'init' || message.type === 'update') {
            setDocument(message.data);
          }
        });

        wsService.onConnect(() => {
          setIsConnected(true);
        });

        wsService.onError((error) => {
          console.error('WebSocket error:', error);
          setIsConnected(false);
        });
      } catch (error) {
        console.error('Failed to connect WebSocket:', error);
        setIsConnected(false);
      }
    };

    setupWebSocket();

    return () => {
      wsService.disconnect();
      setIsConnected(false);
    };
  }, [roomId]);

  // Auto-save functionality
  useEffect(() => {
    if (document.trim()) {
      autoSaveIntervalRef.current = setInterval(() => {
        handleSave();
      }, 30000); // Auto-save every 30 seconds
    }

    return () => {
      if (autoSaveIntervalRef.current) {
        clearInterval(autoSaveIntervalRef.current);
      }
    };
  }, [document]);

  const handleDocumentChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value;
    setDocument(newContent);

    // Send update via WebSocket
    if (wsService.isConnected()) {
      wsService.send({
        type: 'update',
        data: newContent,
      });
    }
  }, []);

  const handleSave = useCallback(async () => {
    if (!document.trim()) return;

    try {
      setIsSaving(true);
      await apiService.saveDocument(roomId, document);
    } catch (error) {
      console.error('Failed to save document:', error);
    } finally {
      setIsSaving(false);
    }
  }, [document, roomId]);

  const handleShare = async () => {
    const shareUrl = UrlUtils.getShareUrl(roomId);
    try {
      await UrlUtils.copyToClipboard(shareUrl);
      setCopiedRoomId(roomId);
      setTimeout(() => setCopiedRoomId(null), 2000);
      onShare();
    } catch (error) {
      console.error('Failed to copy share link:', error);
    }
  };

  return (
    <div className="editor">
      <header className="editor-header">
        <div className="editor-title">
          <h2>{roomName}</h2>
          <div className="connection-status">
            {isConnected ? '🟢 Connected' : '🔴 Disconnected'}
          </div>
        </div>

        <div className="editor-actions">
          <button onClick={onBack} className="btn btn-back">
            ⬅ Back to Home
          </button>
          <button onClick={handleShare} className="btn btn-share">
            {copiedRoomId === roomId ? '✅ Copied' : '🔗 Share'}
          </button>
          <button
            onClick={handleSave}
            disabled={isSaving}
            className={`btn btn-save ${isSaving ? 'btn-saving' : ''}`}
          >
            {isSaving ? '⏳ Saving...' : '💾 Save'}
          </button>
        </div>
      </header>

      <main className="editor-main">
        {isLoading ? (
          <div className="loading-indicator">
            <div className="spinner"></div>
            <p>Loading document...</p>
          </div>
        ) : (
          <textarea
            value={document}
            onChange={handleDocumentChange}
            placeholder="Start typing collaboratively..."
            className="document-editor"
            rows={20}
            cols={80}
          />
        )}
      </main>
    </div>
  );
};

export default Editor;


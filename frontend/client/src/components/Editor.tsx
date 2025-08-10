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
  const [copiedRoomId, setCopiedRoomId] = useState<string | null>(null);
  const autoSaveIntervalRef = useRef<NodeJS.Timeout | null>(null);
  const lastUpdateRef = useRef<string>('');
  const isLocalUpdateRef = useRef<boolean>(false);
  const lastSentContentRef = useRef<string>('');
  const wsConnectedRef = useRef<boolean>(false);

  // Load initial document content
  useEffect(() => {
    const loadDocument = async () => {
      try {
        setIsLoading(true);
        const response = await apiService.getDocument(roomId);
        const content = response.content || '';
        setDocument(content);
        lastUpdateRef.current = content;
        lastSentContentRef.current = content;
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
    let isActive = true;

    const setupWebSocket = async () => {
      try {
        // Clear any existing handlers
        wsService.clearHandlers();
        
        await wsService.connect(roomId);
        
        if (!isActive) return;

        wsConnectedRef.current = true;

        // Set up message handler
        wsService.onMessage((message: WebSocketMessage) => {
          if (!isActive) return;
          
          console.log('Editor received message:', message);
          if (message.type === 'init' || message.type === 'update') {
            // Only update if this is not a local update and content is different
            if (!isLocalUpdateRef.current && message.data !== lastUpdateRef.current) {
              console.log('🔄 Updating document from WebSocket:', message.data.substring(0, 50) + '...');
              setDocument(message.data);
              lastUpdateRef.current = message.data;
              lastSentContentRef.current = message.data;
            }
          }
        });

        wsService.onConnect(() => {
          if (!isActive) return;
          console.log('WebSocket connected in Editor');
          wsConnectedRef.current = true;
        });

        wsService.onError((error) => {
          if (!isActive) return;
          console.error('WebSocket error in Editor:', error);
          wsConnectedRef.current = false;
        });
      } catch (error) {
        if (!isActive) return;
        console.error('Failed to connect WebSocket:', error);
        wsConnectedRef.current = false;
      }
    };

    setupWebSocket();

    return () => {
      isActive = false;
      wsConnectedRef.current = false;
      wsService.disconnect();
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
    
    // Update local state immediately
    setDocument(newContent);
    
    // Only send WebSocket update if content has actually changed and we're connected
    if (wsConnectedRef.current && newContent !== lastSentContentRef.current) {
      isLocalUpdateRef.current = true;
      lastSentContentRef.current = newContent;
      
      console.log('📤 Sending document update via WebSocket');
      wsService.send({
        type: 'update',
        data: newContent,
      });

      // Reset the flag after a longer delay to account for network latency
      setTimeout(() => {
        isLocalUpdateRef.current = false;
      }, 1000); // Increased timeout to 1 second
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


import React, { useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import './App.css';

function App() {
    const [document, setDocument] = useState('');
    const [socket, setSocket] = useState(null);
    const params = new URLSearchParams(window.location.search);
    const room = params.get("room");

    useEffect(() => {
        if (!room) return;

        const newSocket = new WebSocket(`ws://localhost:5000/ws?room=${room}`);
        setSocket(newSocket);

        newSocket.onopen = () => {
            console.log('WebSocket connection established');
        };

        newSocket.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                if (message.type === 'init' || message.type === 'update') {
                    setDocument(message.data);
                }
            } catch (error) {
                console.error('Error parsing message:', error);
            }
        };

        newSocket.onclose = () => {
            console.log('WebSocket connection closed');
        };

        newSocket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        return () => {
            newSocket.close();
        };
    }, [room]);

    const handleChange = (e) => {
        const newDocument = e.target.value;
        setDocument(newDocument);
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ type: 'update', data: newDocument }));
        }
    };

    const createRoom = () => {
        const roomId = uuidv4();
        window.location.href = `/editor?room=${roomId}`;
    };

    return (
        <div className="App">
            {!room ? (
                <div>
                    <h1>Welcome to the Collaborative Editor</h1>
                    <button onClick={createRoom}>Create New Room</button>
                </div>
            ) : (
                <div>
                    <h2>Room: {room}</h2>
                    <textarea
                        value={document}
                        onChange={handleChange}
                        rows="20"
                        cols="80"
                        placeholder="Start typing collaboratively..."
                    />
                </div>
            )}
        </div>
    );
}

export default App;

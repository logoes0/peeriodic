import React, { useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import './App.css';

function App() {
    const [document, setDocument] = useState('');
    const [socket, setSocket] = useState(null);
    const [myRooms, setMyRooms] = useState([]);
    const params = new URLSearchParams(window.location.search);
    const room = params.get("room");

    // Load saved rooms from localStorage on first render
    useEffect(() => {
        const savedRooms = JSON.parse(localStorage.getItem("myRooms") || "[]");
        setMyRooms(savedRooms);
    }, []);

    // WebSocket setup for room
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

    // Create a new room and save it
    const createRoom = () => {
        const roomId = uuidv4();
        const existingRooms = JSON.parse(localStorage.getItem("myRooms") || "[]");
        const updatedRooms = [...existingRooms, roomId];
        localStorage.setItem("myRooms", JSON.stringify(updatedRooms));
        window.location.href = `/editor?room=${roomId}`;
    };

    const goBackToHome = () => {
        window.location.href = '/';
    };

    const handleChange = (e) => {
        const newDocument = e.target.value;
        setDocument(newDocument);
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ type: 'update', data: newDocument }));
        }
    };

    return (
        <div className="App">
            {!room ? (
                <div>
                    <h1>Welcome to the Collaborative Editor</h1>
                    <button onClick={createRoom}>Create New Room</button>

                    {myRooms.length > 0 && (
                        <div style={{ marginTop: '20px' }}>
                            <h2>Your Previous Rooms</h2>
                            <ul>
                                {myRooms.map((roomId) => (
                                    <li key={roomId}>
                                        <a href={`/editor?room=${roomId}`}>{roomId}</a>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    )}
                </div>
            ) : (
                <div>
                    <h2>Room: {room}</h2>
                    <button onClick={goBackToHome} style={{ marginBottom: '10px' }}>
                        ⬅ Back to Home
                    </button>
                    <br />
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

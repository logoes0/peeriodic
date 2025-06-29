import React, { useState, useEffect } from "react";
import { v4 as uuidv4 } from "uuid";
import "./App.css";

function App() {
    const [document, setDocument] = useState("");
    const [socket, setSocket] = useState(null);
    const [myRooms, setMyRooms] = useState([]);
    const [editingRoomId, setEditingRoomId] = useState(null);
    const [editedRoomName, setEditedRoomName] = useState("");
    const [roomName, setRoomName] = useState("");
    const [copiedRoomId, setCopiedRoomId] = useState(null);
    const [isSaving, setIsSaving] = useState(false);
    const [isLoading, setIsLoading] = useState(false);

    const params = new URLSearchParams(window.location.search);
    const room = params.get("room");

    useEffect(() => {
        const savedRooms = JSON.parse(localStorage.getItem("myRooms") || "[]");
        setMyRooms(savedRooms);

        if (room) {
            const matched = savedRooms.find((r) => r.id === room);
            setRoomName(matched ? matched.name : "Shared Room");
        }
    }, []);

    useEffect(() => {
        if (!room) return;

        const loadInitialContent = async () => {
            setIsLoading(true);
            try {
                const response = await fetch(
                    `http://localhost:5000/api/rooms/${room}`
                );
                if (!response.ok) throw new Error("Failed to load");
                const data = await response.json();
                setDocument(data.content || "");
            } catch (error) {
                console.error("Load error:", error);
            } finally {
                setIsLoading(false);
            }
        };

        loadInitialContent();

        const newSocket = new WebSocket(`ws://localhost:5000/ws?room=${room}`);
        setSocket(newSocket);

        newSocket.onopen = () => {
            console.log("✅ WebSocket connected");
        };

        newSocket.onmessage = (event) => {
            try {
                const message =
                    typeof event.data === "string"
                        ? JSON.parse(event.data)
                        : event.data;

                if (message.type === "init" || message.type === "update") {
                    setDocument(message.data);
                }
            } catch (error) {
                console.error(
                    "❌ JSON parse error:",
                    error.message,
                    event.data
                );
            }
        };

        newSocket.onclose = () => {
            console.log("❌ WebSocket closed");
        };

        newSocket.onerror = (error) => {
            console.error("⚠️ WebSocket error:", error);
        };

        const autoSaveInterval = setInterval(() => {
            if (document.trim()) {
                handleSave();
            }
        }, 30000);

        return () => {
            clearInterval(autoSaveInterval);
            newSocket.close();
        };
    }, [room]);

    const createRoom = () => {
        const roomId = uuidv4();
        const newRoom = { id: roomId, name: "New Room" };
        const updatedRooms = [...myRooms, newRoom];
        localStorage.setItem("myRooms", JSON.stringify(updatedRooms));
        setMyRooms(updatedRooms);
        window.location.href = `/editor?room=${roomId}`;
    };

    const goBackToHome = () => {
        window.location.href = "/";
    };

    const handleChange = (e) => {
        const newDocument = e.target.value;
        setDocument(newDocument);
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ type: "update", data: newDocument }));
        }
    };

    const handleSave = async () => {
        setIsSaving(true);
        try {
            const response = await fetch(
                `http://localhost:5000/api/save?room=${room}`,
                {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ content: document }),
                }
            );
            if (!response.ok) throw new Error("Failed to save document");
            await response.json();
        } catch (error) {
            console.error("Save error:", error);
        } finally {
            setIsSaving(false);
        }
    };

    const deleteRoom = (roomIdToDelete) => {
        const updatedRooms = myRooms.filter(
            (room) => room.id !== roomIdToDelete
        );
        setMyRooms(updatedRooms);
        localStorage.setItem("myRooms", JSON.stringify(updatedRooms));
    };

    const startRename = (roomId, currentName) => {
        setEditingRoomId(roomId);
        setEditedRoomName(currentName);
    };

    const saveRename = (roomId) => {
        const updatedRooms = myRooms.map((room) =>
            room.id === roomId ? { ...room, name: editedRoomName } : room
        );
        setMyRooms(updatedRooms);
        localStorage.setItem("myRooms", JSON.stringify(updatedRooms));
        setEditingRoomId(null);
        setEditedRoomName("");
    };

    const copyLink = (roomId) => {
        const url = `${window.location.origin}/editor?room=${roomId}`;
        navigator.clipboard.writeText(url).then(() => {
            setCopiedRoomId(roomId);
            setTimeout(() => setCopiedRoomId(null), 2000);
        });
    };

    return (
        <div className="App">
            {!room ? (
                <div>
                    <h1>Welcome to the Collaborative Editor</h1>
                    <button onClick={createRoom}>Create New Room</button>
                    {myRooms.length > 0 && (
                        <div style={{ marginTop: "20px" }}>
                            <h2>Your Rooms</h2>
                            <ul>
                                {myRooms.map((room) => (
                                    <li
                                        key={room.id}
                                        style={{ marginBottom: "6px" }}
                                    >
                                        <a href={`/editor?room=${room.id}`}>
                                            {editingRoomId === room.id ? (
                                                <input
                                                    value={editedRoomName}
                                                    onChange={(e) =>
                                                        setEditedRoomName(
                                                            e.target.value
                                                        )
                                                    }
                                                    onBlur={() =>
                                                        saveRename(room.id)
                                                    }
                                                    onKeyDown={(e) => {
                                                        if (e.key === "Enter")
                                                            saveRename(room.id);
                                                    }}
                                                    autoFocus
                                                />
                                            ) : (
                                                room.name
                                            )}
                                        </a>
                                        <button
                                            onClick={() =>
                                                startRename(room.id, room.name)
                                            }
                                            style={{ marginLeft: "10px" }}
                                        >
                                            Rename
                                        </button>
                                        <button
                                            onClick={() => copyLink(room.id)}
                                            style={{ marginLeft: "6px" }}
                                        >
                                            {copiedRoomId === room.id
                                                ? "Copied"
                                                : "Share"}
                                        </button>
                                        <button
                                            onClick={() => deleteRoom(room.id)}
                                            style={{
                                                marginLeft: "6px",
                                                color: "white",
                                                backgroundColor: "red",
                                                border: "none",
                                                borderRadius: "4px",
                                                padding: "2px 6px",
                                                cursor: "pointer",
                                            }}
                                        >
                                            Delete
                                        </button>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    )}
                </div>
            ) : (
                <div>
                    <h2>{roomName}</h2>
                    <button
                        onClick={goBackToHome}
                        style={{ marginBottom: "10px" }}
                    >
                        ⬅ Back to Home
                    </button>
                    <button
                        onClick={() => copyLink(room)}
                        style={{
                            marginBottom: "10px",
                            marginLeft: "10px",
                            padding: "6px 12px",
                        }}
                    >
                        {copiedRoomId === room
                            ? "Copied"
                            : "🔗 Copy Share Link"}
                    </button>
                    <button
                        onClick={handleSave}
                        disabled={isSaving}
                        style={{
                            marginBottom: "10px",
                            marginLeft: "10px",
                            padding: "6px 12px",
                            backgroundColor: "#4CAF50",
                            color: "white",
                            border: "none",
                            borderRadius: "4px",
                            cursor: isSaving ? "not-allowed" : "pointer",
                            opacity: isSaving ? 0.7 : 1,
                        }}
                    >
                        {isSaving ? "⏳ Saving..." : "💾 Save Document"}
                    </button>
                    <div>
                        {isLoading ? (
                            <div className="loading-indicator">
                                Loading document...
                            </div>
                        ) : (
                            <textarea
                                value={document}
                                onChange={handleChange}
                                rows="20"
                                cols="80"
                                placeholder="Start typing collaboratively..."
                            />
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}

export default App;

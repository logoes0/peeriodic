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
    const params = new URLSearchParams(window.location.search);
    const room = params.get("room");

    useEffect(() => {
        const savedRooms = JSON.parse(localStorage.getItem("myRooms") || "[]");
        setMyRooms(savedRooms);

        if (room) {
            const matched = savedRooms.find((r) => r.id === room);
            if (matched) setRoomName(matched.name);
            else setRoomName("Shared Room");
        }
    }, []);

    useEffect(() => {
        if (!room) return;

        const newSocket = new WebSocket(`ws://localhost:5000/ws?room=${room}`);
        setSocket(newSocket);

        newSocket.onopen = () => {
            console.log("WebSocket connection established");
        };

        newSocket.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                if (message.type === "init" || message.type === "update") {
                    setDocument(message.data);
                }
            } catch (error) {
                console.error("Error parsing message:", error);
            }
        };

        newSocket.onclose = () => {
            console.log("WebSocket connection closed");
        };

        newSocket.onerror = (error) => {
            console.error("WebSocket error:", error);
        };

        // Auto-save every 30 seconds
        const autoSaveInterval = setInterval(() => {
            if (document && document.trim() !== "") {
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
            console.log("Saving content length:", document.length);

            const response = await fetch(
                `http://localhost:5000/api/save?room=${room}`,
                {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        content: document,
                    }),
                }
            );

            console.log("Response status:", response.status);

            if (!response.ok) {
                // Try to get error message from response
                let errorMsg = "Failed to save document";
                try {
                    const errorData = await response.json();
                    errorMsg = errorData.message || errorMsg;
                } catch (e) {
                    console.error("Error parsing error response:", e);
                }
                throw new Error(errorMsg);
            }

            const result = await response.json();
            console.log("Save successful:", result);
            alert(
                `Document saved successfully! Content length: ${result.contentLength}`
            );
        } catch (error) {
            console.error("Full error details:", {
                error: error.toString(),
                stack: error.stack,
            });
            alert(`Error: ${error.message}`);
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

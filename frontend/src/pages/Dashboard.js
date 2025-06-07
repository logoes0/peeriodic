import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { auth } from "../firebase";
import { signOut, onAuthStateChanged } from "firebase/auth";

let socket;

const Dashboard = () => {
  const [user, setUser] = useState(null);
  const [message, setMessage] = useState("");
  const [response, setResponse] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      if (user) {
        setUser(user);
        const token = await user.getIdToken();
        socket = new WebSocket(`ws://localhost:8000/ws?token=${token}`);

        socket.onmessage = (e) => setResponse(e.data);
        socket.onerror = (e) => console.error("WebSocket error:", e);
      } else {
        navigate("/");
      }
    });

    return () => unsubscribe();
  }, [navigate]);

  const handleSend = () => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(message);
      setMessage("");
    }
  };

  const handleLogout = async () => {
    await signOut(auth);
    navigate("/");
  };

  return (
    <div style={{ padding: "2rem", position: "relative" }}>
      <button onClick={handleLogout} style={{ position: "absolute", right: 20, top: 20 }}>
        Sign Out
      </button>
      <h2>Welcome, {user?.email}</h2>
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        placeholder="Type your message"
      />
      <button onClick={handleSend}>Send</button>
      {response && <p>Server: {response}</p>}
    </div>
  );
};

export default Dashboard;
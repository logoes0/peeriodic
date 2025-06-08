import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { auth } from "../firebase";

const TextEditorPage = () => {
  const navigate = useNavigate();
  const [socket, setSocket] = useState(null);
  const [text, setText] = useState("");

  useEffect(() => {
    const setup = async () => {
      const user = auth.currentUser;
      if (!user) {
        navigate("/");
        return;
      }

      const token = await user.getIdToken();
      const ws = new WebSocket(`ws://localhost:8000/ws?token=${token}`);
      setSocket(ws);

      ws.onmessage = (event) => {
        console.log("📩 Server says:", event.data);
      };

      ws.onerror = (e) => {
        console.error("❌ WebSocket error:", e);
      };

      ws.onclose = () => {
        console.log("🛑 WebSocket closed");
      };

      return () => ws.close();
    };

    setup();
  }, [navigate]);

  const handleChange = (e) => {
    const val = e.target.value;
    setText(val);
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(val);
    }
  };

  const handleSignOut = async () => {
    await auth.signOut();
    navigate("/");
  };

  return (
    <div style={{
      height: "100vh",
      backgroundColor: "#536878",
      display: "flex",
      flexDirection: "column"
    }}>
      {/* Top bar with Sign Out */}
      <div style={{
        padding: "10px 20px",
        display: "flex",
        justifyContent: "flex-end",
        alignItems: "center"
      }}>
        <button
          onClick={handleSignOut}
          style={{
            padding: "8px 16px",
            backgroundColor: "#e74c3c",
            color: "white",
            border: "none",
            borderRadius: "4px",
            cursor: "pointer"
          }}
        >
          Sign Out
        </button>
      </div>

      {/* Text Editor */}
      <textarea
        value={text}
        onChange={handleChange}
        placeholder="Start typing here..."
        style={{
          flex: 1,
          margin: "0 20px 20px",
          padding: "20px",
          backgroundColor: "#ffffff",
          fontSize: "16px",
          fontFamily: "monospace",
          border: "none",
          borderRadius: "6px",
          resize: "none",
          outline: "none",
          boxSizing: "border-box"
        }}
      />
    </div>
  );
};

export default TextEditorPage;

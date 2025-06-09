// ✅ TextSession.js (Viewer-Only)
import React, { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";

const TextSession = () => {
  const { sessionId } = useParams();
  const [text, setText] = useState("");
  const socketRef = useRef(null);

  useEffect(() => {
    socketRef.current = new WebSocket(`ws://localhost:8000/ws/${sessionId}`);

    socketRef.current.onmessage = (event) => {
      setText(event.data);
    };

    return () => socketRef.current.close();
  }, [sessionId]);

  return (
    <div
      style={{
        backgroundColor: "#536878",
        height: "100vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        padding: "1rem"
      }}
    >
      <textarea
        value={text}
        readOnly
        style={{
          width: "90%",
          height: "80vh",
          fontSize: "1.2rem",
          backgroundColor: "black",
          color: "white",
          padding: "1rem",
          border: "none",
          borderRadius: "8px",
          resize: "none"
        }}
      />
    </div>
  );
};

export default TextSession;

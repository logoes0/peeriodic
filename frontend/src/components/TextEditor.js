import React, { useState } from "react";

const TextEditor = ({ socket }) => {
  const [text, setText] = useState("");

  const handleChange = (e) => {
    const newText = e.target.value;
    setText(newText);
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(newText); // Send text to backend
    }
  };

  return (
    <textarea
      value={text}
      onChange={handleChange}
      placeholder="Start typing here..."
      style={{
        width: "100%",
        height: "100%",
        border: "none",
        outline: "none",
        padding: "20px",
        fontSize: "16px",
        fontFamily: "monospace",
        resize: "none",
        boxSizing: "border-box"
      }}
    />
  );
};

export default TextEditor;

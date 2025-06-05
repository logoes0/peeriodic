import React from "react";
import { auth, provider, signInWithPopup } from "./firebase";

const SignIn = () => {
  const handleSignIn = async () => {
    try {
      const result = await signInWithPopup(auth, provider);
      const token = await result.user.getIdToken();

      // 🔌 Connect WebSocket with token
      const socket = new WebSocket(`ws://localhost:8000/ws?token=${token}`);

      socket.onopen = () => {
        console.log("🔌 WebSocket connected");
        socket.send("Hello from React!");
      };

      socket.onmessage = (event) => {
        console.log("💬 Server says:", event.data);
      };

      socket.onerror = (err) => {
        console.error("❌ WebSocket error:", err);
      };
    } catch (error) {
      console.error("Firebase sign-in failed:", error);
    }
  };

  return (
    <div style={{ height: "100vh", display: "flex", justifyContent: "center", alignItems: "center" }}>
      <button onClick={handleSignIn} style={{
        padding: "12px 24px",
        fontSize: "16px",
        backgroundColor: "#4285F4",
        color: "white",
        border: "none",
        borderRadius: "6px",
        cursor: "pointer"
      }}>
        Sign In with Firebase
      </button>
    </div>
  );
};

export default SignIn;

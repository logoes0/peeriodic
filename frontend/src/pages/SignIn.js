import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { signInWithEmailAndPassword } from "firebase/auth";
import { auth } from "../firebase";

const SignIn = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const handleSignIn = async () => {
    try {
      await signInWithEmailAndPassword(auth, email, password);
      navigate("/texteditor");
    } catch (err) {
      setError("Login failed: " + err.message);
    }
  };

  return (
    <div style={{
      backgroundColor: "#536878",
      height: "100vh",
      display: "flex",
      justifyContent: "center",
      alignItems: "center",
    }}>
      <div style={{
        display: "flex",
        flexDirection: "column",
        gap: "16px",
        backgroundColor: "#2c2c3e",
        padding: "32px",
        borderRadius: "8px"
      }}>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Email"
          style={{
            padding: "12px",
            fontSize: "16px",
            backgroundColor: "#000",
            color: "#fff",
            border: "1px solid #444",
            borderRadius: "4px"
          }}
        />
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Password"
          style={{
            padding: "12px",
            fontSize: "16px",
            backgroundColor: "#000",
            color: "#fff",
            border: "1px solid #444",
            borderRadius: "4px"
          }}
        />
        <button
          onClick={handleSignIn}
          style={{
            padding: "12px",
            fontSize: "16px",
            backgroundColor: "red",
            color: "#fff",
            border: "none",
            borderRadius: "6px",
            cursor: "pointer"
          }}
        >
          Sign In
        </button>
        <button
          onClick={() => navigate("/signup")}
          style={{
            padding: "12px",
            fontSize: "16px",
            backgroundColor: "#666",
            color: "#fff",
            border: "none",
            borderRadius: "6px",
            cursor: "pointer"
          }}
        >
          Go to Sign Up
        </button>
      </div>
    </div>
  );
};

export default SignIn;

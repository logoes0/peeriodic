import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createUserWithEmailAndPassword } from "firebase/auth";
import { auth } from "../firebase";

const SignUp = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const handleSignUp = async () => {
    try {
      await createUserWithEmailAndPassword(auth, email, password);
      navigate("/texteditor");
    } catch (err) {
      setError("Signup failed: " + err.message);
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
          onClick={handleSignUp}
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
          Sign Up
        </button>
        <button
          onClick={() => navigate("/")}
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
          Back to Sign In
        </button>
      </div>
    </div>
  );
};

export default SignUp;

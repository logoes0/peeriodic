// src/App.js
import React, { useEffect, useState } from "react";
import { onAuthStateChanged } from "firebase/auth";
import { auth } from "./firebase";
import SignIn from "./SignIn";

function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, (user) => {
      setUser(user);
      setLoading(false);
    });

    return () => unsubscribe();
  }, []);

  if (loading) return <div>Loading...</div>;

  if (!user) return <SignIn />;

  return (
    <div style={{ padding: "2rem" }}>
      <h2>Welcome, {user.email}</h2>
      <button
        onClick={() => auth.signOut()}
        style={{ marginTop: "1rem", padding: "10px 20px", cursor: "pointer" }}
      >
        Sign Out
      </button>
    </div>
  );
}

export default App;

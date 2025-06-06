// SignIn.js
import React, { useEffect } from "react";
import { auth } from "./firebase";
import * as firebaseui from "firebaseui";
import "firebaseui/dist/firebaseui.css";

function SignIn() {
  useEffect(() => {
    const ui =
      firebaseui.auth.AuthUI.getInstance() ||
      new firebaseui.auth.AuthUI(auth);

    ui.start("#firebaseui-auth-container", {
      signInFlow: "popup",
      signInOptions: [
        {
          provider: "password",
          requireDisplayName: true,
        },
      ],
      callbacks: {
        signInSuccessWithAuthResult: async (authResult) => {
          const token = await authResult.user.getIdToken();
          const socket = new WebSocket(`ws://localhost:8000/ws?token=${token}`);

          socket.onopen = () => {
            console.log("🔌 WebSocket connected");
            socket.send("Hello from React!");
          };

          socket.onmessage = (event) => {
            console.log("💬 Server says:", event.data);
          };

          return false; // prevent redirect
        },
      },
    });
  }, []);

  return (
    <div
      style={{
        height: "100vh",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <div id="firebaseui-auth-container" />
    </div>
  );
}

export default SignIn;

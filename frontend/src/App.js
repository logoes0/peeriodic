import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import SignIn from "./pages/SignIn";
import SignUp from "./pages/SignUp";
import TextEditor from "./components/TextEditor"; // ✅ Direct import

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<SignIn />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/texteditor" element={<TextEditor />} /> {/* ✅ Editor directly */}
      </Routes>
    </Router>
  );
}

export default App;

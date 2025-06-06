// firebase.js
import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";

const firebaseConfig = {
  apiKey: "AIzaSyAgQ79xJ5-_94lBW1KWkUobH7_Hg0dJSAs",
  authDomain: "peeriodic-auth.firebaseapp.com",
  projectId: "peeriodic-auth",
  storageBucket: "peeriodic-auth.firebasestorage.app",
  messagingSenderId: "606136600211",
  appId: "1:606136600211:web:e72d3b3a42632771c96f66",
};

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);

export { auth };

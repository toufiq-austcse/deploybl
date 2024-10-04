// Import the functions you need from the SDKs you need
import { initializeApp } from 'firebase/app';
import { getAnalytics } from 'firebase/analytics';
// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
  apiKey: 'AIzaSyCV8Clzoh_VSPKkbQqrZH2vM_c62FZoTfM',
  authDomain: 'deploybl.firebaseapp.com',
  projectId: 'deploybl',
  storageBucket: 'deploybl.appspot.com',
  messagingSenderId: '510719812678',
  appId: '1:510719812678:web:6952cea1f087ad47141647',
  measurementId: 'G-51M77H1W7B'
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const analytics = getAnalytics(app);
export default app;
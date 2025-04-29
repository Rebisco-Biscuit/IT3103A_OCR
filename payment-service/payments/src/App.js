// src/App.js
import './custom.css';
import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'; // Import BrowserRouter
import Navbar from './components/Navbar';
import History from './components/History'; 
import CartPage from './components/CartPage';

function App() {
  return (
    <Router> {/* Wrap the application with BrowserRouter */}
      <Navbar />
      <Routes>
        <Route path="/cart" element={<CartPage />} />
        <Route path="/history" element={<History />} />
      </Routes>
    </Router>
  );
}

export default App;

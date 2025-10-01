import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router'
import Landing from './pages/Landing'
import Login from './pages/Login'
import Register from './pages/Register'
import Home from  './pages/Home'

const root = createRoot(document.getElementById('root')!)
root.render(
  <StrictMode>
    <BrowserRouter>
        <Routes>
            <Route path="/" element={<Landing />} /> 
            <Route path="/login" element={<Login />} /> 
            <Route path="/register" element={<Register />} /> 
            <Route path="/home" element={<Home />} /> 
        </Routes>
    </BrowserRouter>
  </StrictMode>,
)

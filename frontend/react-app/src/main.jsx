import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'

// Change from '#root' to '#predictor-root'
ReactDOM.createRoot(document.getElementById('predictor-root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
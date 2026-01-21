import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'

// Placeholder - will be replaced during TDD implementation
function App() {
  return (
    <div className="min-h-screen bg-base-100 flex items-center justify-center">
      <div className="text-center">
        <h1 className="text-2xl font-bold">Music Search Engine</h1>
        <p className="text-base-content/60">Frontend implementation in progress...</p>
      </div>
    </div>
  )
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)

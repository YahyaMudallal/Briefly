
import { UserProvider } from './context/UserContext';
import AuthPage from './components/AuthPage'
import HomePage from './components/HomePage';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import ProfilePage from './components/ProfilePage';
import ProtectedRoute from './components/ProtectedRoute';

// --- MAIN APP COMPONENT ---
// This component sets up the routing and context for the entire application
// It uses React Router for navigation and UserContext for managing user state across the app
// The UserProvider wraps the entire app to provide user state to all components
// The Routes component defines the different pages of the app and their corresponding paths
// The Navigate component is used as a fallback to redirect any unknown paths back to the home page

function App() {
  return(
    <UserProvider>
      <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/auth" element={<AuthPage />} />
        <Route path="/profile" element={
          <ProtectedRoute>
            <ProfilePage />
          </ProtectedRoute>
        } />
        {/* Fallback route */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
      </BrowserRouter>
    </UserProvider>

  )
}

export default App

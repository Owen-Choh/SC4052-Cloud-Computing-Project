import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import ProtectedRoute from "./auth/ProtectedRoute";
import { AuthProvider } from "./auth/useAuth";
import Dashboard from "./pages/Dashboard";
import ConversationPage from "./pages/ConversationPage";
import { ChatbotProvider } from "./context/ChatbotContext";

/**
 * App Component
 *
 * The main component of the SimpleChat application. It sets up the routing
 * and authentication context for the entire application.
 */
function App() {
  return (
    /**
     * BrowserRouter: Enables client-side routing using URLs.
     */
    <Router>
      /**
       * AuthProvider: Provides authentication context to the application,
       * making authentication state and functions available to all components.
       */
      <AuthProvider>
        <Routes>
          {/* Route to redirect from root to /login */}
          <Route path="/" element={<Navigate to="/login" />} />
          {/* Route for the login page */}
          <Route path="/login" element={<LoginPage />} />
          {/* Route for the conversation page, accessible via /chat/:username/:chatbotname */}
          <Route
            path="/chat/:username/:chatbotname"
            element={<ConversationPage />}
          />

          {/* Protected routes that require authentication */}
          <Route element={<ProtectedRoute />}>
            {/* Route for the dashboard, accessible via /Dashboard */}
            <Route
              path="/Dashboard"
              element=
              /**
               * ChatbotProvider: Provides chatbot context to the Dashboard,
               * making chatbot state and functions available.
               */
              {
                <ChatbotProvider>
                  <Dashboard />
                </ChatbotProvider>
              }
            />
          </Route>

          {/* Catch-all route for unknown paths, redirects to /login */}
          <Route path="*" element={<Navigate to="/login" />} />
        </Routes>
      </AuthProvider>
    </Router>
  );
}

export default App;
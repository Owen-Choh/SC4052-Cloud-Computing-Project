// Import necessary modules from the react-router-dom library for handling navigation.
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
// Import the LoginPage component.
import LoginPage from "./pages/LoginPage";
// Import the ProtectedRoute component for handling protected routes.
import ProtectedRoute from "./auth/ProtectedRoute";
// Import the AuthProvider component for providing authentication context.
import { AuthProvider } from "./auth/useAuth";
// Import the Dashboard component.
import Dashboard from "./pages/Dashboard";
// Import the ConversationPage component.
import ConversationPage from "./pages/ConversationPage";
// Import the ChatbotProvider component for managing chatbot context.
import { ChatbotProvider } from "./context/ChatbotContext";

/**
 * The main application component.
 * This component sets up the routing configuration for the application.
 */
function App() {
  return (
    // Wrap the application with the BrowserRouter to enable routing.
    <Router>
      {/* Provide authentication context to the application. */}
      <AuthProvider>
        {/* Define the routes for the application. */}
        <Routes>
          {/* Redirect the root path to the login page. */}
          <Route path="/" element={<Navigate to="/login" />} />
          {/* Define the route for the login page. */}
          <Route path="/login" element={<LoginPage />} />
          {/* Define the route for the conversation page, which displays a chat interface.
           *  It takes the username and chatbotname as parameters.
           */}
          <Route
            path="/chat/:username/:chatbotname"
            element={<ConversationPage />}
          />

          {/* Define a protected route that requires authentication. */}
          <Route element={<ProtectedRoute />}>
            {/* Define the route for the dashboard, which is only accessible to authenticated users. */}
            <Route
              path="/Dashboard"
              element={
                // Provide chatbot context to the Dashboard component.
                <ChatbotProvider>
                  <Dashboard />
                </ChatbotProvider>
              }
            />
          </Route>

          {/* Catch-all route for unknown paths, redirecting to the login page. */}
          <Route path="*" element={<Navigate to="/login" />} />
        </Routes>
      </AuthProvider>
    </Router>
  );
}

export default App;
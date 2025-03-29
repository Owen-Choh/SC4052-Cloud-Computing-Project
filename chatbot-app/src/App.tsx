import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import TempPage from "./pages/TempPage";
import ProtectedRoute from "./auth/ProtectedRoute";
import { AuthProvider } from "./auth/useAuth";
import Dashboard from "./pages/Dashboard";
import ConversationPage from "./pages/ConversationPage";
import { ChatbotProvider } from "./context/ChatbotContext";
import GeminiStreamDisplay from "./pages/GeminiStreamDisplay";

function App() {
  return (
    <Router>
      <AuthProvider>
        <Routes>
          <Route
            path="/test/chat/:username/:chatbotname"
            element={<GeminiStreamDisplay />}
          />

          <Route path="/" element={<Navigate to="/login" />} />
          <Route path="/login" element={<LoginPage />} />
          <Route
            path="/chat/:username/:chatbotname"
            element={<ConversationPage />}
          />

          <Route element={<ProtectedRoute />}>
            <Route path="/TempPage" element={<TempPage />} />
            <Route
              path="/Dashboard"
              element={
                <ChatbotProvider>
                  <Dashboard />
                </ChatbotProvider>
              }
            />
          </Route>
        </Routes>
      </AuthProvider>
    </Router>
  );
}

export default App;

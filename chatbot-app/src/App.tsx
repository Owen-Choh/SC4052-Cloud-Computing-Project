import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import Login from "./pages/Login";
import TempPage from "./pages/TempPage";
import ProtectedRoute from "./auth/ProtectedRoute";
import { AuthProvider } from "./auth/useAuth";
import Dashboard from "./pages/Dashboard";
import ConversationPage from "./pages/ConversationPage";

function App() {  
  return (
    <Router>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<Navigate to='/login' />} />
          <Route path="/login" element={<Login />} />
          <Route path="/chat/:username/:chatbotname" element={<ConversationPage />} />

          <Route element={<ProtectedRoute />}>
            <Route path="/TempPage" element={<TempPage />} />
            <Route path="/Dashboard" element={<Dashboard />} />
          </Route>
        </Routes>
      </AuthProvider>
    </Router>
  );
}

export default App;

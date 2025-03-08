import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Login from "./components/Login";
import TempPage from "./pages/TempPage";
import ProtectedRoute from "./auth/ProtectedRoute";
import { AuthProvider } from "./auth/useAuth";
import Dashboard from "./pages/Dashboard";

function App() {  
  return (
    <Router>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<Login />} />
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

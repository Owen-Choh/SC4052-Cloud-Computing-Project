import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Login from "./components/Login";
import TempPage from "./pages/TempPage";
import ProtectedRoute from "./auth/ProtectedRoute";

function App() {  
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route element={<ProtectedRoute />}>
          <Route path="/TempPage" element={<TempPage />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;

import { useNavigate } from "react-router-dom";
import useAuth from "../auth/useAuth";
import { User } from "../auth/User";
import { useEffect } from "react";
import Sidebar from "../components/Sidebar";


function Dashboard() {
  const { currentUser } = useAuth();
  const navigate = useNavigate();
  
  return (
    <div className="flex h-screen flex-1 w-full">
      <Sidebar />
      <div className="w-full p-4">
        <h1>Dashboard</h1>
        <p>Welcome, {currentUser?.username}!</p>
      </div>
    </div>
  )
}

export default Dashboard;
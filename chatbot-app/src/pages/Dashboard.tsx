import { useNavigate } from "react-router-dom";
import useAuth from "../auth/useAuth";
import { User } from "../auth/User";
import { useEffect } from "react";


function Dashboard() {
  const { currentUser } = useAuth();
  const navigate = useNavigate();
  
  return (
    <div>
      <h1>Dashboard</h1>
      <p>Welcome, {currentUser?.username}!</p>
    </div>
  )
}

export default Dashboard;
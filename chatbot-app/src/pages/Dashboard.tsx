import { useNavigate } from "react-router-dom";
import useAuth from "../auth/useAuth";
import { User } from "../auth/User";
import { useEffect } from "react";
import Sidebar from "../components/Sidebar";
import Botconfigs from "../components/Botconfigs";


function Dashboard() {
  const { currentUser } = useAuth();
  const navigate = useNavigate();
  
  return (
    <div className="flex h-screen flex-1 w-full">
      <Sidebar />
      <div className="w-full">
        <Botconfigs />
      </div>
    </div>
  )
}

export default Dashboard;
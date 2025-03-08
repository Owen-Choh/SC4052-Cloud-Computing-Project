import React from "react";
import useAuth from "../auth/useAuth";
import { Link } from "react-router-dom";


const Sidebar: React.FC = () => {
  const { currentUser } = useAuth();
  const sidebarItems = [
    { name: "Dashboard", path: "/Dashboard" },
    { name: "TempPage", path: "/TempPage" },
  ];

  return (
    <div className="bg-gray-800 text-white w-64 flex flex-col h-screen">
      <h1 className="text-2xl font-bold p-4">Welcome</h1>
      <ul className="flex-grow overflow-y-auto">
        {sidebarItems.map((item) => (
          <li key={item.name} className="p-4 hover:bg-gray-700">
            <Link to={item.path}>{item.name}</Link>
          </li>
        ))}
      </ul>
      <p className="p-4">Logged in as: {currentUser?.username}</p>
    </div>
  )
}

export default Sidebar;
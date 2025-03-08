import React from "react";
import useAuth from "../auth/useAuth";


const Sidebar: React.FC = () => {
  const { currentUser } = useAuth();
  const sidebarItems = [
    { name: "Dashboard", path: "/Dashboard" },
    { name: "TempPage", path: "/TempPage" },
  ];

  return (
    <div className="bg-gray-800 text-white w-64">
      <h1 className="text-2xl font-bold p-4">Sidebar</h1>
      <ul>
        {sidebarItems.map((item) => (
          <li key={item.name} className="p-4 hover:bg-gray-700">
            <a href={item.path}>{item.name}</a>
          </li>
        ))}
      </ul>
      <p className="p-4">Logged in as: {currentUser?.username}</p>
    </div>
  )
}

export default Sidebar;
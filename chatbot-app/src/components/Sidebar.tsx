import React, { useEffect } from "react";
import useAuth from "../auth/useAuth";
import { Link } from "react-router-dom";
import { User } from "../auth/auth";
import { getChatbotsApi } from "../api/apiConfig";
import { Chatbot } from "../api/chatbot";

const Sidebar: React.FC = () => {
  const { currentUser, token } = useAuth();

  const fetchChatbots = async () => {
    const chatbotsResponse: Chatbot[] = await getChatbotsApi.get("",{headers: {Authorization: `Bearer ${token}`}});
    console.log("chatbot response:");
    console.log(chatbotsResponse);
    return chatbotsResponse;
  };

  const sidebarItems = [
    { name: "Dashboard", path: "/Dashboard" },
    { name: "TempPage", path: "/TempPage" },
  ];

  useEffect(() => {
    fetchChatbots();
  }, [currentUser]);

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
  );
};

export default Sidebar;

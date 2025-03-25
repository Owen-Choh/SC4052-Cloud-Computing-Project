import React from "react";
import { useChatbotContext } from "../context/ChatbotContext";
import LogoutIcon from "@mui/icons-material/Logout";
import useAuth from "../auth/useAuth";
interface SidebarProps {
  currentUsername: string | undefined;
}

const Sidebar: React.FC<SidebarProps> = ({ currentUsername }) => {
  const { chatbots, setSelectedChatbot, setIsCreatingChatbot } =
    useChatbotContext();

  const { doLogout } = useAuth();

  return (
    <div className="bg-gray-800 text-white w-64 flex flex-col h-screen">
      <button
        onClick={() => {
          setIsCreatingChatbot(true);
          setSelectedChatbot(null);
        }}
        className="m-4 bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
      >
        + New Chatbot
      </button>
      <ul className="flex-grow overflow-y-auto">
        {chatbots.map((chatbot) => (
          <li
            key={chatbot.chatbotid}
            onClick={() => {
              setSelectedChatbot(chatbot);
              setIsCreatingChatbot(false);
            }}
            className="m-2 rounded-lg p-4 hover:bg-gray-700 cursor-pointer"
          >
            {chatbot.chatbotname}
          </li>
        ))}
      </ul>
      <div className="flex justify-center items-center p-4">
        <button onClick={doLogout} className="!p-1 m-0 flex">
          <LogoutIcon />
        </button>
        <p className="p-4">Logged in as: {currentUsername}</p>
      </div>
    </div>
  );
};

export default Sidebar;

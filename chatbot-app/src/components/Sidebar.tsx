import React from "react";
import { useChatbotContext } from "../context/ChatbotContext";
interface SidebarProps {
  currentUsername: string | undefined;
}

const Sidebar: React.FC<SidebarProps> = ({ currentUsername }) => {
  const { chatbots, setSelectedChatbot, setIsCreatingChatbot } = useChatbotContext();

  return (
    <div className="bg-gray-800 text-white w-64 flex flex-col h-screen">
      <h1 className="text-2xl font-bold p-4">Welcome</h1>
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
            className="p-4 hover:bg-gray-700 cursor-pointer"
          >
            {chatbot.chatbotname}
          </li>
        ))}
      </ul>

      <p className="p-4">Logged in as: {currentUsername}</p>
    </div>
  );
};

export default Sidebar;
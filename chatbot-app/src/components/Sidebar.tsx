import React from "react";
import { Chatbot } from "../api/chatbot";

interface SidebarProps {
  currentUsername: string | undefined;
  chatbots: Chatbot[];
  onCreateNewChatbot: () => void;
  selectChatbot: (chatbotid: number | null) => void;
}

const Sidebar: React.FC<SidebarProps> = ({
  currentUsername,
  chatbots,
  onCreateNewChatbot,
  selectChatbot,
}) => {
  return (
    <div className="bg-gray-800 text-white w-64 flex flex-col h-screen">
      <h1 className="text-2xl font-bold p-4">Welcome</h1>
      <ul className="flex-grow overflow-y-auto">
        {/* New Chatbot Button */}
        <button
          onClick={onCreateNewChatbot}
          className="m-4 bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          + New Chatbot
        </button>
        {/* <li className="p-4 hover:bg-gray-700">
          <Link to="/Dashboard">Dashboard</Link>
        </li>
        <li className="p-4 hover:bg-gray-700">
          <Link to="/TempPage">TempPage</Link>
        </li> */}

        {/* Dynamically list chatbots */}
        {chatbots.map((chatbot) => (
          <li
            key={chatbot.chatbotid}
            onClick={() => selectChatbot(chatbot.chatbotid)}
            className="p-4 hover:bg-gray-700"
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

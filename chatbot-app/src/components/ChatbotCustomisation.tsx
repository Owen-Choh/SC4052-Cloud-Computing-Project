import React, { useState } from "react";

interface ChatbotCustomisationProps {
  chatbotName: string;
  isShared: boolean;
  chatbotLink: string;
}

const ChatbotCustomisation: React.FC<ChatbotCustomisationProps> = ({ chatbotName, isShared, chatbotLink }) => {
  const [currentName, setCurrentName] = useState(chatbotName);
  const [currentShared, setCurrentShared] = useState(isShared);

  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-2xl font-bold">This page is to customise your Chatbot's behaviour</h2>
      <div className="flex flex-row gap-4 items-center">
        <p className="text-lg">Your Chatbot's Name: </p>
        <input 
          type="text"
          value={currentName}
          onChange={(e) => setCurrentName(e.target.value)}
          className="p-2 border rounded bg-gray-900 text-white"
        />
      </div>
      <div className="flex flex-row items-center gap-4">
        <p className="text-lg">Do we share your chatbot?</p>
        <input
          type="checkbox"
          checked={currentShared}
          onChange={(e) => setCurrentShared(e.target.checked)}
          className="w-5 h-5"
        />
      </div>

      <div className="flex flex-row items-center gap-4">
        <p className="text-lg">Your Chatbot's unique link:</p>
        <div>{chatbotLink}</div>
      </div>

      <button 
        className="bg-blue-600 p-2 rounded mt-4 hover:bg-blue-700"
        onClick={() => alert("Settings button isnt working!")}
      >
        Save Changes
      </button>
    </div>
  )
};

export default ChatbotCustomisation;
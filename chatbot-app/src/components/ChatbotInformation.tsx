import React, { useEffect, useState } from "react";
import { Chatbot } from "../api/chatbot";
import Botconfigs from "./Botconfigs";

interface ChatbotInformationProps {
  chatbotName: string;
  isShared: boolean;
  chatbotEndpoint: string;
  updateChatbotLink: (chatbotName: string) => void;
  updateChatbotInfo: (chatbotName: string, isShared: boolean) => void;
  saveChatbot: () => void;
}

const ChatbotInformation: React.FC<ChatbotInformationProps> = ({
  chatbotName,
  isShared,
  chatbotEndpoint,
  updateChatbotLink,
  updateChatbotInfo,
  saveChatbot,
}) => {
  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-2xl font-bold">Chatbot Information and Settings</h2>
      <div className="flex flex-row gap-4 items-center">
        <p className="text-lg">Your Chatbot's Name: </p>
        <input
          type="text"
          value={chatbotName}
          onChange={(e) => {
            updateChatbotInfo(e.target.value, isShared);
            updateChatbotLink(e.target.value);
          }}
          className="p-2 border rounded bg-gray-900 text-white"
        />
      </div>
      <div className="flex flex-row items-center gap-4">
        <p className="text-lg">Do we share your chatbot?</p>
        <input
          type="checkbox"
          checked={isShared}
          onChange={(e) => {
            updateChatbotInfo(chatbotName, e.target.checked);
          }}
          className="w-5 h-5"
        />
      </div>

      <div className="flex flex-row items-center gap-4">
        <p className="text-lg">Your Chatbot's unique link:</p>
        <div>{window.location.origin + chatbotEndpoint}</div>
      </div>

      <button
        className="bg-blue-600 p-2 rounded mt-4 hover:bg-blue-700"
        onClick={() => saveChatbot()}
      >
        Save Changes
      </button>
    </div>
  );
};

export default ChatbotInformation;

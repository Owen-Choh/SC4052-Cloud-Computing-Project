import React, { useEffect, useState } from "react";
import { Chatbot } from "../api/chatbot";
import Botconfigs from "./Botconfigs";

interface ChatbotInformationProps {
  chatbotName: string;
  isShared: boolean;
  chatbotEndpoint: string;
  description: string;
  createdDate: string;
  updatedDate: string;
  lastUsed: string;
  updateChatbotLink: (chatbotName: string) => void;
  updateChatbotInfo: (
    chatbotName: string,
    isShared: boolean,
    description: string
  ) => void;
}

const ChatbotInformation: React.FC<ChatbotInformationProps> = ({
  chatbotName,
  isShared,
  chatbotEndpoint,
  description,
  createdDate,
  updatedDate,
  lastUsed,
  updateChatbotLink,
  updateChatbotInfo,
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
            updateChatbotInfo(e.target.value, isShared, description);
            updateChatbotLink(e.target.value);
          }}
          className="flex-grow p-2 border rounded bg-gray-900 text-white"
        />
      </div>
      <div className="flex flex-row gap-4 items-center">
        <p className="text-lg">Description of chatbot: </p>
        <textarea
          value={description}
          onChange={(e) => {
            updateChatbotInfo(chatbotName, isShared, e.target.value);
          }}
          className="flex-grow p-2 border rounded bg-gray-900 text-white"
        />
      </div>
      <div className="flex flex-row items-center gap-4">
        <p className="text-lg">Do we share your chatbot?</p>
        <input
          type="checkbox"
          checked={isShared}
          onChange={(e) => {
            updateChatbotInfo(chatbotName, e.target.checked, description);
          }}
          className="w-5 h-5"
        />
      </div>
      <div className="flex flex-row items-center gap-4">
        <p className="text-lg">
          Your Chatbot's unique link (has to be shared in order to access):
        </p>
        <div>{window.location.origin + chatbotEndpoint}</div>
      </div>
      {createdDate && (
        <div className="flex flex-row items-center gap-4">
          <p className="text-lg">Your Chatbot's created on: </p>
          <p className="text-lg">{createdDate}</p>
        </div>
      )}
      {updatedDate && (
        <div className="flex flex-row items-center gap-4">
          <p className="text-lg">Your Chatbot's last updated on: </p>
          <p className="text-lg">{updatedDate}</p>
        </div>
      )}
      {lastUsed && (
        <div className="flex flex-row items-center gap-4">
          <p className="text-lg">Your Chatbot's last used on: </p>
          <p className="text-lg">{lastUsed}</p>
          <p className="text-lg">(not real time updates)</p>
        </div>
      )}
    </div>
  );
};

export default ChatbotInformation;

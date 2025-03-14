import React, { useState, useEffect } from "react";
import FileUpload from "./FileUpload";

interface ChatbotCustomisationProps {
  chatbotBehaviour: string;
  chatbotContext: string;
  chatbotDocument: string;
  updateChatbotCustomisation: (behaviour: string, context: string) => void;
  updateChatbotFile: (document: File | null) => void;
  saveChatbotCustomisation: () => void;
}

const ChatbotCustomisation: React.FC<ChatbotCustomisationProps> = ({
  chatbotBehaviour,
  chatbotContext,
  chatbotDocument,
  updateChatbotCustomisation,
  updateChatbotFile,
  saveChatbotCustomisation,
}) => {
  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-2xl font-bold">
        This page is to customise your Chatbot's behaviour
      </h2>
      <div className="flex flex-col gap-4">
        <p className="text-xl">How should your Chatbot behave?</p>
        <textarea
          name="behaviour_input"
          rows={2}
          value={chatbotBehaviour}
          onChange={(e) => updateChatbotCustomisation(e.target.value, chatbotContext)}
          className="p-2 border rounded bg-gray-900 text-white"
        />
      </div>
      <div className="flex flex-col gap-4">
        <p className="text-xl">
          Any specific information that your chatbot should know about?
        </p>
        <textarea
          name="context_input"
          rows={5}
          value={chatbotContext}
          onChange={(e) => updateChatbotCustomisation(chatbotBehaviour, e.target.value)}
          className="p-2 border rounded bg-gray-900 text-white"
        />
      </div>
      <div className="flex flex-col gap-4">
        <p className="text-xl">Any documents that your chatbot should use?</p>
        {chatbotDocument ? (
          <p>
            Chatbot currently has:{" "}
            <span className="font-bold">{chatbotDocument}</span>
          </p>
        ) : (
          <p>No document uploaded</p>
        )}
        <FileUpload onFileSelect={updateChatbotFile} />
      </div>
      <button
        className="bg-blue-600 p-2 rounded mt-4 hover:bg-blue-700"
        onClick={() => saveChatbotCustomisation()}
      >
        Save Changes
      </button>
    </div>
  );
};

export default ChatbotCustomisation;

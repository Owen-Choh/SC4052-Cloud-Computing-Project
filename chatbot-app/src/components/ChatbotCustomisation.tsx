import React, { useState, useEffect } from "react";
import FileUpload from "./FileUpload";

interface ChatbotCustomisationProps {
  chatbotBehaviour: string;
  chatbotContext: string;
  chatbotDocument: string;
  updateChatbotCustomisation: (behaviour: string, context: string) => void;
  updateChatbotFile: (document: File | null) => void;
  excludeFile: boolean;
  toggleExcludeFile: () => void;
}

const ChatbotCustomisation: React.FC<ChatbotCustomisationProps> = ({
  chatbotBehaviour,
  chatbotContext,
  chatbotDocument,
  updateChatbotCustomisation,
  updateChatbotFile,
  excludeFile,
  toggleExcludeFile,
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
          onChange={(e) =>
            updateChatbotCustomisation(e.target.value, chatbotContext)
          }
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
          onChange={(e) =>
            updateChatbotCustomisation(chatbotBehaviour, e.target.value)
          }
          className="p-2 border rounded bg-gray-900 text-white"
        />
      </div>
      <div className="flex flex-col gap-4">
        <p className="text-xl">
          Any documents that your chatbot should use? (There can only be one at
          any time)
        </p>
        {chatbotDocument ? (
          <div className="flex flex-row gap-2 items-center">
            <p>
              Chatbot currently has:{" "}
              <span className="font-bold">{chatbotDocument}</span>
            </p>
            <button
              className={`rounded ${
                excludeFile
                  ? "bg-red-600 hover:bg-red-700"
                  : "bg-green-600 hover:bg-green-700"
              }`}
              onClick={toggleExcludeFile}
            >
              {excludeFile ? "Exclude Original File" : "Include Original File"}
            </button>
          </div>
        ) : (
          <p>No document uploaded</p>
        )}
        <FileUpload onFileSelect={updateChatbotFile} />
      </div>
    </div>
  );
};

export default ChatbotCustomisation;

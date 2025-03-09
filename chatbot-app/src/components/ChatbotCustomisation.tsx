import React, { useState } from "react";
import FileUpload from "./FileUpload";

interface ChatbotCustomisationProps {
  chatbotBehaviour: string;
  chatbotContext: string;
  chatbotDocument: string;
}

const ChatbotCustomisation: React.FC<ChatbotCustomisationProps> = ({
  chatbotBehaviour,
  chatbotContext,
  chatbotDocument,
}) => {
  const [currentBehaviour, setCurrentBehaviour] = useState(chatbotBehaviour);
  const [currentContext, setCurrentContext] = useState(chatbotContext);
  const [currentDocument, setCurrentDocument] = useState<File | null>(null);

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
          value={currentBehaviour}
          onChange={(e) => setCurrentBehaviour(e.target.value)}
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
          value={currentContext}
          onChange={(e) => setCurrentContext(e.target.value)}
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
        <FileUpload onFileSelect={setCurrentDocument} />
      </div>
      <button
        className="bg-blue-600 p-2 rounded mt-4 hover:bg-blue-700"
        onClick={() =>
          alert(
            "Settings button isnt working! " +
              "Selected file for upload: " +
              currentDocument?.name
          )
        }
      >
        Save Changes
      </button>
    </div>
  );
};

export default ChatbotCustomisation;

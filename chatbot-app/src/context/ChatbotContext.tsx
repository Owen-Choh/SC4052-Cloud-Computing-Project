import React, { createContext, useContext, useState, ReactNode } from "react";
import { Chatbot } from "../api/chatbot";

interface ChatbotContextType {
  chatbots: Chatbot[];
  setChatbots: React.Dispatch<React.SetStateAction<Chatbot[]>>;
  selectedChatbot: Chatbot | null;
  setSelectedChatbot: React.Dispatch<React.SetStateAction<Chatbot | null>>;
  isCreatingChatbot: boolean;
  setIsCreatingChatbot: React.Dispatch<React.SetStateAction<boolean>>;
}

const ChatbotContext = createContext<ChatbotContextType | undefined>(undefined);

export const useChatbotContext = () => {
  const context = useContext(ChatbotContext);
  if (!context) {
    throw new Error("useChatbotContext must be used within ChatbotProvider");
  }
  return context;
};

export const ChatbotProvider = ({ children }: { children: ReactNode }) => {
  const [chatbots, setChatbots] = useState<Chatbot[]>([]);
  const [selectedChatbot, setSelectedChatbot] = useState<Chatbot | null>(null);
  const [isCreatingChatbot, setIsCreatingChatbot] = useState(false);

  return (
    <ChatbotContext.Provider
      value={{
        chatbots,
        setChatbots,
        selectedChatbot,
        setSelectedChatbot,
        isCreatingChatbot,
        setIsCreatingChatbot,
      }}
    >
      {children}
    </ChatbotContext.Provider>
  );
};

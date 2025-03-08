import React from "react";

interface ChatbotInformationProps {
  chatbotName: string;
  isShared: boolean;
  chatbotLink: string;
}

const ChatbotInformation: React.FC<ChatbotInformationProps> = ({ chatbotName, isShared, chatbotLink }) => {
  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-2xl font-bold">Chatbot Information and Settings</h2>
      <div className="flex flex-row gap-4">
        <p>Your Chatbot's Name: </p>
        <div>{ chatbotName }</div>
      </div>
      <div className="flex flex-row gap-4">
        <p>Do we share your chatbot? </p>
        <div>{ isShared }</div>
      </div>
      <div className="flex flex-row gap-4">
        <p>Your Chatbot's unique link: </p>
        <div>{ chatbotLink }</div>
      </div>
    </div>
  )
};

export default ChatbotInformation;
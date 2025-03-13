import React from "react";
import Tab from "./ui/Tab";
import TabPanel from "./ui/TabPanel";
import ChatbotInformation from "./ChatbotInformation";
import ChatbotCustomisation from "./ChatbotCustomisation";
import { Chatbot } from "../api/chatbot";

interface BotconfigsProps {
  username: string;
  chatbot: Chatbot;
}

const Botconfigs: React.FC<BotconfigsProps> = ({ username, chatbot }) => {
  const [activeTab, setActiveTab] = React.useState("chatInfo");
  const [chatbotLink, setChatbotEndpoint] = React.useState("/chat/" + username + "/" + chatbot.Chatbotname);

  const updateChatbotLink = (chatbotName: string) => {
    setChatbotEndpoint("/chat/" + username + "/" + chatbotName);
  };

  const chatbotDocument = "project-details.pdf";

  return (
    <div className="flex flex-col w-full h-full p-4 bg-gray-900 gap-4">
      <div className="flex gap-4">
        <Tab
          label="Chatbot information"
          isActive={activeTab === "chatInfo"}
          onClick={() => setActiveTab("chatInfo")}
        />
        <Tab
          label="Customise"
          isActive={activeTab === "customisation"}
          onClick={() => setActiveTab("customisation")}
        />
      </div>

      <div className="border-b-2 border-gray-700"></div>

      <div className="w-full flex-grow overflow-y-auto">
        <TabPanel activeTab={activeTab} tabKey="chatInfo">
          <ChatbotInformation
            chatbotName={chatbot.Chatbotname}
            isShared={chatbot.IsShared}
            chatbotEndpoint={chatbotLink}
            updateChatbotLink={(chatbotName) => updateChatbotLink(chatbotName)}
          />
        </TabPanel>
        <TabPanel activeTab={activeTab} tabKey="customisation">
          <ChatbotCustomisation
            chatbotBehaviour={chatbot.Behaviour}
            chatbotContext={chatbot.Usercontext}
            chatbotDocument={chatbot.Filepath}
          />
        </TabPanel>
      </div>
    </div>
  );
};

export default Botconfigs;

import React from "react";
import Tab from "./ui/Tab";
import TabPanel from "./ui/TabPanel";
import ChatbotInformation from "./ChatbotInformation";
import ChatbotCustomisation from "./ChatbotCustomisation";
import { Chatbot } from "../api/chatbot";
import { createChatbotsApi } from "../api/apiConfig";
import useAuth from "../auth/useAuth";

interface BotconfigsProps {
  username: string;
  chatbot: Chatbot;
}

const Botconfigs: React.FC<BotconfigsProps> = ({ username, chatbot }) => {
  const { token } = useAuth();
  const [activeTab, setActiveTab] = React.useState("chatInfo");
  const [chatbotLink, setChatbotEndpoint] = React.useState("/chat/" + username + "/" + chatbot.Chatbotname);
  const [currentChatbot, setCurrentChatbot] = React.useState(chatbot);

  const updateChatbotLink = (chatbotName: string) => {
    setChatbotEndpoint("/chat/" + username + "/" + chatbotName);
  };

  const updateChatbotInfo = (chatbotName: string, isShared: boolean) => {
    setCurrentChatbot({
      ...currentChatbot,
      Chatbotname: chatbotName,
      IsShared: isShared,
    });
    console.log("Chatbot updated: ", currentChatbot);
  };

  const updateChatbotCustomisation = (behaviour: string, context: string, document: File | null) => {
    setCurrentChatbot({
      ...currentChatbot,
      Behaviour: behaviour,
      Usercontext: context,
      Filepath: document ? document.name : currentChatbot.Filepath,
    });
    console.log("Chatbot customisation updated: ", currentChatbot);
  };

  const saveChatbot = () => {
    // Save the chatbot to the database
    console.log("Chatbot saved: ", currentChatbot);
    // alert("button not ready");
    createChatbotsApi.post("", currentChatbot, {
      headers: { Authorization: `Bearer ${token}` },
    });
  };

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
            updateChatbotInfo={(chatbotName, isShared) => updateChatbotInfo(chatbotName, isShared)}
            saveChatbot={() => saveChatbot()}
          />
        </TabPanel>
        <TabPanel activeTab={activeTab} tabKey="customisation">
          <ChatbotCustomisation
            chatbotBehaviour={chatbot.Behaviour}
            chatbotContext={chatbot.Usercontext}
            chatbotDocument={chatbot.Filepath}
            updateChatbotCustomisation={(behaviour, context, document) => updateChatbotCustomisation(behaviour, context, document)}
            saveChatbotCustomisation={() => saveChatbot()}
          />
        </TabPanel>
      </div>
    </div>
  );
};

export default Botconfigs;

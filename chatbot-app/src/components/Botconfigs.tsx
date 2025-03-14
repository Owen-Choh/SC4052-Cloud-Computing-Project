import React from "react";
import Tab from "./ui/Tab";
import TabPanel from "./ui/TabPanel";
import ChatbotInformation from "./ChatbotInformation";
import ChatbotCustomisation from "./ChatbotCustomisation";
import { Chatbot, CreateChatbotPayload } from "../api/chatbot";
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
      File: document ? document : currentChatbot.File,
    });
    console.log("Chatbot customisation updated: ", currentChatbot);
  };

  const saveChatbot = () => {
    // Save the chatbot to the database
    console.log("Chatbot saved: ", currentChatbot);
    const formData = new FormData();
    formData.append("Chatbotname", currentChatbot.Chatbotname);
    formData.append("Behaviour", currentChatbot.Behaviour);
    formData.append("Usercontext", currentChatbot.Usercontext);
    formData.append("IsShared", currentChatbot.IsShared.toString());
    
    if (currentChatbot.File) {
      formData.append("File", currentChatbot.File); // Append file if available
    }
  
    createChatbotsApi.post("", formData, {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "multipart/form-data", // Important for file upload
      },
    });

    // const requestChatbot: CreateChatbotPayload ={
    //   Chatbotname: currentChatbot.Chatbotname,
    //   Behaviour: currentChatbot.Behaviour,
    //   Usercontext: currentChatbot.Usercontext,
    //   IsShared: currentChatbot.IsShared,
    //   File: null,
    // }
    // createChatbotsApi.post("", requestChatbot, {
    //   headers: { Authorization: `Bearer ${token}` },
    // });
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

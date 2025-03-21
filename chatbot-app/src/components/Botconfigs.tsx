import React, { useEffect } from "react";
import Tab from "./ui/Tab";
import TabPanel from "./ui/TabPanel";
import ChatbotInformation from "./ChatbotInformation";
import ChatbotCustomisation from "./ChatbotCustomisation";
import { Chatbot } from "../api/chatbot";
import { chatbotsApi } from "../api/apiConfig";
import useAuth from "../auth/useAuth";

interface BotconfigsProps {
  username: string;
  chatbot: Chatbot;
}

const Botconfigs: React.FC<BotconfigsProps> = ({ username, chatbot }) => {
  const { token } = useAuth();
  const [activeTab, setActiveTab] = React.useState("chatInfo");
  const [chatbotLink, setChatbotEndpoint] = React.useState("/chat/" + username + "/" + chatbot.chatbotname);
  const [currentChatbot, setCurrentChatbot] = React.useState(chatbot);
  const [success, setSuccess] = React.useState("");
  const [error, setError] = React.useState("");
  const [originalFilepath, setOriginalFilepath] = React.useState(chatbot.filepath);

  const updateChatbotLink = (chatbotName: string) => {
    setChatbotEndpoint("/chat/" + username + "/" + chatbotName);
  };

  const updateChatbotInfo = (chatbotName: string, isShared: boolean) => {
    setCurrentChatbot({
      ...currentChatbot,
      chatbotname: chatbotName,
      isShared: isShared,
    });
    console.log("Chatbot updated: ", currentChatbot);
  };

  const updateChatbotCustomisation = (behaviour: string, context: string) => {
    setCurrentChatbot({
      ...currentChatbot,
      behaviour: behaviour,
      usercontext: context,
    });
    console.log("Chatbot customisation updated: ", currentChatbot);
  };

  const updateChatbotFile = (document: File | null) => {
    setCurrentChatbot({
      ...currentChatbot,
      filepath: document ? document.name : currentChatbot.filepath,
      file: document ? document : currentChatbot.file,
    });
  }

  const saveChatbot = async () => {
    // Save the chatbot to the database
    console.log("Chatbot saved: ", currentChatbot);
    const formData = new FormData();
    formData.append("chatbotname", currentChatbot.chatbotname);
    formData.append("description", currentChatbot.description);
    formData.append("behaviour", currentChatbot.behaviour);
    formData.append("usercontext", currentChatbot.usercontext);
    formData.append("isShared", currentChatbot.isShared.toString());
    
    if (currentChatbot.file) {
      formData.append("file", currentChatbot.file); // Append file if available
    }

    try {
      if(currentChatbot.chatbotid == null) {
        const response = await chatbotsApi.post("", formData, {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "multipart/form-data",
          },
        });
      } else {
        const response = await chatbotsApi.put("", formData, {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "multipart/form-data",
          },
        });
      }
  
      console.log("Chatbot saved successfully:", response.data);
      setSuccess("Chatbot saved successfully!");
      setError("");
    } catch (err: any) {
      setSuccess("");
      console.error("Failed to save chatbot:", err);
      if (err.response?.status) {
        setError("Failed to save chatbot. Error: "+err.response?.status + " " + err.response?.data?.error);
      } else {
        setError("Failed to save chatbot. Unknown Error occured");
      }
    }
  };
  
  useEffect(() => {
    setCurrentChatbot(chatbot);
    updateChatbotLink(chatbot.chatbotname);
    setSuccess("");
    setError("");
    setOriginalFilepath(chatbot.filepath);
    console.log("Chatbot updated: ", currentChatbot);
  }, [chatbot]);

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
        <Tab
          label={chatbot.chatbotname}
          isActive={activeTab === "customisation"}
          onClick={() => setActiveTab("customisation")}
        />
      </div>

      <div className="border-b-2 border-gray-700"></div>

      <div className="w-full flex-grow overflow-y-auto">
        <TabPanel activeTab={activeTab} tabKey="chatInfo">
          <ChatbotInformation
            chatbotName={currentChatbot.chatbotname}
            isShared={currentChatbot.isShared}
            chatbotEndpoint={chatbotLink}
            updateChatbotLink={(chatbotName) => updateChatbotLink(chatbotName)}
            updateChatbotInfo={(chatbotName, isShared) => updateChatbotInfo(chatbotName, isShared)}
            saveChatbot={() => saveChatbot()}
          />
        </TabPanel>
        <TabPanel activeTab={activeTab} tabKey="customisation">
          <ChatbotCustomisation
            chatbotBehaviour={currentChatbot.behaviour}
            chatbotContext={currentChatbot.usercontext}
            chatbotDocument={originalFilepath}
            updateChatbotCustomisation={(behaviour, context) => updateChatbotCustomisation(behaviour, context)}
            updateChatbotFile={(document) => updateChatbotFile(document)}
            saveChatbotCustomisation={() => saveChatbot()}
          />
        </TabPanel>
        {success && <p className="p-4 text-green-500">{success}</p>}
        {error && <p className="p-4 text-red-500">{error}</p>}
      </div>
    </div>
  );
};

export default Botconfigs;

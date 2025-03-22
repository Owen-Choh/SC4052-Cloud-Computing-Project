import React, { useEffect, useState } from "react";
import Tab from "./ui/Tab";
import TabPanel from "./ui/TabPanel";
import ChatbotInformation from "./ChatbotInformation";
import ChatbotCustomisation from "./ChatbotCustomisation";
import { Chatbot } from "../api/chatbot";
import { chatbotsApi } from "../api/apiConfig";
import useAuth from "../auth/useAuth";
import { useChatbotContext } from "../context/ChatbotContext";

interface BotconfigsProps {
  username: string;
  chatbot: Chatbot;
  setChatbot: React.Dispatch<React.SetStateAction<Chatbot | null>>;
  excludeFile: boolean;
  setExcludeFile: React.Dispatch<React.SetStateAction<boolean>>;
}

const Botconfigs: React.FC<BotconfigsProps> = ({
  username,
  chatbot,
  setChatbot,
  excludeFile,
  setExcludeFile,
}) => {
  const { token } = useAuth();
  const {
    isCreatingChatbot,
    setIsCreatingChatbot,
    updateChatbotInContext,
  } = useChatbotContext();
  const [activeTab, setActiveTab] = useState("chatInfo");
  const [chatbotLink, setChatbotLink] = useState(`/chat/${username}/${chatbot.chatbotname}`);
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");
  
  const updateChatbotInfo = (chatbotName: string, isShared: boolean, description: string) => {
    setChatbot(prev => prev ? { ...prev, chatbotname: chatbotName, isShared, description } : prev);
    setChatbotLink(`/chat/${username}/${chatbotName}`);
  };

  const updateChatbotCustomisation = (behaviour: string, context: string) => {
    setChatbot(prev => prev ? { ...prev, behaviour, usercontext: context } : prev);
  };

  const updateChatbotFile = (document: File | null) => {
    setChatbot(prev => prev ? { ...prev, filepath: document ? document.name : prev.filepath, file: document } : prev);
  };

  const saveChatbot = async () => {
    if (!chatbot) return;
    const formData = new FormData();
    formData.append("chatbotname", chatbot.chatbotname);
    formData.append("description", chatbot.description);
    formData.append("behaviour", chatbot.behaviour);
    formData.append("usercontext", chatbot.usercontext);
    formData.append("isShared", chatbot.isShared.toString());
    if (chatbot.file) formData.append("file", chatbot.file);
    if (excludeFile) formData.append("removeFile", "true");

    try {
      const response = !isCreatingChatbot
        ? await chatbotsApi.put(`/${chatbot.chatbotid}`, formData, {
            headers: { Authorization: `Bearer ${token}`, "Content-Type": "multipart/form-data" },
          })
        : await chatbotsApi.post("/", formData, {
            headers: { Authorization: `Bearer ${token}`, "Content-Type": "multipart/form-data" },
          });

      console.log("Chatbot saved successfully:", response.data);
      setSuccess("Chatbot saved successfully!");
      setError("");
      if (isCreatingChatbot) {
        // Update chatbot id if user create new chatbot
        setChatbot(response.data.chatbotid);
        setIsCreatingChatbot(false);
      }
      updateChatbotInContext(chatbot);
    } catch (err: any) {
      console.error("Failed to save chatbot:", err);
      setSuccess("");
      setError("Failed to save chatbot. " + (err.response?.data?.error || "Unknown error"));
    }
  };

  useEffect(()=>{
    setSuccess("");
    setError("");
  }, [chatbot]);

  return (
    <div className="flex flex-col w-full h-full p-4 bg-gray-900 gap-4">
      <div className="flex gap-4">
        <Tab label="Chatbot information" isActive={activeTab === "chatInfo"} onClick={() => setActiveTab("chatInfo")} />
        <Tab label="Customise" isActive={activeTab === "customisation"} onClick={() => setActiveTab("customisation")} />
      </div>

      <div className="border-b-2 border-gray-700"></div>

      <div className="w-full flex-grow overflow-y-auto">
        <TabPanel activeTab={activeTab} tabKey="chatInfo">
          <ChatbotInformation
            chatbotName={chatbot.chatbotname}
            isShared={chatbot.isShared}
            chatbotEndpoint={chatbotLink}
            description={chatbot.description}
            updateChatbotLink={setChatbotLink}
            updateChatbotInfo={updateChatbotInfo}
          />
        </TabPanel>
        <TabPanel activeTab={activeTab} tabKey="customisation">
          <ChatbotCustomisation
            chatbotBehaviour={chatbot.behaviour}
            chatbotContext={chatbot.usercontext}
            chatbotDocument={chatbot.filepath}
            excludeFile={excludeFile}
            toggleExcludeFile={() => setExcludeFile(prev => !prev)}
            updateChatbotCustomisation={updateChatbotCustomisation}
            updateChatbotFile={updateChatbotFile}
          />
        </TabPanel>
        {success && <p className="p-4 text-green-500">{success}</p>}
        {error && <p className="p-4 text-red-500">{error}</p>}
      </div>
      <div className="border-b-2 border-gray-700"></div>
      <button
        className="bg-green-600 p-2 rounded hover:bg-green-700"
        onClick={saveChatbot}
      >
        Save Changes
      </button>
      </div>
  );
};

export default Botconfigs;

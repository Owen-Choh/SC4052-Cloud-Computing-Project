import { useNavigate } from "react-router-dom";
import useAuth from "../auth/useAuth";
import { useEffect, useState } from "react";
import Sidebar from "../components/Sidebar";
import Botconfigs from "../components/Botconfigs";

import { getChatbotsApi } from "../api/apiConfig";
import { Chatbot } from "../api/chatbot";

function Dashboard() {
  const { currentUser, token } = useAuth();
  const [selectedChatbot, setSelectedChatbot] = useState<Chatbot | null>(null);
  const [selectedChatbotID, setSelectedChatbotID] = useState<number | null>(
    null
  );
  const [isCreatingChatbot, setIsCreatingChatbot] = useState(false);
  const [chatbots, setChatbots] = useState<Chatbot[]>([]);

  const newBot: Chatbot = {
    Chatbotid: 0,
    Userid: currentUser?.userid ? currentUser.userid : 0,
    Chatbotname: "",
    Behaviour: "",
    Usercontext: "",
    IsShared: false,
    Createddate: "",
    Updateddate: "",
    Lastused: "",
    Filepath: "",
  };

  const username = currentUser?.username ? currentUser.username : "";

  const fetchChatbots = async () => {
    try {
      const response = await getChatbotsApi.get("", {
        headers: { Authorization: `Bearer ${token}` },
      });
      setChatbots(response.data);
    } catch (error) {
      console.error("Failed to fetch chatbots:", error);
    }
  };

  const findSelectedChatbot = (chatbotID: number | null) => {
    if (!chatbotID) {
      return null;
    }
    const found = chatbots.find((chatbot) => chatbot.Chatbotid === chatbotID);

    if (found) {
      return found;
    } else {
      console.log("error: Chatbot id not found");
      return null;
    }
  };

  useEffect(() => {
    if (currentUser) {
      fetchChatbots();
    }
  }, [currentUser]);

  return (
    <div className="flex h-screen flex-1 w-full">
      <Sidebar
        currentUsername={username}
        chatbots={chatbots}
        onCreateNewChatbot={() => {
          setIsCreatingChatbot(true);
          setSelectedChatbotID(null);
          setSelectedChatbot(findSelectedChatbot(selectedChatbotID));
        }}
        selectChatbot={(selectedChatbotID) => {
          setIsCreatingChatbot(false);
          setSelectedChatbotID(selectedChatbotID);
          setSelectedChatbot(findSelectedChatbot(selectedChatbotID));
        }}
      />
      <div className="w-full">
      {isCreatingChatbot ? (
        <Botconfigs username={username} chatbot={newBot} />
      ) : selectedChatbotID && selectedChatbot ? (
        <Botconfigs username={username} chatbot={selectedChatbot} />
      ) : (
        <h1 className="text-2xl font-bold p-4">Select a chatbot to view</h1>
      )}
      </div>
    </div>
  );
}

export default Dashboard;

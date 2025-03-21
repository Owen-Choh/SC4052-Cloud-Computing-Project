import useAuth from "../auth/useAuth";
import { useEffect, useState } from "react";
import Sidebar from "../components/Sidebar";
import Botconfigs from "../components/Botconfigs";

import { getChatbotsListApi } from "../api/apiConfig";
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
    chatbotid: null,
    userid: currentUser?.userid ? currentUser.userid : 0,
    chatbotname: "",
    description: "",
    behaviour: "",
    usercontext: "",
    isShared: false,
    createddate: "",
    updateddate: "",
    lastused: "",
    filepath: "",
    file: null,
  };

  const username = currentUser?.username ? currentUser.username : "";

  const fetchChatbots = async () => {
    try {
      const response = await getChatbotsListApi.get("", {
        headers: { Authorization: `Bearer ${token}` },
      });
      setChatbots(response.data);
      console.log("Chatbots fetched:", response.data);
    } catch (error) {
      console.error("Failed to fetch chatbots:", error);
    }
  };

  const findSelectedChatbot = (chatbotID: number | null) => {
    console.log("findSelectedChatbot: ", chatbotID);
    if (!chatbotID) {
      return null;
    }
    const found = chatbots.find((chatbot) => chatbot.chatbotid === chatbotID);

    if (found) {
      return { ...found };
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

  useEffect(() => {
    setSelectedChatbot(findSelectedChatbot(selectedChatbotID));
  }, [selectedChatbotID, chatbots]);

  return (
    <div className="flex h-screen flex-1 w-full">
      <Sidebar
        currentUsername={username}
        chatbots={chatbots}
        onCreateNewChatbot={() => {
          setIsCreatingChatbot(true);
          setSelectedChatbotID(null);
        }}
        selectChatbot={(selectedChatbotID) => {
          setIsCreatingChatbot(false);
          setSelectedChatbotID(selectedChatbotID);
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

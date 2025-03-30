import useAuth from "../auth/useAuth";
import { useEffect, useState } from "react";
import Sidebar from "../components/Sidebar";
import Botconfigs from "../components/Botconfigs";

import { getChatbotsListApi } from "../api/apiConfig";
import { Chatbot, ChatbotServerResponse } from "../api/chatbot";
import { useChatbotContext } from "../context/ChatbotContext";

function Dashboard() {
  const { currentUser, doLogout } = useAuth();
  const { setChatbots, selectedChatbot, isCreatingChatbot } =
    useChatbotContext();

  const [currentChatbot, setCurrentChatbot] = useState<Chatbot | null>(null);
  const [excludeFile, setExcludeFile] = useState(false);

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

  if (!currentUser) {
    doLogout();
    return;
  }

  const fetchChatbots = async () => {
    try {
      const response = await getChatbotsListApi.get<ChatbotServerResponse>("", {
        withCredentials: true,
      });

      const chatbotsData = Array.isArray(response.data)
        ? response.data.map((bot) => {
            return {
              ...bot,
              prevFilePath: bot.filepath,
              file: null,
            };
          })
        : [];
      setChatbots(chatbotsData);
      // console.log("Chatbots fetched:", response.data);
    } catch (error) {
      console.error("Failed to fetch chatbots:", error);
    }
  };

  useEffect(() => {
    if (currentUser) {
      fetchChatbots();
    }
  }, [currentUser]);

  useEffect(() => {
    if (isCreatingChatbot) {
      setCurrentChatbot(newBot);
    } else if (selectedChatbot) {
      setCurrentChatbot(selectedChatbot);
    }
  }, [isCreatingChatbot, selectedChatbot]);

  return (
    <div className="flex h-screen flex-1 w-full items-center">
      <Sidebar currentUsername={currentUser?.username} />
      <div className="w-full h-full flex items-center justify-center">
        {currentChatbot ? (
          <Botconfigs
            chatbot={currentChatbot}
            setChatbot={setCurrentChatbot}
            excludeFile={excludeFile}
            setExcludeFile={setExcludeFile}
          />
        ) : (
          <h1 className="text-2xl font-bold p-4 text-center">
            Click on the sidebar to create a new chatbot or <br /> select an
            existing chatbot to view your customisations
          </h1>
        )}
      </div>
    </div>
  );
}

export default Dashboard;

import React, { useEffect, useState } from "react";
import Tab from "./ui/Tab";
import TabPanel from "./ui/TabPanel";
import ChatbotInformation from "./ChatbotInformation";
import ChatbotCustomisation from "./ChatbotCustomisation";
import { Chatbot } from "../api/chatbot";
import { chatbotsApi } from "../api/apiConfig";
import useAuth from "../auth/useAuth";
import { useChatbotContext } from "../context/ChatbotContext";
import DeleteModal from "./DeleteModal";

interface BotconfigsProps {
  chatbot: Chatbot;
  setChatbot: React.Dispatch<React.SetStateAction<Chatbot | null>>;
  excludeFile: boolean;
  setExcludeFile: React.Dispatch<React.SetStateAction<boolean>>;
}

const Botconfigs: React.FC<BotconfigsProps> = ({
  chatbot,
  setChatbot,
  excludeFile,
  setExcludeFile,
}) => {
  const { currentUser } = useAuth();
  const {
    isCreatingChatbot,
    setIsCreatingChatbot,
    addChatbotInContext,
    updateChatbotInContext,
    deleteChatbotInContext,
  } = useChatbotContext();
  const [activeTab, setActiveTab] = useState("chatInfo");
  const [chatbotLink, setChatbotLink] = useState(
    `/chat/${currentUser?.username}/${chatbot.chatbotname}`
  );
  const [success, setSuccess] = useState("");
  const [error, setError] = useState("");
  const [resetMessages, setResetMessages] = useState(true);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);

  const updateChatbotInfo = (
    chatbotName: string,
    isShared: boolean,
    description: string
  ) => {
    setChatbot((prev) =>
      prev ? { ...prev, chatbotname: chatbotName, isShared, description } : prev
    );
    setChatbotLink(`/chat/${currentUser?.username}/${chatbotName}`);
  };

  const updateChatbotCustomisation = (behaviour: string, context: string) => {
    setChatbot((prev) =>
      prev ? { ...prev, behaviour, usercontext: context } : prev
    );
  };

  const updateChatbotFile = (document: File | null) => {
    setChatbot((prev) =>
      prev
        ? {
            ...prev,
            filepath: document ? document.name : prev.filepath,
            file: document,
          }
        : prev
    );
  };

  const saveChatbot = async () => {
    if (!chatbot) return;
    setSuccess("");
    setError("");

    if (!chatbot.chatbotname || chatbot.chatbotname.length < 1) {
      setError("Chatbot name is required.");
      return;
    } else if (/^[a-zA-Z0-9_-]*$/.test(chatbot.chatbotname) === false) {
      setError(
        `Chatbot name ${chatbot.chatbotname} can only contain alphanumeric, _ or - characters. It cannot contain special characters or spaces.`
      );
      return;
    }

    if (chatbot.file) {
      if (chatbot.file.size > 10 * 1024 * 1024) {
        setError("File size exceeds 10MB limit.");
        return;
      } else if (/^[a-zA-Z0-9_\-\. ]+$/.test(chatbot.file.name) === false) {
        setError(
          `File name ${chatbot.file.name} must be alphanumeric and cannot contain special characters or spaces.`
        );
        return;
      } else if (
        ![
          "application/pdf",
          "text/plain",
          "image/jpeg",
          "audio/mpeg",
          "video/mp4",
        ].includes(chatbot.file.type)
      ) {
        // can only be pdf, txt, jpg, mp3, mp4
        setError(
          "Invalid file type. Please upload a valid file. Only pdf, txt, jpg, mp3, mp4 are allowed."
        );
        return;
      }
    }

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
            headers: {
              "Content-Type": "multipart/form-data",
            },
            withCredentials: true,
          })
        : await chatbotsApi.post("/", formData, {
            headers: {
              "Content-Type": "multipart/form-data",
            },
            withCredentials: true,
          });

      // console.log("Chatbot saved successfully:", response.data);
      setSuccess("Chatbot saved successfully!");
      setError("");
      if (isCreatingChatbot) {
        // Update chatbot id if user create new chatbot
        const updatedChatbot = {
          ...chatbot,
          chatbotid: response.data.chatbotid,
          file: null,
          prevFilePath: chatbot.filepath,
        };
        setChatbot(updatedChatbot);
        setIsCreatingChatbot(false);
        addChatbotInContext(updatedChatbot);
      } else {
        setIsCreatingChatbot(false);
        updateChatbotInContext(chatbot);
      }
      setResetMessages((prev) => !prev);
    } catch (err: any) {
      console.error("Failed to save chatbot:", err);
      setSuccess("");
      setError(
        "Failed to save chatbot. " +
          (err.response?.data?.error || "Unknown error")
      );
    }
  };

  const deleteChatbot = async () => {
    if (!chatbot) return;
    try {
      const response = await chatbotsApi.delete(`/${chatbot.chatbotid}`, {
        withCredentials: true,
      });
      if (response.status !== 200) {
        throw new Error("Failed to delete chatbot.");
      }

      // console.log("Chatbot deleted successfully:", response.data);
      setDeleteModalOpen(false);
      // setSuccess("Chatbot deleted successfully!");
      setError("");
      setChatbot(null);
      setIsCreatingChatbot(false);
      setResetMessages((prev) => !prev);
      deleteChatbotInContext(chatbot);
      alert("Chatbot deleted successfully!");
    } catch (err: any) {
      console.error("Failed to delete chatbot:", err);
      setSuccess("");
      setError(
        "Failed to delete chatbot. " +
          (err.response?.data?.error || "Unknown error")
      );
      setDeleteModalOpen(false);
    }
  };

  useEffect(() => {
    if (resetMessages) {
      setResetMessages((prev) => !prev);
    } else {
      setSuccess("");
      setError("");
    }
    setChatbotLink(`/chat/${currentUser?.username}/${chatbot.chatbotname}`);
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
            chatbotDocument={chatbot.prevFilePath}
            excludeFile={excludeFile}
            toggleExcludeFile={() => setExcludeFile((prev) => !prev)}
            updateChatbotCustomisation={updateChatbotCustomisation}
            updateChatbotFile={updateChatbotFile}
          />
        </TabPanel>
        {success && <p className="p-4 text-green-500">{success}</p>}
        {error && <p className="p-4 text-red-500">{error}</p>}
      </div>
      <div className="border-b-2 border-gray-700"></div>
      <div className="flex w-full justify-end gap-4">
        <button
          className="bg-green-600 p-2 rounded hover:bg-green-700"
          onClick={saveChatbot}
        >
          Save Changes
        </button>
        {!isCreatingChatbot && (
          <>
            <button
              className="bg-red-600 p-2 rounded hover:bg-red-700"
              onClick={() => setDeleteModalOpen(true)}
            >
              Delete Chatbot
            </button>
            <DeleteModal
              open={deleteModalOpen}
              handleClose={() => setDeleteModalOpen(false)}
              handleDelete={() => deleteChatbot()}
              title={chatbot.chatbotname}
            ></DeleteModal>
          </>
        )}
      </div>
    </div>
  );
};

export default Botconfigs;

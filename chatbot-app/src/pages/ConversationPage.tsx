import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getConversationIdApi, getTestConversationApi } from "../api/apiConfig";
import ReactMarkdown from "react-markdown";
import SendIcon from "@mui/icons-material/Send";

const ConversationPage = () => {
  const { username, chatbotname } = useParams(); // Extract URL params
  const [conversationID, setConversationID] = useState<number | null>(null);
  const [conversation, setConversation] = useState<string[]>([]);
  const [userInput, setUserInput] = useState<string>("");

  const getConversationID = async () => {
    // Fetch conversation ID from server
    const conversationIDresponse = await getConversationIdApi.get("");
    console.log("Conversation ID:", conversationIDresponse);
    setConversationID(conversationIDresponse.data.conversationid);
  };

  const getTestConversation = async () => {
    // Fetch conversation from server, need to be informat /{username}/{chatbotname}
    if (userInput == "") {
      console.log("User input is empty");
      return;
    }

    // Add user input as a new message (from the user)
    setConversation((prev) => [...prev, `**You:**\n> ${userInput}`]);
    const userText = userInput; // Store before clearing input
    setUserInput("");

    try {
      const testConversationResponse = await getTestConversationApi.post(
        `/${username}/${chatbotname}`,
        { userInput: userText }
      );
      console.log("Test conversation:", testConversationResponse.data.response);
      // Append chatbot response
      setConversation((prev) => [
        ...prev,
        `**${chatbotname}:**\n> ${testConversationResponse.data.response}`,
      ]);
    } catch (error) {
      console.error("Error fetching chatbot response:", error);
    }
  };

  useEffect(
    () => {
      getConversationID();
    },
    [] // Run once
  );

  return (
    <div className="p-4 flex flex-col h-screen gap-2">
      <div>
        <h1 className="text-2xl font-bold">
          Chatting with {chatbotname} by user {username}
        </h1>
        <p>Conversation ID: {conversationID ? conversationID : "Loading..."}</p>
      </div>
      <div className="border-b-2 border-gray-700"></div>

      <div className="flex flex-col flex-grow gap-2">
        {conversation.map((line, index) => {
          const isUser = line.startsWith("**You:**");
          return (
            <div
              key={index}
              className={`border p-4 rounded-lg ${
                isUser ? "bg-gray-700" : ""
              }`}
            >
              <ReactMarkdown>{line}</ReactMarkdown>
            </div>
          );
        })}
      </div>
      <div className="border p-4 rounded-lg sticky bottom-0 bg-gray-800 flex">
        <textarea
          className="border rounded flex-grow p-2 "
          placeholder="Type a message..."
          value={userInput}
          onChange={(e) => setUserInput(e.target.value)}
        />
        <button
          className="border rounded-lg m-2 flex gap-1 items-center"
          onClick={getTestConversation}
        >
          Send
          <SendIcon />
        </button>
      </div>
    </div>
  );
};

export default ConversationPage;

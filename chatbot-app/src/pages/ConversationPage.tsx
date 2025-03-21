import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getConversationIdApi, chatConversationApi } from "../api/apiConfig";
import ReactMarkdown from "react-markdown";
import SendIcon from "@mui/icons-material/Send";
import axios, { AxiosError, HttpStatusCode } from "axios";

export type ConversationSuccessResponse = {
  conversationid: string;
  description: string;
};

const ConversationPage = () => {
  const { username, chatbotname } = useParams(); // Extract URL params
  const [conversationID, setConversationID] = useState<string>("");
  const [conversation, setConversation] = useState<string[]>([]);
  const [userInput, setUserInput] = useState<string>("");
  const [chatbotDescription, setChatbotDescription] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getConversationID = async () => {
    try {
      // Fetch conversation ID from server
      const response = await getConversationIdApi.get(
        `/${username}/${chatbotname}`
      );

      const conversationResponse = response.data as ConversationSuccessResponse;
      console.log("Conversation start response object:", conversationResponse);
      setConversationID(conversationResponse.conversationid);
      setChatbotDescription(conversationResponse.description);
    } catch (error) {
      if (axios.isAxiosError(error)) {
        // if error is due to timeout
        if (
          error.code === AxiosError.ECONNABORTED ||
          error.code === AxiosError.ERR_NETWORK
        ) {
          setError("Unable to reach chatbot server. Please try again later :(");
        } else if (error.response?.status === HttpStatusCode.NotFound) {
          setError("This chatbot does not exist. Is your url correct?");
        } else if (error.response?.status === HttpStatusCode.Forbidden) {
          setError("This chatbot is not shared. Please check with the owner.");
        } else if (
          error.response?.status === HttpStatusCode.InternalServerError
        ) {
          setError(
            "An error occurred while starting the conversation. Please try again later :("
          );
        } else {
          setError("An unknown error occurred. Please try again later :(");
        }
      }
    }
  };

  const sendConversation = async () => {
    // Fetch conversation from server, need to be informat /{username}/{chatbotname}
    if (userInput == "") {
      console.log("User input is empty");
      return;
    }

    // Add user input as a new message (from the user)
    setConversation((prev) => [...prev, `**You:**\n> ${userInput}`]);
    const userText = userInput; // Store before clearing input
    setUserInput("");
    setLoading(true);

    try {
      const chatConversationResponse = await chatConversationApi.post(
        `/${username}/${chatbotname}`,
        {
          conversationid: conversationID,
          message: userText,
        }
      );
      // Append chatbot response
      setConversation((prev) => [
        ...prev,
        `**${chatbotname}:**\n> ${chatConversationResponse.data.response}`,
      ]);
    } catch (error) {
      if (axios.isAxiosError(error)) {
        setConversation((prev) => [
          ...prev,
          `**${chatbotname}:**\n> Error ${error.response?.status}: ${error.response?.data.error}`,
        ]);
      }
    } finally {
      setLoading(false);
    }
  };

  const downloadConversationAsText = () => {
    const content = conversation.join("\n\n");
    const blob = new Blob([content], { type: "text/plain;charset=utf-8" });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${chatbotname}_conversation_${conversationID}.txt`;
    a.click();
    window.URL.revokeObjectURL(url);
  };

  const downloadConversationAsMarkdown = () => {
    const content = conversation.join("\n\n"); // Spacing between messages
    const blob = new Blob([content], { type: "text/markdown;charset=utf-8" });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${chatbotname}_conversation_${conversationID}.md`;
    a.click();
    window.URL.revokeObjectURL(url);
  };

  useEffect(
    () => {
      getConversationID();
    },
    [] // Run once
  );

  return (
    <div className="p-4 flex flex-col h-screen gap-2">
      <div className="flow-root">
        <div className="float-left w-3/4">
          <h1 className="text-2xl font-bold underline">
            Chatting with {chatbotname} by user {username}
          </h1>
          {error ? (
            <div className="text-red-500 font-bold text-1xl">{error}</div>
          ) : (
            <>
              <p>
                Conversation ID:{" "}
                {conversationID && conversationID != ""
                  ? conversationID
                  : "Loading..."}
              </p>
              <p>
                Description of chatbot:{" "}
                {chatbotDescription && chatbotDescription != ""
                  ? chatbotDescription
                  : "Loading..."}
              </p>
            </>
          )}
        </div>
        {conversationID == "" ? null : (
          <div className="float-right flex flex-col gap-2">
            <button
              className="border rounded-lg p-1 bg-green-600 hover:bg-green-700 text-white"
              onClick={downloadConversationAsMarkdown}
            >
              Download as Markdown
            </button>
            <button
              className="border rounded-lg p-1 bg-green-600 hover:bg-green-700 text-white"
              onClick={downloadConversationAsText}
            >
              Download as Text file
            </button>
          </div>
        )}
      </div>
      <div className="border-b-2 border-gray-700"></div>

      <div className="flex flex-col flex-grow gap-2">
        {conversation.map((line, index) => {
          const isUser = line.startsWith("**You:**");
          return (
            <div
              key={index}
              className={`border p-4 rounded-lg ${isUser ? "bg-gray-700" : ""}`}
            >
              <ReactMarkdown>{line}</ReactMarkdown>
            </div>
          );
        })}
      </div>
      {loading && (
        <div className="border p-4 rounded-lg bg-gray-600 italic text-gray-300">
          {`${chatbotname} is thinking very hard, this may take up to a minute...`}
        </div>
      )}
      <div className="border p-4 rounded-lg sticky bottom-0 bg-gray-800 flex">
        <textarea
          className="border rounded flex-grow p-2 "
          placeholder="Type a message..."
          value={userInput}
          onChange={(e) => setUserInput(e.target.value)}
        />
        <button
          className="border rounded-lg m-2 flex gap-1 items-center disabled:opacity-50"
          onClick={sendConversation}
          disabled={loading || userInput == ""}
        >
          Send
          <SendIcon />
        </button>
      </div>
    </div>
  );
};

export default ConversationPage;

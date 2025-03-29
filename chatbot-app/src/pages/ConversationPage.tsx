import { useEffect, useState, useRef } from "react";
import { useParams } from "react-router-dom";
import {
  getConversationIdApi,
  chatConversationApi,
  chatStreamConversationApiUrl, // Import the streaming API
} from "../api/apiConfig";
import ReactMarkdown from "react-markdown";
import SendIcon from "@mui/icons-material/Send";
import axios, { AxiosError, HttpStatusCode } from "axios";
import remarkGfm from "remark-gfm";

export type ConversationSuccessResponse = {
  conversationid: string;
  description: string;
};

const ConversationPage = () => {
  const { username, chatbotname } = useParams(); // Extract URL params
  const [conversationID, setConversationID] = useState<string>("");
  const [conversation, setConversation] = useState<
    { role: "user" | "chatbot"; content: string }[]
  >([]); // Array of objects to manage user/chatbot messages
  const [userInput, setUserInput] = useState<string>("");
  const [chatbotDescription, setChatbotDescription] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isStreaming, setIsStreaming] = useState(false); // For streaming response
  const [geminiResponse, setGeminiResponse] = useState(""); // For streaming response display
  const responseAreaRef = useRef(null); // Ref for response display area during streaming

  const getConversationID = async () => {
    try {
      // Fetch conversation ID from server
      const response = await getConversationIdApi.get(
        `/${username}/${chatbotname}`
      );

      const conversationResponse = response.data as ConversationSuccessResponse;
      console.log("Conversation start response object:", conversationResponse);
      // setConversationID(conversationResponse.conversationid);
      setConversationID("6176875e-e0ca-4bf8-a8f2-8f1a59ba36b5");
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
          setError(
            "This chatbot is not available to use. Please check with the owner."
          );
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
    if (userInput === "") {
      console.log("User input is empty");
      return;
    }

    const userMessage = userInput;
    setUserInput("");
    setLoading(true);
    setError(null); // Clear any previous errors
    setGeminiResponse(""); // Clear previous streaming response
    setConversation((prev) => [
      ...prev,
      { role: "user", content: userMessage }, // Add user message to conversation state
    ]);

    if (isStreaming) {
      try {
        const response = await fetch(
          chatStreamConversationApiUrl + `/${username}/${chatbotname}`, // Use absolute URL
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              conversationid: conversationID,
              message: userMessage,
            }),
          }
        );

        if (!response.ok) {
          console.error("HTTP error!", response.status, response.statusText);
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        if (!response.body) {
          console.error("Response body is null or undefined.");
          setError(
            "Error: Response body is null or undefined. Please try again."
          );
          return;
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let accumulatedResponse = "";
        let chatbotFullResponse = ""; // To store the full response for conversation history

        const processStream = () => {
          reader
            .read()
            .then(({ done, value }) => {
              if (done) {
                console.log("Stream completed.");
                console.log("Full response:", chatbotFullResponse);
                setConversation((prev) => [
                  ...prev,
                  { role: "chatbot", content: chatbotFullResponse }, // Add full chatbot response
                ]);
                setLoading(false); // Remove loading state
                setGeminiResponse("");
                return;
              }
              const decodedChunk = decoder.decode(value, { stream: true });
              // remove data: prefix from the decoded chunk
              console.log("Decoded chunk:", decodedChunk);
              var cleanedChunk = decodedChunk.replace(/^data:\s/, "");
              console.log("Cleaned chunk:", cleanedChunk);

              if (cleanedChunk.endsWith('\n\n')) {
                console.log("the pair of trailing newlines detected, removing them.");
                cleanedChunk = cleanedChunk.slice(0, -2); // Remove trailing \n\n if it exists
              }

              if (cleanedChunk === "event: close\ndata: done") {
                console.log("Stream completed.");
                console.log("Full response:", chatbotFullResponse);
                setConversation((prev) => [
                  ...prev,
                  { role: "chatbot", content: chatbotFullResponse }, // Add full chatbot response
                ]);
                setLoading(false); // Remove loading state
                setGeminiResponse("");
                return;
              }

              accumulatedResponse += cleanedChunk;
              chatbotFullResponse += cleanedChunk; // Append to full response
              
              console.log("Accumulated response:", accumulatedResponse);

              setGeminiResponse(accumulatedResponse); // Update streaming UI
              processStream(); // Continue reading the stream
            })
            .catch((streamError) => {
              console.error("Stream reading error:", streamError);
              setError(
                "Error streaming response from chatbot. Please try again."
              );
              setIsStreaming(false); // Streaming stopped due to error
              setLoading(false);
              reader.cancel(); // Cancel the reader on error
            });
        };

        processStream(); // Start streaming

        // Clear geminiResponse here, before stream starts to ensure clean UI
        setGeminiResponse("");
      } catch (fetchError) {
        console.error("Fetch error:", fetchError);
        setError(
          "Error fetching streaming response. Please check your network and try again."
        );
        setIsStreaming(false); // Streaming stopped due to fetch error
        setLoading(false);
      }
    } else {
      // Non-streaming API call
      try {
        const chatConversationResponse = await chatConversationApi.post(
          `/${username}/${chatbotname}`,
          {
            conversationid: conversationID,
            message: userMessage,
          }
        );
        // Append chatbot response
        setConversation((prev) => [
          ...prev,
          {
            role: "chatbot",
            content: chatConversationResponse.data.response,
          },
        ]);
      } catch (error) {
        if (axios.isAxiosError(error)) {
          setConversation((prev) => [
            ...prev,
            {
              role: "chatbot",
              content: `Error ${error.response?.status}: ${error.response?.data.error}`,
            },
          ]);
        }
      } finally {
        setLoading(false);
      }
    }
  };

  const downloadConversationAsText = () => {
    const content = conversation
      .map(
        (msg) =>
          `**${msg.role === "user" ? "You" : chatbotname}:**\n> ${msg.content}`
      )
      .join("\n\n");
    const blob = new Blob([content], { type: "text/plain;charset=utf-8" });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${chatbotname}_conversation_${conversationID}.txt`;
    a.click();
    window.URL.revokeObjectURL(url);
  };

  const downloadConversationAsMarkdown = () => {
    const content = conversation
      .map(
        (msg) =>
          `**${msg.role === "user" ? "You" : chatbotname}:**\n> ${msg.content}`
      )
      .join("\n\n");
    const blob = new Blob([content], { type: "text/markdown;charset=utf-8" });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${chatbotname}_conversation_${conversationID}.md`;
    a.click();
    window.URL.revokeObjectURL(url);
  };

  const renderers = {
    code({ node, inline, className, children, ...props }: any) {
      if (inline) {
        return <code className="inline-code">{children}</code>; // Correctly renders inline code
      }
  
      const match = /language-(\w+)/.exec(className || "");
      const language = match ? match[1].toUpperCase() : null; // Detect language
  
      if (!language) {
        return <code className="inline-code">{children}</code>; // Return inline code if no language is specified
      }

      return (
        <div className="code-block-container m-2">
          {language && <div className="code-language-label w-full">{language}</div>}
          <pre {...props} className={className}>
            <code className="!p-2">{children}</code>
          </pre>
        </div>
      );
    },
  };

  useEffect(() => {
    getConversationID();
  }, []); // Run once on component mount

  return (
    <div className="p-4 flex flex-col h-screen gap-2">
      <div className="flow-root">
        <div className="float-left w-3/4">
          <h1 className="text-2xl font-bold underline">
            Chatting with {chatbotname} by user {username}
          </h1>
          {error ? (
            <div className="text-red-500 font-bold text-2xl">{error}</div>
          ) : (
            <>
              <p>
                Conversation ID:{" "}
                {conversationID && conversationID !== ""
                  ? conversationID
                  : "Loading..."}
              </p>
              {conversationID !== "" && chatbotDescription === "" ? (
                <p>No chatbot description provided</p>
              ) : (
                <p>
                  Description of chatbot:{" "}
                  {chatbotDescription && chatbotDescription !== ""
                    ? chatbotDescription
                    : "Loading..."}
                </p>
              )}
            </>
          )}
        </div>
        {conversationID === "" ? null : (
          <div className="float-right flex flex-col gap-2">
            <button
              className="border rounded-lg p-1 bg-green-600 hover:bg-green-700 text-white"
              onClick={downloadConversationAsMarkdown}
            >
              Download chat as <span className="font-bold">md</span> file
            </button>
            <button
              className="border rounded-lg p-1 bg-green-600 hover:bg-green-700 text-white"
              onClick={downloadConversationAsText}
            >
              Download chat as <span className="font-bold">txt</span> file
            </button>
          </div>
        )}
      </div>
      <div className="border-b-2 border-gray-700"></div>

      <div className="flex flex-col flex-grow gap-2">
        {conversation.map((msg, index) => {
          const isUser = msg.role === "user";
          return (
            <div
              key={index}
              className={`border p-4  rounded-lg markdown-body !py-0 ${isUser ? "bg-gray-700" : ""}`}
            >
              <ReactMarkdown 
              skipHtml={true} 
              remarkPlugins={[remarkGfm]}
              components={renderers} 
              >{`**${
                isUser ? "You" : chatbotname
              }:**\n\n ${msg.content}`}</ReactMarkdown>
            </div>
          );
        })}
        {loading && isStreaming && (
          <div ref={responseAreaRef} className="border p-4 rounded-lg">
            <ReactMarkdown skipHtml={true}>{`**${chatbotname} (Streaming):**\n> ${geminiResponse}`}</ReactMarkdown>
          </div>
        )}
      </div>
      {loading &&
        !isStreaming && ( // Show loading only for non-streaming
          <div className="border p-4 rounded-lg bg-gray-600 italic text-gray-300">
            {`${chatbotname} is thinking very hard, this may take up to a minute...`}
          </div>
        )}
      {loading &&
        isStreaming && ( // Show streaming loading message
          <div className="border p-4 rounded-lg bg-gray-600 italic text-gray-300">
            {`${chatbotname} is responding in real-time...`}
          </div>
        )}
      <div className="border p-4 rounded-lg sticky bottom-0 bg-gray-800 flex">
        <textarea
          className="border rounded flex-grow p-2 "
          placeholder="Type a message..."
          value={userInput}
          onChange={(e) => setUserInput(e.target.value)}
        />
        <div className="w-fit">
          <button
            className="border rounded-lg flex gap-1 items-center m-auto disabled:opacity-50"
            onClick={sendConversation}
            disabled={loading || userInput === ""}
          >
            Send
            <SendIcon />
          </button>
          <input
            type="checkbox"
            className="m-2"
            id="streamingCheckbox"
            checked={isStreaming}
            disabled={loading}
            onChange={(e) => setIsStreaming(e.target.checked)}
          />
          <label htmlFor="streamingCheckbox">Stream Response?</label>
        </div>
      </div>
    </div>
  );
};

export default ConversationPage;

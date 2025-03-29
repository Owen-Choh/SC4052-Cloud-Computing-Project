import React, { useState, useEffect, useRef } from "react";
import { useParams } from "react-router-dom";

function GeminiStreamDisplay() {
  const { username, chatbotname } = useParams(); // Extract URL params
  const [geminiResponse, setGeminiResponse] = useState("");
  const [messageInput, setMessageInput] = useState(""); // Input for message
  const [conversationIdInput, setConversationIdInput] =
    useState("qw-rghr5341-5136"); // Input for conversation ID
  const responseAreaRef = useRef(null); // Ref for response display area

  const fetchData = async () => {
    try {
      const response = await fetch(
        `http://localhost:8080/api/conversation/chat/stream/${username}/${chatbotname}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json", // Indicate JSON payload
          },
          body: JSON.stringify({
            // Send JSON payload
            conversationid: conversationIdInput,
            message: messageInput,
          }),
        }
      );

      if (!response.ok) {
        console.error("HTTP error!", response.status, response.statusText);
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let accumulatedResponse = "";

      const processStream = () => {
        reader
          .read()
          .then(({ done, value }) => {
            if (done) {
              console.log("Stream completed.");
              return;
            }
            accumulatedResponse += decoder.decode(value, { stream: true });
            setGeminiResponse(accumulatedResponse);
            responseAreaRef.current.scrollTop =
              responseAreaRef.current.scrollHeight; // Scroll to bottom
            processStream(); // Continue reading the stream
          })
          .catch((error) => {
            console.error("Stream reading error:", error);
            reader.cancel(); // Cancel the reader on error
          });
      };

      processStream(); // Start processing the stream
    } catch (error) {
      console.error("Fetch error:", error);
      setGeminiResponse(`Error: ${error.message}`); // Display error in UI
    }
  };

  useEffect(() => {
    if (!conversationIdInput) return; // Don't start stream if conversation ID is not set
    if (!messageInput) return; // Don't start stream if message is empty

    setGeminiResponse(""); // Clear previous response
  }, []); // Re-run effect when these change

  const handleSendMessage = () => {
    if (conversationIdInput && messageInput) {
      // The useEffect will handle sending the request when conversationIdInput or messageInput changes
      fetchData();
    } else {
      alert("Please enter Conversation ID and Message.");
    }
  };

  return (
    <div>
      <h1>Gemini Streaming Response (POST Request):</h1>

      <div>
        <label htmlFor="conversationId">Conversation ID:</label>
        <input
          type="text"
          id="conversationId"
          value={conversationIdInput}
          onChange={(e) => setConversationIdInput(e.target.value)}
          placeholder="Enter Conversation ID"
        />
      </div>
      <div>
        <label htmlFor="message">Message:</label>
        <input
          type="text"
          id="message"
          value={messageInput}
          onChange={(e) => setMessageInput(e.target.value)}
          placeholder="Enter your message"
        />
      </div>
      <button onClick={handleSendMessage}>Send Message</button>

      <div
        ref={responseAreaRef}
        style={{
          border: "1px solid #ccc",
          padding: "10px",
          marginTop: "10px",
          height: "300px",
          overflowY: "scroll",
          whiteSpace: "pre-wrap",
        }}
      >
        {geminiResponse}
      </div>
    </div>
  );
}

export default GeminiStreamDisplay;

import React from "react";
import Tab from "./ui/Tab";
import TabPanel from "./ui/TabPanel";
import ChatbotInformation from "./ChatbotInformation";
import ChatbotCustomisation from "./ChatbotCustomisation";

const Botconfigs: React.FC = () => {
  const [activeTab, setActiveTab] = React.useState("chatInfo");
  const chatbotName = "My Demo Chatbot";
  const isShared = false;
  const chatbotEndpoint = "/testuser/my-demo-chatbot";
  const chatbotLink = "http://localhost:5173/testuser/my-demo-chatbot";

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
          <ChatbotInformation chatbotName={chatbotName} isShared={isShared} chatbotEndpoint={chatbotEndpoint} />
        </TabPanel>
        <TabPanel activeTab={activeTab} tabKey="customisation">
          <ChatbotCustomisation chatbotName={chatbotName} isShared={isShared} chatbotLink={chatbotLink} />
        </TabPanel>
      </div>
    </div>
  );
};

export default Botconfigs;

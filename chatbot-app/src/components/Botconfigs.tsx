import React from "react";
import Tab from "./ui/Tab";
import TabPanel from "./ui/TabPanel";
import ChatbotInformation from "./ChatbotInformation";
import ChatbotCustomisation from "./ChatbotCustomisation";

const Botconfigs: React.FC = () => {
  const [activeTab, setActiveTab] = React.useState("chatInfo");

  // TODO: replace mock for chatbot information component
  const chatbotName = "My Demo Chatbot";
  const isShared = false;
  const chatbotEndpoint = "/testuser/my-demo-chatbot";
  
  // TODO: replace mock for chatbot customise component
  const chatbotBehaviour = "Your friendly internet chatbot";
  const chatbotContext = "Project due in a month";
  const chatbotDocument = "project-details.pdf";

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
          <ChatbotCustomisation chatbotBehaviour={chatbotBehaviour} chatbotContext={chatbotContext} chatbotDocument={chatbotDocument} />
        </TabPanel>
      </div>
    </div>
  );
};

export default Botconfigs;

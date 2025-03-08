import React from "react";
import Tab from "./ui/Tab";
import TabPanel from "./ui/TabPanel";

const Botconfigs: React.FC = () => {
  const [activeTab, setActiveTab] = React.useState("chatInfo");

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
          <p>Content for Chatbot Information</p>
        </TabPanel>
        <TabPanel activeTab={activeTab} tabKey="customisation">
          <p>Content for Customisation</p>
        </TabPanel>
      </div>
    </div>
  );
};

export default Botconfigs;

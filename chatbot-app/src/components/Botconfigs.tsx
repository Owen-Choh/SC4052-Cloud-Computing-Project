import React from "react";


const Botconfigs: React.FC = () => {
  const [activeTab, setActiveTab] = React.useState("chatInfo");

  return (
    <div className="flex flex-col w-full h-full">
      <div className="flex gap-4">
        <h1 onClick={() => setActiveTab("chatInfo")}>Chat info</h1>
        <h1 onClick={() => setActiveTab("customisation")}>customise</h1>
      </div>
      <div className="w-full flex-grow overflow-y-auto">
        {activeTab === "chatInfo" && <p>Content for Tab 1</p>}
        {activeTab === "customisation" && <p>Content for Tab 2</p>}
      </div>
    </div>
  )
}

export default Botconfigs;
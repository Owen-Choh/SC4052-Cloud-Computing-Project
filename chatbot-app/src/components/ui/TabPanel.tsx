import React from "react";

interface TabPanelProps {
  activeTab: string;
  tabKey: string;
  children: React.ReactNode;
}

const TabPanel: React.FC<TabPanelProps> = ({ activeTab, tabKey, children }) => {
  if (activeTab !== tabKey) return null; // Hide content if not active
  return <div className="p-4">{children}</div>;
};

export default TabPanel;

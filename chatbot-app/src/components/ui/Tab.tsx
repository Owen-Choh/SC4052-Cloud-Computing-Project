import React from "react";

interface TabProps {
  label: string;
  isActive: boolean;
  onClick: () => void;
}

const Tab: React.FC<TabProps> = ({ label, isActive, onClick }) => {
  return (
    <p
      onClick={onClick}
      className={`flex-1 text-center cursor-pointer rounded-md p-2 text-2xl transition-colors ${
        isActive ? "bg-blue-700 text-white" : "bg-blue-900 text-gray-300"
      }`}
    >
      {label}
    </p>
  );
};

export default Tab;
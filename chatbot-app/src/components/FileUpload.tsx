import React from "react";

interface FileUploadProps {
  fileInputRef: React.RefObject<HTMLInputElement|null>;
  onFileSelect: (file: File | null) => void;
}

const FileUpload: React.FC<FileUploadProps> = ({ fileInputRef, onFileSelect }) => {
  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0] || null;
    if (file) {
      onFileSelect(file);
    }
  };

  return (
    <div className="flex gap-4 items-center">
      <p>Upload a file</p>
      <input
        ref={fileInputRef}
        type="file"
        className="p-2 border rounded"
        onChange={handleFileChange}
      />
    </div>
  );
};

export default FileUpload;

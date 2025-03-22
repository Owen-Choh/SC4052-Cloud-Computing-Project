import Modal from "@mui/material/Modal";
import React from "react";

export interface DeleteModalProps {
  open: boolean;
  handleClose: () => void;
  handleDelete: () => void;
  title: string;
}

const DeleteModal: React.FC<DeleteModalProps> = ({
  open,
  handleClose,
  handleDelete,
  title,
}) => {
  return (
    <Modal
      open={open}
      onClose={handleClose}
      className="flex m-auto bg-blue-300 w-fit h-fit font-bold rounded-lg overflow-hidden"
    >
      <div>
        <h2 className="bg-blue-800 p-4">Delete {title}</h2>
        <h2 className="flex-grow p-4 pb-0">
          Are you sure you want to delete {title}?
        </h2>
        <div className="flex w-full justify-center gap-4 p-4">
          <button
            className="bg-green-600 p-2 rounded hover:bg-green-700"
            onClick={handleClose}
          >
            Close
          </button>
          <button
            className="bg-red-600 p-2 rounded hover:bg-red-700"
            onClick={handleDelete}
          >
            Delete
          </button>
        </div>
      </div>
    </Modal>
  );
};

export default DeleteModal;

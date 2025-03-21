export interface Chatbot {
  chatbotid: number | null;
  userid: number;
  chatbotname: string;
  description: string;
  behaviour: string;
  usercontext: string;
  isShared: boolean;
  createddate: string;
  updateddate: string;
  lastused: string;
  filepath: string;
  file: File | null;
}

export interface CreateChatbotPayload {
  chatbotname: string;
  behaviour: string;
  usercontext: string;
  isShared: boolean;
  file: File | null;
}

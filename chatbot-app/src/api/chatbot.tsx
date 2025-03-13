export interface Chatbot {
  Chatbotid: number;
  Userid: number;
  Chatbotname: string;
  Behaviour: string;
  Usercontext: string;
  IsShared: boolean;
  Createddate: string;
  Updateddate: string;
  Lastused: string;
  Filepath: string;
}

export interface CreateChatbotPayload {
  Chatbotname: string;
  Behaviour: string;
  Usercontext: string;
  IsShared: boolean;
  Filepath: string;
}

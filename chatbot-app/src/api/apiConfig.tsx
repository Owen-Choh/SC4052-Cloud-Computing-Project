import axios from "axios";

export const baseURL: string = import.meta.env.VITE_BASE_URL;

export const loginApi = axios.create({
  baseURL: `${baseURL}/user/login`,
});

export const getChatbotsApi = axios.create({
  baseURL: `${baseURL}/chatbot/list`
});

export const createChatbotsApi = axios.create({
  baseURL: `${baseURL}/chatbot/newchatbot`
});

export const getConversationIdApi = axios.create({
  baseURL: `${baseURL}/conversation/start`,
  timeout: 10000,
});

export const chatConversationApi = axios.create({
  baseURL: `${baseURL}/conversation/chat`
});
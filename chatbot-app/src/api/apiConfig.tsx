import axios from "axios";

export const baseURL: string = import.meta.env.VITE_BASE_URL;

export const loginApi = axios.create({
  baseURL: `${baseURL}/user/login`,
});

export const getChatbotsListApi = axios.create({
  baseURL: `${baseURL}/chatbot/list`
});

export const chatbotsApi = axios.create({
  baseURL: `${baseURL}/chatbot`
});

export const getConversationIdApi = axios.create({
  baseURL: `${baseURL}/conversation/start`,
  timeout: 10000,
});

export const chatConversationApi = axios.create({
  baseURL: `${baseURL}/conversation/chat`
});
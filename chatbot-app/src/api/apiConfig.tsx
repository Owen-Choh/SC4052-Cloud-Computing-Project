import axios from "axios";

export const baseURL: string = import.meta.env.VITE_BASE_URL;

export const loginApi = axios.create({
  baseURL: `${baseURL}/user/login`,
});

export const registerApi = axios.create({
  baseURL: `${baseURL}/user/register`,
});

export const logoutApi = axios.create({
  baseURL: `${baseURL}/user/logout`,
});

export const checkAuthApi = axios.create({
  baseURL: `${baseURL}/user/auth/check`,
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
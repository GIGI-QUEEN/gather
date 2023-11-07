import axios from 'axios';

const backend_url = 'http://localhost:8080';

export const api = axios.create({
  baseURL: 'http://localhost:8080',
  withCredentials: true,
});

export const postURL = async (url, data) => {
  return api.post(`http://localhost:8080${url}`, data);
};

export const getURL = async (url) => {
  return await api.get(`http://localhost:8080${url}`);
};

export const myAxios = axios.create({
  baseURL: backend_url,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

export const uploadAxios = axios.create({
  baseURL: backend_url,
  headers: {
    'Content-Type': 'multipart/form-data',
  },
  withCredentials: true,
});

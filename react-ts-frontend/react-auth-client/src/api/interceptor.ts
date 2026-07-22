import api from "./axios";

import {
  getAccessToken,
  setAccessToken,
  clearAccessToken,
} from "../auth/token";

api.interceptors.request.use((config) => {
  const token = getAccessToken();

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
});

let refreshing = false;

let requests: any[] = [];

function resolveRequests(token: string) {
  requests.forEach((callback) => callback(token));

  requests = [];
}

api.interceptors.response.use(
  (response) => response,

  async (error) => {
    const original = error.config;

    if (error.response?.status === 401 && !original._retry) {
      original._retry = true;

      if (refreshing) {
        return new Promise((resolve) => {
          requests.push((token: string) => {
            original.headers.Authorization = `Bearer ${token}`;

            resolve(api(original));
          });
        });
      }

      refreshing = true;

      try {
        const response = await api.post("/auth/refresh");

        const token = response.data.access_token;

        setAccessToken(token);

        resolveRequests(token);

        original.headers.Authorization = `Bearer ${token}`;

        return api(original);
      } catch (error) {
        clearAccessToken();

        return Promise.reject(error);
      } finally {
        refreshing = false;
      }
    }

    return Promise.reject(error);
  },
);

export default api;

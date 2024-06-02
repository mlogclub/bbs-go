import axios from 'axios';
import type { AxiosRequestConfig, AxiosResponse } from 'axios';
import { Message, Modal } from '@arco-design/web-vue';
import { useUserStore } from '@/store';
import { getToken } from '@/utils/auth';

export interface HttpResponse<T = unknown> {
  success: boolean;
  errorCode: number;
  message?: string;
  data?: T;
}

if (import.meta.env.VITE_API_BASE_URL) {
  axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL;
}

axios.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    const token = getToken();
    if (token) {
      if (!config.headers) {
        config.headers = {};
      }
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
// add response interceptors
axios.interceptors.response.use(
  (response: AxiosResponse<HttpResponse>) => {
    const res = response.data;
    if (!res.success) {
      Message.error({
        content: res.message || 'Error',
        duration: 5 * 1000,
      });
      if ([1].includes(res.errorCode)) {
        Modal.error({
          title: 'Confirm logout',
          content:
            'You have been logged out, you can cancel to stay on this page, or log in again',
          okText: 'Re-Login',
          async onOk() {
            const userStore = useUserStore();

            await userStore.logout();
            window.location.reload();
          },
        });
      }
      return Promise.reject(new Error(res.message || 'Error'));
    }
    return res.data;
  },
  (error) => {
    Message.error({
      content: error.message || 'Request Error',
      duration: 5 * 1000,
    });
    return Promise.reject(error);
  }
);

// Extend Axios interface to include the custom method
declare module 'axios' {
  interface AxiosInstance {
    postForm<T = unknown>(
      url: string,
      data: FormData,
      config?: AxiosRequestConfig
    ): Promise<T>;
  }
}

// Add a method for handling form requests
export const postForm = async <T = unknown>(
  url: string,
  data: FormData,
  config?: AxiosRequestConfig
): Promise<T> => {
  const response = await axios.post<HttpResponse<T>>(url, data, {
    ...config,
    headers: {
      'Content-Type': 'multipart/form-data',
      ...(config?.headers || {}),
    },
  });
  return response as T;
};

// Now, if you want to use axios.postForm, you can manually cast axios to any
(axios as any).postForm = postForm;
(axios as any).form = postForm;

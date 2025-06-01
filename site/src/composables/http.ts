export interface ApiResponse<T = any> {
  success: boolean;
  data: T;
  message?: string;
}

// 定义请求选项接口
interface HttpOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | 'HEAD' | 'OPTIONS';
  headers?: Record<string, string>;
  body?: any;
  params?: Record<string, any>;
  initialCache?: boolean;
  [key: string]: any;
}

// 定义参数接口
interface GetParams {
  params?: Record<string, any>;
}

function applyOptions(options: HttpOptions = {}): HttpOptions {
  options.initialCache = options.initialCache ?? false;
  options.headers = options.headers || {};
  options.method = options.method || "GET";

  const token = useCookie("bbsgo_token");
  if (token.value) {
    options.headers["X-User-Token"] = token.value;
  }
  return options;
}

//-------------------------------------------
// 使用 封装useFetch 封装 HTTP 请求
//-------------------------------------------

export function useMyFetch<T = any>(url: string, options: HttpOptions = {}) {
  options = applyOptions(options);

  return useFetch(url, {
    ...options,
    // 使用 transform 来处理响应数据
    transform: (response: ApiResponse<T> | null) => {
      if (!response) {
        return null;
      }

      // 当 success 为 true 时，只返回 data 中的数据
      if (response.success) {
        return response.data;
      }
      // 当 success 为 false 时，返回完整的响应数据
      return response;
    },
    // 自定义错误处理
    onResponseError({ response }: { response: any }) {
      // 如果响应结构中 success 为 false，可以在这里处理错误
      if (response._data && typeof response._data === 'object' && 'success' in response._data) {
        const apiResponse = response._data as ApiResponse<any>;
        if (!apiResponse.success) {
          console.error("API Error:", apiResponse.message);
          throw new Error(apiResponse.message || "API请求失败");
        }
      }
      throw new Error(`${response.status} ${response.statusText}`);
    },
  });
}

//-------------------------------------------
// 使用 $fetch 封装 HTTP 请求
//-------------------------------------------

export function useHttp<T = any>(url: string, options: HttpOptions = {}): Promise<T> {
  options = applyOptions(options);

  return new Promise((resolve, reject) => {
    $fetch(url, options as any)
      .then((resp: unknown) => {
        const apiResponse = resp as ApiResponse<T>;
        if (apiResponse.success) {
          resolve(apiResponse.data);
        } else {
          reject(apiResponse);
        }
      })
      .catch((err: any) => {
        reject(err);
      });
  });
}

export function useHttpGet<T = any>(url: string, { params }: GetParams = {}): Promise<T> {
  return useHttp<T>(url, {
    method: "GET",
    params,
  });
}

export function useJsonToForm(json: Record<string, any>): FormData {
  const formData = new FormData();
  for (const [key, value] of Object.entries(json)) {
    if (Array.isArray(value)) {
      value.forEach((item) => formData.append(`${key}[]`, item));
    } else {
      formData.append(key, value);
    }
  }
  return formData;
}

export function useHttpPost<T = any>(url: string, body: any): Promise<T> {
  return useHttp<T>(url, {
    method: "POST",
    body,
  });
}

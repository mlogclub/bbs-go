const CONFIG = {
  baseURL: process.env.BASE_URL_TEST,
};

// 请求体封装
function applyOptions(options = {}) {
  options.baseURL = options.baseURL ?? CONFIG.baseURL;
  options.initialCache = options.initialCache ?? false;
  options.headers = options.headers || {};
  options.method = options.method || "GET";

  const token = useCookie("bbsgo_token");
  if (token.value) {
    options.headers["X-User-Token"] = token.value;
  }

  // options.params = options.params || {}
  // options.params.userToken = token.value

  return options;
}

// useFetch 封装
export function useMyFetch(url, options = {}) {
  options = applyOptions(options);

  return new Promise((resolve, reject) => {
    useFetch(url, options)
      .then(({ data, pending, refresh, execute, error }) => {
        if (error.value) {
          reject(error.value);
          return;
        }
        if (data.value == null) {
          reject(new Error(`请求错误 ${url}`));
          return;
        }
        if (data.value.success) {
          resolve(data.value.data);
        } else {
          reject(data.value);
        }
      })
      .catch((err) => {
        reject(err);
      });
  });
}

// POST请求
export function useMyFetchPost(url, { body } = {}) {
  return useMyFetch(url, {
    method: "POST",
    body,
  });
}

// POST请求(application/x-www-form-urlencoded)
export function useMyFetchPostForm(url, { body = {} } = {}) {
  let bodyData = "";
  if (body && typeof body === "object") {
    for (const name in body) {
      if (bodyData.length > 0) {
        bodyData += "&";
      }
      bodyData += `${encodeURIComponent(name)}=${encodeURIComponent(
        body[name]
      )}`;
    }
  }
  return useMyFetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: bodyData,
  });
}

// POST请求(multipart/form-data)
export function useMyFetchPostMultipart(url, { body = {} } = {}) {
  const formData = new FormData();
  if (body && typeof body === "object") {
    for (const name in body) {
      formData.append(name, body[name]);
    }
  }

  return useMyFetch(url, {
    method: "POST",
    body: formData,
  });
}

export function useHttp(url, options = {}) {
  options = applyOptions(options);

  return new Promise((resolve, reject) => {
    $fetch(url, options)
      .then((resp) => {
        if (resp.success) {
          resolve(resp.data);
        } else {
          reject(resp);
        }
      })
      .catch((err) => {
        reject(err);
      });
  });
}

export function useHttpGet(url, { params } = {}) {
  return useHttp(url, {
    method: "GET",
    params,
  });
}

export function useHttpPost(url, { body } = {}) {
  return useHttp(url, {
    method: "POST",
    body,
  });
}

// POST请求(application/x-www-form-urlencoded)
export function useHttpPostForm(url, { body = {} } = {}) {
  let bodyData = "";
  if (body && typeof body === "object") {
    for (const name in body) {
      if (bodyData.length > 0) {
        bodyData += "&";
      }
      bodyData += `${encodeURIComponent(name)}=${encodeURIComponent(
        body[name]
      )}`;
    }
  }
  return useHttp(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: bodyData,
  });
}

// POST请求(multipart/form-data)
export function useHttpPostMultipart(url, formData) {
  return useHttp(url, {
    method: "POST",
    body: formData,
  });
}

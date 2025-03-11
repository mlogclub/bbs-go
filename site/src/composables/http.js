// 请求体封装
function applyOptions(options = {}) {
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

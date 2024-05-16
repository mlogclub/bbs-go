export function useFormatDate(timestamp, fmt) {
  fmt = fmt || "yyyy-MM-dd HH:mm:ss";
  const date = new Date(timestamp);
  const o = {
    "M+": date.getMonth() + 1,
    "d+": date.getDate(),
    "h+": date.getHours() % 12,
    "H+": date.getHours(),
    "m+": date.getMinutes(),
    "s+": date.getSeconds(),
    "q+": Math.floor((date.getMonth() + 3) / 3),
    S: date.getMilliseconds(),
  };
  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(
      RegExp.$1,
      `${date.getFullYear()}`.substr(4 - RegExp.$1.length)
    );
  }
  for (const k in o) {
    if (new RegExp(`(${k})`).test(fmt)) {
      fmt = fmt.replace(
        RegExp.$1,
        RegExp.$1.length === 1 ? o[k] : `00${o[k]}`.substr(`${o[k]}`.length)
      );
    }
  }
  return fmt;
}

export function usePrettyDate(timestamp) {
  const minute = 1000 * 60;
  const hour = minute * 60;
  const day = hour * 24;
  const diffValue = new Date().getTime() - timestamp;
  if (diffValue / minute < 1) {
    return "刚刚";
  } else if (diffValue / minute < 60) {
    return `${Number.parseInt(diffValue / minute)}分钟前`;
  } else if (diffValue / hour <= 24) {
    return `${Number.parseInt(diffValue / hour)}小时前`;
  } else if (diffValue / day <= 30) {
    return `${Number.parseInt(diffValue / day)}天前`;
  }
  return useFormatDate(timestamp, "yyyy-MM-dd HH:mm:ss");
}

export function useLinkTo(path) {
  const router = useRouter();
  router.push(path);
}

export function useToSignIn(redirect) {
  if (!redirect && import.meta.client) {
    redirect = window.location.pathname;
  }
  useLinkTo("/user/signin?redirect=" + encodeURIComponent(redirect));
}

/**
 * 弹出错误消息，然后执行登录
 */
export function useMsgSignIn() {
  useMsg({
    type: "error",
    message: "请先登录",
    onClose() {
      useToSignIn();
    },
  });
}

export function useCatchError(e) {
  if (e.errorCode === 1) {
    useMsgSignIn();
  } else {
    useMsgError(e.message || e);
  }
}

/**
 * Convert the given string to a unique color.
 *
 * @param {String} string
 * @return {String}
 */
export function useStringToColor(string) {
  if (string === "") return "#fff";

  let num = 0;
  // Convert the username into a number based on the ASCII value of each
  // character.
  for (let i = 0; i < string.length; i++) {
    num += string.charCodeAt(i);
  }

  // Construct a color using the remainder of that number divided by 360, and
  // some predefined saturation and value values.
  const hue = num % 360;
  const rgb = hsvToRgb(hue / 360, 0.3, 0.9);

  function hsvToRgb(h, s, v) {
    let r;
    let g;
    let b;

    const i = Math.floor(h * 6);
    const f = h * 6 - i;
    const p = v * (1 - s);
    const q = v * (1 - f * s);
    const t = v * (1 - (1 - f) * s);

    switch (i % 6) {
      case 0:
        r = v;
        g = t;
        b = p;
        break;
      case 1:
        r = q;
        g = v;
        b = p;
        break;
      case 2:
        r = p;
        g = v;
        b = t;
        break;
      case 3:
        r = p;
        g = q;
        b = v;
        break;
      case 4:
        r = t;
        g = p;
        b = v;
        break;
      case 5:
        r = v;
        g = p;
        b = q;
        break;
    }

    return {
      r: Math.floor(r * 255),
      g: Math.floor(g * 255),
      b: Math.floor(b * 255),
    };
  }

  return "" + rgb.r.toString(16) + rgb.g.toString(16) + rgb.b.toString(16);
}

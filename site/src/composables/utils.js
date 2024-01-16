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
  if (!redirect && process.client) {
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

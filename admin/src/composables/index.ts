export function jsonToFormData(json: Record<string, any>): FormData {
  const formData = new FormData();

  Object.keys(json).forEach((key) => {
    const value = json[key];
    if (value !== null && value !== undefined) {
      formData.append(key, value);
    }
  });

  return formData;
}

export function useFormatDate(
  input: Date | number,
  format = 'yyyy-MM-dd HH:mm:ss'
): string {
  const date = typeof input === 'number' ? new Date(input) : input;

  const map: Record<string, number | string> = {
    'M+': date.getMonth() + 1, // 月份
    'd+': date.getDate(), // 日
    'H+': date.getHours(), // 24小时制的小时
    'h+': date.getHours() % 12 === 0 ? 12 : date.getHours() % 12, // 12小时制的小时
    'm+': date.getMinutes(), // 分
    's+': date.getSeconds(), // 秒
    'q+': Math.floor((date.getMonth() + 3) / 3), // 季度
    'S': date.getMilliseconds(), // 毫秒
  };

  if (/(y+)/.test(format)) {
    format = format.replace(
      RegExp.$1,
      String(date.getFullYear()).substr(4 - RegExp.$1.length)
    );
  }

  Object.keys(map).forEach((key) => {
    const regex = new RegExp(`(${key})`);
    if (regex.test(format)) {
      format = format.replace(regex, (match, ...args) =>
        String(
          args[0].length === 1
            ? map[key]
            : `00${map[key]}`.substr(`${map[key]}`.length)
        )
      );
    }
  });

  return format;
}

export function isNotBlank(str: string | null | undefined): boolean {
  return str !== null && str !== undefined && str.trim().length > 0;
}

export function isBlank(str: string | null | undefined): boolean {
  return !isNotBlank(str);
}

export function isAnyBlank(
  ...strings: Array<string | null | undefined>
): boolean {
  return strings.some(isBlank);
}

export function useSiteUrl(url: string) {
  const base = import.meta.env.VITE_API_SITE_URL || '';
  return base + url;
}

export function useTableHeight() {
  const { proxy } = getCurrentInstance() as any;

  if (!proxy) {
    return;
  }

  const containerHeaderRef =
    proxy.$el.getElementsByClassName('container-header')[0];
  const containerMainRef =
    proxy.$el.getElementsByClassName('container-main')[0];
  const containerFooterRef =
    proxy.$el.getElementsByClassName('container-footer')[0];

  if (!containerMainRef) {
    return;
  }

  function calcHeight() {
    const magicHeight = 118; // 需要减去的其他高度，这些高度可能是一些边边角角的margin/padding
    const minHeight = 300; // 最低高度

    let height = document.documentElement.clientHeight - magicHeight;
    if (containerHeaderRef) {
      height -= containerHeaderRef.clientHeight;
    }
    if (containerFooterRef) {
      height -= containerFooterRef.clientHeight;
    }
    // height -= 32;
    containerMainRef.style.height = `${Math.max(height, minHeight)}px`;
  }

  calcHeight();
  window.onresize = calcHeight;
}

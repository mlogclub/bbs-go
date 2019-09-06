class Format {
  showMore(title, content) {
    window.vue.$alert(`<pre >${content}</pre>`, title, {
      dangerouslyUseHTMLString: true,
      customClass: 'show-more-dialog',
    });
  }

  formatDate(timestamp, fmt) {
    fmt = fmt || 'yyyy-MM-dd HH:mm:ss';
    const date = new Date(timestamp);
    const o = {
      'M+': date.getMonth() + 1,
      'd+': date.getDate(),
      'h+': date.getHours() % 12,
      'H+': date.getHours(),
      'm+': date.getMinutes(),
      's+': date.getSeconds(),
      'q+': Math.floor((date.getMonth() + 3) / 3),
      S: date.getMilliseconds(),
    };
    if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (`${date.getFullYear()}`).substr(4 - RegExp.$1.length));
    for (const k in o) if (new RegExp(`(${k})`).test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : ((`00${o[k]}`).substr((`${o[k]}`).length)));
    return fmt;
  }
}

export default new Format();

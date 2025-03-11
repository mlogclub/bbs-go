/**
 * 滚动到顶部
 */
export function useScrollToTop(duration, elementId) {
  // 获取当前滚动条垂直位置
  let scrollTop = 0;
  let element = window;

  if (!elementId) {
    scrollTop =
      window.pageYOffset ||
      document.documentElement.scrollTop ||
      document.body.scrollTop;
  } else {
    element = document.getElementById(elementId);
    scrollTop = element.scrollTop || 0;
  }
  const start = scrollTop;

  // 计算滚动条滚动的距离
  const distance = -scrollTop;
  let startTime = null;

  function animation(currentTime) {
    if (startTime === null) {
      startTime = currentTime;
    }
    const timeElapsed = currentTime - startTime;
    const run = ease(timeElapsed, start, distance, duration);
    element.scrollTo(0, run);
    if (timeElapsed < duration) {
      requestAnimationFrame(animation);
    }
  }

  // 缓动函数
  function ease(t, b, c, d) {
    t /= d / 2;
    if (t < 1) return (c / 2) * t * t + b;
    t--;
    return (-c / 2) * (t * (t - 2) - 1) + b;
  }
  requestAnimationFrame(animation);
}

/**
 * 判断元素是否在可视视口内
 * @param {HTMLElement} el html元素
 * @returns
 */
export function useCheckInView(el) {
  if (el) {
    // 获取元素的视口
    const rect = el.getBoundingClientRect();
    return (
      rect.top < window.innerHeight &&
      rect.bottom > 0 &&
      rect.left < window.innerWidth &&
      rect.right > 0
    );
  }
  return false;
}

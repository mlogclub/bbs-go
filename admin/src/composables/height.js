export function useTableHeight() {
  const instance = getCurrentInstance();

  if (!instance) {
    return;
  }

  const containerHeaderRef =
    instance.ctx.$el.getElementsByClassName('container-header')[0];
  const containerMainRef =
    instance.ctx.$el.getElementsByClassName('container-main')[0];

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
    containerMainRef.style.height = `${Math.max(height, minHeight)}px`;
  }

  calcHeight();
  window.onresize = calcHeight;
}

export function empty(instance) {
  console.log('unsupported');
}

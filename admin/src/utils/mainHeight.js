/**
 * 计算页面mainContent内容高度
 * @param that 页面实例
 * @param handle 用于其他处理，这个函数中可以改变页面高度
 */
export default function calcMainHeight(that, handle) {
  calc();
  window.onresize = calc;

  function calc() {
    const magicHeight = 156; // 需要减去的其他高度，这些高度可能是一些边边角角的margin/padding
    const minHeight = 300;
    let height = document.documentElement.clientHeight - magicHeight;
    if (that.$refs.toolbar) {
      height -= that.$refs.toolbar.clientHeight;
    }
    if (that.$refs.pagebar) {
      height -= that.$refs.pagebar.clientHeight;
    }
    if (handle) {
      height = handle(height);
    }
    that.mainHeight = `${Math.max(height, minHeight)}px`;
  }
}

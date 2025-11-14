import { ElMessage, ElMessageBox, ElLoading } from "element-plus";

export function useMsg({ type = "success", message, duration = 800, onClose }) {
  ElMessage({
    duration: duration,
    type,
    message,
    onClose,
  });
}

export function useMsgSuccess(content) {
  ElMessage.success(content);
}

export function useMsgError(content) {
  ElMessage.error(content);
}

export function useMsgWarning(content) {
  ElMessage.warning(content);
}

export function useConfirm(content, options = {}) {
  // 默认选项，如果调用时没有提供 t 函数，则使用这些默认的按钮文本
  const defaultOptions = {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  };
  const mergedOptions = { ...defaultOptions, ...options };
  return ElMessageBox.confirm(content, mergedOptions.confirmButtonText, mergedOptions);
}

export function useLoading(text) {
  return ElLoading.service({
    lock: true,
    text: text || "Loading",
    background: "rgba(0, 0, 0, 0.7)",
  });
}
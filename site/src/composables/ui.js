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

export function useConfirm(content) {
  return ElMessageBox.confirm(content);
}

export function useLoading(text) {
  return ElLoading.service({
    lock: true,
    text: text || "Loading",
    background: "rgba(0, 0, 0, 0.7)",
  });
}

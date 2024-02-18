import { Notification, NotificationConfig } from '@arco-design/web-vue';
import { AppContext } from 'vue';

export default {
  Notification,
};

export function useNotificationInfo(
  config: string | NotificationConfig,
  appContext?: AppContext | undefined
) {
  Notification.info(config, appContext);
}

export function useNotificationSuccess(
  config: string | NotificationConfig,
  appContext?: AppContext | undefined
) {
  Notification.success(config, appContext);
}

export function useNotificationError(
  config: string | NotificationConfig,
  appContext?: AppContext | undefined
) {
  Notification.error(config, appContext);
}

export function useNotificationWarning(
  config: string | NotificationConfig,
  appContext?: AppContext | undefined
) {
  Notification.warning(config, appContext);
}

export function useHandleError(e: any) {
  useNotificationError(e.message || e);
}

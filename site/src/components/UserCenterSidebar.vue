<template>
  <div class="left-container">
    <my-counts :user="localUser" />
    <my-profile :user="localUser" />
    <fans-widget :user="localUser" />
    <follow-widget :user="localUser" />

    <div v-if="isAdmin" class="widget">
      <div class="widget-header">
        {{ t("component.userCenterSidebar.operations") }}
      </div>
      <div class="widget-content">
        <ul class="operations">
          <li v-if="localUser.forbidden">
            <i class="iconfont icon-forbidden" />
            <a @click="removeForbidden"
              >&nbsp;{{ t("component.userCenterSidebar.removeForbidden") }}</a
            >
          </li>
          <template v-else>
            <li>
              <i class="iconfont icon-forbidden" />
              <a @click="forbidden(7)"
                >&nbsp;{{ t("component.userCenterSidebar.forbidden7Days") }}</a
              >
            </li>
            <li>
              <i v-if="isSiteOwner" class="iconfont icon-forbidden" />
              <a @click="forbidden(-1)"
                >&nbsp;{{
                  t("component.userCenterSidebar.forbiddenForever")
                }}</a
              >
            </li>
          </template>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ElMessageBox } from "element-plus";
const { t } = useI18n();
const userStore = useUserStore();
const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
});
const localUser = ref(props.user);

const isSiteOwner = computed(() => {
  return userIsOwner(userStore.user);
});

const isAdmin = computed(() => {
  return userIsOwner(userStore.user) || userIsAdmin(userStore.user);
});

const forbidden = (days) => {
  const msg =
    days > 0
      ? t("component.userCenterSidebar.confirmForbidden")
      : t("component.userCenterSidebar.confirmForbiddenForever");
  ElMessageBox.confirm(msg, t("component.userCenterSidebar.dialogTitle"), {
    confirmButtonText: t("component.userCenterSidebar.confirmButtonText"),
    cancelButtonText: t("component.userCenterSidebar.cancelButtonText"),
    type: "warning",
  })
    .then(() => {
      doForbidden(days);
    })
    .catch(() => {});
};

const doForbidden = async (days) => {
  try {
    await useHttpPost(
      "/api/user/forbidden",
      useJsonToForm({
        userId: localUser.value.id,
        days,
      })
    );
    localUser.value.forbidden = true;
    useMsgSuccess(t("component.userCenterSidebar.forbiddenSuccess"));
  } catch (e) {
    useMsgError(t("component.userCenterSidebar.forbiddenFailed"));
  }
};

const removeForbidden = async () => {
  try {
    await useHttpPost(
      "/api/user/forbidden",
      useJsonToForm({
        userId: localUser.value.id,
        days: 0,
      })
    );
    localUser.value.forbidden = false;
    useMsgSuccess(t("component.userCenterSidebar.removeForbiddenSuccess"));
  } catch (e) {
    useMsgError(t("component.userCenterSidebar.removeForbiddenFailed"));
  }
};
</script>

<style lang="scss" scoped>
.img-avatar {
  margin-top: 5px;
  border: 1px dotted var(--border-color);
  border-radius: 5%;
}

.operations {
  list-style: none;

  li {
    padding-left: 3px;

    font-size: 13px;
    &:hover {
      cursor: pointer;
      background-color: #fcf8e3;
      color: #8a6d3b;
      font-weight: bold;
    }

    a {
      color: var(--text-link-color);
    }
  }
}
</style>

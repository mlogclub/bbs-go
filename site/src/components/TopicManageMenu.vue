<template>
  <el-dropdown
    v-if="menus && menus.length"
    trigger="click"
    @command="handleCommand"
  >
    <span class="el-dropdown-link">{{ t("component.topicManageMenu.manage") }}</span>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="item in menus"
          :key="item.command"
          :command="item.command"
          >{{ item.label }}</el-dropdown-item
        >
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
const { t } = useI18n();

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
});

const topic = ref(props.modelValue);

const emits = defineEmits(["update:modelValue"]);

const userStore = useUserStore();
const isOwner = userIsOwner(userStore.user);
const isAdmin = userIsAdmin(userStore.user);

const menus = computed(() => {
  const isTopicOwner =
    userStore.user && userStore.user.id === topic.value.user.id;
  const items = [];
  if (isTopicOwner && topic.value.type === 0) {
    items.push({
      command: "edit",
      label: t("component.topicManageMenu.edit"),
    });
  }
  if (isTopicOwner || isOwner || isAdmin) {
    items.push({
      command: "delete",
      label: t("component.topicManageMenu.delete"),
    });
  }
  if (isOwner || isAdmin) {
    items.push({
      command: "recommend",
      label: topic.value.recommend 
        ? t("component.topicManageMenu.cancelRecommend") 
        : t("component.topicManageMenu.recommend"),
    });
  }
  if (isOwner || isAdmin) {
    items.push({
      command: "sticky",
      label: topic.value.sticky 
        ? t("component.topicManageMenu.cancelSticky") 
        : t("component.topicManageMenu.sticky"),
    });
  }
  if (isOwner || isAdmin) {
    items.push({
      command: "forbidden7Days",
      label: t("component.topicManageMenu.forbidden7Days"),
    });
  }
  if (isOwner) {
    items.push({
      command: "forbiddenForever",
      label: t("component.topicManageMenu.forbiddenForever"),
    });
  }
  return items;
});

const handleCommand = async (command) => {
  if (command === "edit") {
    editTopic();
  } else if (command === "delete") {
    deleteTopic();
  } else if (command === "recommend") {
    switchRecommend();
  } else if (command === "sticky") {
    switchSticky();
  } else if (command === "forbidden7Days") {
    await forbidden(7);
  } else if (command === "forbiddenForever") {
    await forbidden(-1);
  } else {
    console.log("click on item " + command);
  }
};

const forbidden = async (days) => {
  try {
    await useHttpPost(
      "/api/user/forbidden",
      useJsonToForm({
        userId: topic.value.user.id,
        days,
      })
    );
    useMsgSuccess(t("component.topicManageMenu.forbiddenSuccess"));
  } catch (e) {
    useMsgError(t("component.topicManageMenu.forbiddenFailed"));
  }
};

const deleteTopic = () => {
  useConfirm(t("component.topicManageMenu.confirmDelete")).then(function () {
    useHttpPost(`/api/topic/delete/${topic.value.id}`)
      .then(() => {
        useMsg({
          message: t("component.topicManageMenu.deleteSuccess"),
          onClose() {
            useLinkTo("/topics");
          },
        });
      })
      .catch((e) => {
        useMsgError(t("component.topicManageMenu.deleteFailed") + (e.message || e));
      });
  });
};

const editTopic = () => {
  useLinkTo(`/topic/edit/${topic.value.id}`);
};

const switchRecommend = () => {
  const action = topic.value.recommend 
    ? t("component.topicManageMenu.cancelRecommend") 
    : t("component.topicManageMenu.recommend");
  useConfirm(t("component.topicManageMenu.confirmAction", { action })).then(function () {
    const recommend = !topic.value.recommend;
    useHttpPost(
      `/api/topic/recommend/${topic.value.id}`,
      useJsonToForm({
        recommend,
      })
    )
      .then(() => {
        topic.value.recommend = recommend;
        emits("update:modelValue", topic.value);
        useMsgSuccess({
          message: t("component.topicManageMenu.actionSuccess", { action }),
        });
      })
      .catch((e) => {
        useMsgError(t("component.topicManageMenu.actionFailed", { action }) + (e.message || e));
      });
  });
};

const switchSticky = () => {
  const action = topic.value.sticky 
    ? t("component.topicManageMenu.cancelSticky") 
    : t("component.topicManageMenu.sticky");
  useConfirm(t("component.topicManageMenu.confirmAction", { action })).then(function () {
    const sticky = !topic.value.sticky;
    useHttpPost(
      `/api/topic/sticky/${topic.value.id}`,
      useJsonToForm({
        sticky,
      })
    )
      .then(() => {
        topic.value.sticky = sticky;
        emits("update:modelValue", topic.value);
        useMsgSuccess({
          message: t("component.topicManageMenu.actionSuccess", { action }),
        });
      })
      .catch((e) => {
        useMsgError(t("component.topicManageMenu.actionFailed", { action }) + (e.message || e));
      });
  });
}
</script>

<style lang="scss" scoped>
.el-dropdown-link {
  cursor: pointer;
  color: var(--text-color3);
  font-size: 12px;
}
.el-dropdown-menu__item {
  font-size: 12px;
}
</style>

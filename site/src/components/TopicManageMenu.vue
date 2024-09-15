<template>
  <!-- <ClientOnly> -->
  <el-dropdown
    v-if="menus && menus.length"
    trigger="click"
    @command="handleCommand"
  >
    <span class="el-dropdown-link">管理</span>
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
  <!-- </ClientOnly> -->
</template>

<script setup>
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
      label: "修改",
    });
  }
  if (isTopicOwner || isOwner || isAdmin) {
    items.push({
      command: "delete",
      label: "删除",
    });
  }
  if (isOwner || isAdmin) {
    items.push({
      command: "recommend",
      label: topic.value.recommend ? "取消推荐" : "推荐",
    });
  }
  if (isOwner || isAdmin) {
    items.push({
      command: "sticky",
      label: topic.value.sticky ? "取消置顶" : "置顶",
    });
  }
  if (isOwner || isAdmin) {
    items.push({
      command: "forbidden7Days",
      label: "禁言7天",
    });
  }
  if (isOwner) {
    items.push({
      command: "forbiddenForever",
      label: "永久禁言",
    });
  }
  return items;
});

async function handleCommand(command) {
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
}
async function forbidden(days) {
  try {
    await useHttpPostForm("/api/user/forbidden", {
      body: {
        userId: topic.value.user.id,
        days,
      },
    });
    useMsgSuccess("禁言成功");
  } catch (e) {
    useMsgError("禁言失败");
  }
}
function deleteTopic() {
  useConfirm("是否确认删除该帖子？").then(function () {
    useHttpPost(`/api/topic/delete/${topic.value.id}`)
      .then(() => {
        useMsg({
          message: "删除成功",
          onClose() {
            useLinkTo("/topics");
          },
        });
      })
      .catch((e) => {
        useMsgError("删除失败：" + (e.message || e));
      });
  });
}
function editTopic() {
  useLinkTo(`/topic/edit/${topic.value.id}`);
}
function switchRecommend() {
  const action = topic.value.recommend ? "取消推荐" : "推荐";
  useConfirm(`是否确认${action}该帖子？`).then(function () {
    const recommend = !topic.value.recommend;
    useHttpPostForm(`/api/topic/recommend/${topic.value.id}`, {
      body: {
        recommend,
      },
    })
      .then(() => {
        topic.value.recommend = recommend;
        emits("update:modelValue", topic.value);
        useMsgSuccess({
          message: `${action}成功`,
        });
      })
      .catch((e) => {
        useMsgError(`${action}失败：` + (e.message || e));
      });
  });
}
function switchSticky() {
  const action = topic.value.sticky ? "取消置顶" : "置顶";
  useConfirm(`是否确认${action}该帖子？`).then(function () {
    const sticky = !topic.value.sticky;
    useHttpPostForm(`/api/topic/sticky/${topic.value.id}`, {
      body: {
        sticky,
      },
    })
      .then(() => {
        topic.value.sticky = sticky;
        emits("update:modelValue", topic.value);
        useMsgSuccess({
          message: `${action}成功`,
        });
      })
      .catch((e) => {
        useMsgError(`${action}失败：` + (e.message || e));
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

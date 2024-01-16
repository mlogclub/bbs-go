<template>
  <ClientOnly>
    <el-dropdown v-if="hasPermission" trigger="click" @command="handleCommand">
      <span class="el-dropdown-link">管理</span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item v-if="topic.type === 0" command="edit"
            >修改</el-dropdown-item
          >
          <el-dropdown-item command="delete">删除</el-dropdown-item>
          <el-dropdown-item v-if="isOwner || isAdmin" command="recommend">{{
            topic.recommend ? "取消推荐" : "推荐"
          }}</el-dropdown-item>
          <el-dropdown-item v-if="isOwner || isAdmin" command="sticky">{{
            topic.sticky ? "取消置顶" : "置顶"
          }}</el-dropdown-item>
          <el-dropdown-item v-if="isOwner || isAdmin" command="forbidden7Days"
            >禁言7天</el-dropdown-item
          >
          <el-dropdown-item v-if="isOwner" command="forbiddenForever"
            >永久禁言</el-dropdown-item
          >
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </ClientOnly>
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
const isTopicOwner = computed(() => {
  return userStore.user && userStore.user.id === topic.value.user.id;
});
const hasPermission = computed(() => {
  return isTopicOwner || isOwner || isAdmin;
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

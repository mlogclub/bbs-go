<template>
  <ClientOnly>
    <el-dropdown v-if="hasPermission" trigger="click" @command="handleCommand">
      <span class="el-dropdown-link">管理</span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="edit">修改</el-dropdown-item>
          <el-dropdown-item command="delete">删除</el-dropdown-item>
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
  article: {
    type: Object,
    required: true,
  },
});

const userStore = useUserStore();
const isOwner = userIsOwner(userStore.user);
const isAdmin = userIsAdmin(userStore.user);
const isArticleOwner = computed(() => {
  return userStore.user && userStore.user.id === props.article.user.id;
});
const hasPermission = computed(() => {
  return isArticleOwner || isOwner || isAdmin;
});

async function handleCommand(command) {
  if (command === "edit") {
    editArticle();
  } else if (command === "delete") {
    deleteArticle();
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
        userId: props.article.user.id,
        days,
      },
    });
    useMsgSuccess("禁言成功");
  } catch (e) {
    useMsgError("禁言失败");
  }
}
function deleteArticle() {
  useConfirm("是否确认删除该文章？").then(function () {
    useHttpPost(`/api/article/delete/${props.article.id}`)
      .then(() => {
        useMsg({
          message: "删除成功",
          onClose() {
            useLinkTo("/articles");
          },
        });
      })
      .catch((e) => {
        useMsgError("删除失败：" + (e.message || e));
      });
  });
}
function editArticle() {
  useLinkTo(`/article/edit/${props.article.id}`);
}
</script>

<style lang="scss" scoped>
.el-dropdown-link {
  cursor: pointer;
  color: var(--text-color3);
  font-size: 12px;
}
</style>

<template>
  <ClientOnly>
    <el-dropdown v-if="showDropdown" trigger="click" @command="handleCommand">
      <span class="el-dropdown-link">{{$t("component.articleManageMenu.manage")}}</span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item v-if="canEdit" command="edit">{{$t("component.articleManageMenu.edit")}}</el-dropdown-item>
          <el-dropdown-item v-if="canDelete" command="delete">{{$t("component.articleManageMenu.delete")}}</el-dropdown-item>
          <el-dropdown-item v-if="canForbid" command="forbidden7Days">{{$t("component.articleManageMenu.forbidden7Days")}}</el-dropdown-item>
          <el-dropdown-item v-if="canForbid" command="forbiddenForever">{{$t("component.articleManageMenu.forbiddenForever")}}</el-dropdown-item>
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

const { t } = useI18n();
const userStore = useUserStore();
const isOwner = computed(() => userIsOwner(userStore.user));
const isAdmin = computed(() => userIsAdmin(userStore.user));
const isArticleOwner = computed(() => {
  return userStore.user && userStore.user.id === props.article.user.id;
});
const canEdit = computed(() => isArticleOwner.value || isOwner.value || isAdmin.value);
const canDelete = computed(() => isArticleOwner.value || isOwner.value || isAdmin.value);
const canForbid = computed(() => isOwner.value || isAdmin.value);
const showDropdown = computed(() => {
  return canEdit.value || canDelete.value || canForbid.value;
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
    await useHttpPost(
      "/api/user/forbidden",
      useJsonToForm({
        userId: props.article.user.id,
        days,
      })
    );
    useMsgSuccess(t("component.articleManageMenu.forbiddenSuccess"));
  } catch (e) {
    useMsgError(t("component.articleManageMenu.forbiddenFailed"));
  }
}
function deleteArticle() {
  useConfirm(t("component.articleManageMenu.confirmDelete")).then(function () {
    useHttpPost(`/api/article/delete/${props.article.id}`)
      .then(() => {
        useMsg({
          message: t("component.articleManageMenu.deleteSuccess"),
          onClose() {
            useLinkTo("/articles");
          },
        });
      })
      .catch((e) => {
        useMsgError(
          t("component.articleManageMenu.deleteFailed") + (e.message || e)
        );
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

<template>
  <section class="main">
    <div class="container">
      <article v-if="isNeedEmailVerify" class="message is-warning">
        <div class="message-header">
          <p>请先验证邮箱</p>
        </div>
        <div class="message-body">
          发表话题前，请先前往
          <strong
            ><nuxt-link
              to="/user/profile/account"
              style="color: var(--text-link-color)"
              >个人中心 &gt; 账号设置</nuxt-link
            ></strong
          >
          页面设置邮箱，并完成邮箱认证。
        </div>
      </article>
      <div v-else class="publish-form">
        <div class="form-title">
          <div class="form-title-name">发文章</div>
        </div>

        <div class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input"
              type="text"
              placeholder="标题"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <markdown-editor
              v-model="postForm.content"
              placeholder="请输入内容，将图片复制或拖入编辑器可上传"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <image-upload v-model="postForm.cover" :limit="1" size="120px">
              <template #add-image-button>
                <div class="cover-add-btn">
                  <i class="iconfont icon-add" />
                  <span>封面</span>
                </div>
              </template>
            </image-upload>
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <a
              v-if="publishing"
              :class="{ 'is-loading': publishing }"
              disabled
              class="button is-primary"
              >发表</a
            >
            <a
              v-else
              :class="{ 'is-loading': publishing }"
              class="button is-primary"
              @click="submitCreate"
              >发表</a
            >
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const publishing = ref(false); // 当前是否正处于发布中...
const postForm = ref({
  title: "",
  tags: [],
  content: "",
});

const userStore = useUserStore();
const configStore = useConfigStore();
const isNeedEmailVerify = computed(() => {
  return (
    configStore.config.createArticleEmailVerified &&
    !userStore.user.emailVerified
  );
});

if (!configStore.config.modules.article) {
  showError("😱 文章功能未开启");
}

useHead({
  title: useSiteTitle("发表文章"),
});

definePageMeta({
  middleware: "auth",
});

if (!configStore.isEnabledArticle) {
  throw createError({
    statusCode: 403,
    message: "已关闭文章功能",
  });
}

async function submitCreate() {
  if (publishing.value) {
    return;
  }
  publishing.value = true;
  try {
    const article = await useHttpPost(
      "/api/article/create",
      useJsonToForm({
        title: postForm.value.title,
        content: postForm.value.content,
        tags: postForm.value.tags ? postForm.value.tags.join(",") : "",
        cover:
          postForm.value.cover && postForm.value.cover.length
            ? JSON.stringify(postForm.value.cover[0])
            : null,
      })
    );
    useMsg({
      message: "提交成功",
      onClose() {
        useLinkTo(`/article/${article.id}`);
      },
    });
  } catch (e) {
    useMsgError(e.message || e);
    publishing.value = false;
  }
}
</script>

<style lang="scss" scoped>
.cover-add-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  i {
    font-size: 18px;
    color: var(--text-color3);
  }

  span {
    font-size: 14px;
    color: var(--text-color3);
    font-weight: 500;
  }
}
</style>

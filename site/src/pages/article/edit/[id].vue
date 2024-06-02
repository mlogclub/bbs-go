<template>
  <section class="main">
    <div class="container" v-if="postForm">
      <div class="article-create-form">
        <h1 class="title">修改文章</h1>

        <div class="field">
          <div class="control">
            <div
              v-for="node in nodes"
              :key="node.id"
              class="article-tag"
              :class="{ selected: postForm.nodeId === node.id }"
              @click="postForm.nodeId = node.id"
            >
              <span>{{ node.name }}</span>
            </div>
          </div>
        </div>

        <div class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input article-title"
              type="text"
              placeholder="请输入文章标题"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <markdown-editor
              v-model="postForm.content"
              placeholder="可空，将图片复制或拖入编辑器可上传"
            />
          </div>
        </div>

        <div v-if="isEnableHideContent || postForm.hideContent" class="field">
          <div class="control">
            <markdown-editor
              ref="mdEditor"
              v-model="postForm.hideContent"
              height="200px"
              placeholder="隐藏内容，评论后可见"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <a
              v-if="publishing"
              :class="{ 'is-loading': publishing }"
              disabled
              class="button is-success"
              >提交更改</a
            >
            <a
              v-else
              :class="{ 'is-loading': publishing }"
              class="button is-success"
              @click="submitCreate"
              >提交更改</a
            >
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
definePageMeta({
  middleware: "auth",
});
useHead({
  title: useSiteTitle("修改文章"),
});

const route = useRoute();
const configStore = useConfigStore();

const isEnableHideContent = computed(() => {
  return configStore.config.enableHideContent;
});

const { data: nodes } = useAsyncData("nodes", () =>
  useMyFetch("/api/article/nodes")
);
const { data: postForm } = useAsyncData(() =>
  useMyFetch(`/api/article/edit/${route.params.id}`)
);
const publishing = ref(false);

async function submitCreate() {
  if (publishing.value) {
    return;
  }
  publishing.value = true;

  try {
    useHttpPostForm(`/api/article/edit/${postForm.value.id}`, {
      body: {
        title: postForm.value.title,
        content: postForm.value.content,
        cover: postForm.value.cover,
        tags: postForm.value.tags ? postForm.value.tags.join(",") : "",
      },
    });
    useMsg({
      message: "修改成功",
      onClose() {
        useLinkTo(`/article/${postForm.value.id}`);
      },
    });
  } catch (e) {
    publishing.value = false;
    useMsgError("提交失败：" + (e.message || e));
  }
}
</script>

<style lang="scss" scoped>
.article-create-form {
  background-color: var(--bg-color);
  padding: 30px;

  .article-form-title {
    font-size: 36px;
    font-weight: 500;
    margin-bottom: 10px;
  }
  .field {
    margin-bottom: 10px;

    input {
      &:focus-visible {
        outline-width: 0;
      }
    }
  }
}

.cover-add-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  i {
    font-size: 24px;
    color: #1878f3;
  }

  span {
    font-size: 14px;
    color: var(--text-color3);
  }
}
</style>

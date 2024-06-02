<template>
  <section class="main">
    <div class="container">
      <div class="topic-create-form" v-if="postForm">
        <div class="topic-form-title">修改话题</div>

        <div class="field">
          <div class="control">
            <div
              v-for="node in nodes"
              :key="node.id"
              class="topic-tag"
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
              class="input topic-title"
              type="text"
              placeholder="请输入帖子标题"
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
  title: useSiteTitle("修改话题"),
});

const route = useRoute();
const configStore = useConfigStore();

const isEnableHideContent = computed(() => {
  return configStore.config.enableHideContent;
});

const { data: nodes } = useAsyncData("nodes", () =>
  useMyFetch("/api/topic/nodes")
);
const { data: postForm } = useAsyncData(() =>
  useMyFetch(`/api/topic/edit/${route.params.id}`)
);
const publishing = ref(false);

async function submitCreate() {
  if (publishing.value) {
    return;
  }
  publishing.value = true;

  try {
    useHttpPostForm(`/api/topic/edit/${postForm.value.id}`, {
      body: {
        nodeId: postForm.value.nodeId,
        title: postForm.value.title,
        content: postForm.value.content,
        hideContent: postForm.value.hideContent,
        tags: postForm.value.tags ? postForm.value.tags.join(",") : "",
      },
    });
    useMsg({
      message: "修改成功",
      onClose() {
        useLinkTo(`/topic/${postForm.value.id}`);
      },
    });
  } catch (e) {
    publishing.value = false;
    useMsgError("提交失败：" + (e.message || e));
  }
}
</script>

<style lang="scss" scoped></style>

<template>
  <section class="main">
    <div class="container">
      <div class="publish-form" v-if="postForm">
        <div class="form-title">
          <div class="form-title-name">修改帖子</div>
          <div class="form-title-switch" @click="switchEditor">
            <div v-if="postForm.contentType === 'markdown'" class="editor-type">
              <img src="~/assets/images/markdown.svg" />
              <span>Markdown</span>
            </div>
            <div v-else class="editor-type">
              <img src="~/assets/images/html.svg" />
              <span>HTML</span>
            </div>
            <i class="iconfont icon-switch" />
          </div>
        </div>

        <div class="topic-tags">
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
              v-if="postForm.contentType === 'markdown'"
              v-model="postForm.content"
              placeholder="请输入你要发表的内容..."
            />
            <MEditor
              v-else
              v-model="postForm.content"
              :uploadImage="useUploadImage"
            />
          </div>
        </div>

        <div v-if="isEnableHideContent || postForm.hideContent" class="field">
          <div class="control">
            <MEditor v-model="postForm.hideContent" height="200px" />
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
              class="button is-primary"
              >提交更改</a
            >
            <a
              v-else
              :class="{ 'is-loading': publishing }"
              class="button is-primary"
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

const { data: nodes } = await useMyFetch("/api/topic/nodes");
const { data: postForm } = await useMyFetch(
  `/api/topic/edit/${route.params.id}`
);
const publishing = ref(false);

const switchEditor = () => {
  useConfirm("切换编辑器将会清空当前内容，是否继续？")
    .then(() => {
      postForm.value.content = "";
      if (postForm.value.contentType === "markdown") {
        postForm.value.contentType = "html";
      } else {
        postForm.value.contentType = "markdown";
      }
    })
    .catch(() => {});
};

async function submitCreate() {
  if (publishing.value) {
    return;
  }
  publishing.value = true;

  try {
    useHttpPost(
      `/api/topic/edit/${postForm.value.id}`,
      useJsonToForm({
        nodeId: postForm.value.nodeId,
        title: postForm.value.title,
        content: postForm.value.content,
        hideContent: postForm.value.hideContent,
        tags: postForm.value.tags ? postForm.value.tags.join(",") : "",
      })
    );
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

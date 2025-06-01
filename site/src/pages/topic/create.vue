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
          <div class="form-title-name">
            {{ postForm.type === 0 ? "发帖" : "发动态" }}
          </div>
          <div
            v-if="postForm.type === 0"
            class="form-title-switch"
            @click="switchEditor"
          >
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

        <div v-if="postForm.type === 0" class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input topic-title"
              type="text"
              placeholder="请输入帖子标题"
            />
          </div>
        </div>

        <div v-if="postForm.type === 0" class="field">
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

        <div v-if="postForm.type === 0 && isEnableHideContent" class="field">
          <div class="control">
            <MEditor v-model="postForm.hideContent" height="200px" />
          </div>
        </div>

        <div v-if="postForm.type === 1" class="field">
          <div class="control">
            <simple-editor
              ref="simpleEditorComponent"
              v-model:content="postForm.content"
              v-model:imageList="postForm.imageList"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="form-footer">
          <a
            :class="{ 'is-loading': publishing }"
            class="button is-primary btn-publish"
            @click="publish"
            >{{ postForm.type === 1 ? "发表动态" : "发表帖子" }}</a
          >
        </div>
      </div>
    </div>

    <CaptchaDialog ref="captchaDialog" />
  </section>
</template>

<script setup>
definePageMeta({
  middleware: "auth",
});

const userStore = useUserStore();
const configStore = useConfigStore();
const route = useRoute();
const router = useRouter();

const type = Number.parseInt(route.query.type) || 0;
const nodeId =
  parseInt(route.query.nodeId) || configStore.config.defaultNodeId || 0;

if (type === 1 && !configStore.config.modules.tweet) {
  showError("😱 动态功能未开启");
}
if (type === 0 && !configStore.config.modules.topic) {
  showError("😱 帖子功能未开启");
}

const postForm = ref({
  type: type,
  nodeId: nodeId,
  title: "",
  tags: [],
  contentType: "html",
  content: "",
  hideContent: "",
  imageList: [],

  captchaId: "",
  captchaCode: "",
  captchaProtocol: 2,
});
const publishing = ref(false);
const simpleEditorComponent = ref(null);
const captchaDialog = ref(null);

const isNeedEmailVerify = computed(() => {
  return (
    configStore.config.createTopicEmailVerified && !userStore.user.emailVerified
  );
});

const isEnableHideContent = computed(() => {
  return configStore.config.enableHideContent;
});

const topicCaptchaEnabled = computed(() => {
  return configStore.config.topicCaptcha;
});

const { data: nodes } = await useMyFetch("/api/topic/nodes");

watch(
  () => route.query,
  (newQuery, oldQuery) => {
    init();
  },
  { deep: true }
);

const init = () => {
  postForm.value.type = Number.parseInt(route.query.type) || 0;

  let contentType = route.query.contentType;
  if (!contentType) {
    contentType = "html";
  }
  postForm.value.contentType = "contentType";

  useHead({
    title: useSiteTitle(type === 0 ? "发帖子" : "发动态"),
  });
};

init();

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

const publish = () => {
  if (publishing.value) {
    return;
  }

  console.log(configStore.config);

  console.log(topicCaptchaEnabled.value);

  if (topicCaptchaEnabled.value) {
    captchaDialog.value.show().then((captcha) => {
      publishSubmit(captcha);
    });
  } else {
    publishSubmit();
  }
};

const publishSubmit = async (captcha) => {
  if (publishing.value) {
    return;
  }

  if (postForm.value.type === 1) {
    if (simpleEditorComponent.value.loading) {
      useMsgWarning("图片上传中,请稍后重试...");
      return;
    }
  }

  if (captcha) {
    postForm.value.captchaId = captcha.captchaId;
    postForm.value.captchaCode = captcha.captchaCode;
    postForm.value.captchaProtocol = 2;
  }

  publishing.value = true;
  try {
    const topic = await useHttpPost("/api/topic/create", postForm.value);
    router.push(`/topic/${topic.id}`);
  } catch (e) {
    useMsgError(e.message || e);
    publishing.value = false;
  }
};
</script>

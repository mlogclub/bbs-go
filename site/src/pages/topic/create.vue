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
          {{ postForm.type === 0 ? "发帖" : "发动态" }}
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
              v-model="postForm.content"
              placeholder="请输入你要发表的内容..."
            />
          </div>
        </div>

        <div v-if="postForm.type === 0 && isEnableHideContent" class="field">
          <div class="control">
            <markdown-editor
              v-model="postForm.hideContent"
              height="200px"
              placeholder="请输入隐藏内容，隐藏内容，评论后可见"
            />
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
            class="button is-success btn-publish"
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

const { data: nodes } = useAsyncData("nodes", () =>
  useHttpGet("/api/topic/nodes")
);

watch(
  () => route.query,
  (newQuery, oldQuery) => {
    init();
  },
  { deep: true }
);

const init = () => {
  postForm.value.type = Number.parseInt(route.query.type) || 0;
  useHead({
    title: useSiteTitle(type === 0 ? "发帖子" : "发动态"),
  });
};

init();

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
    const topic = await useHttpPost("/api/topic/create", {
      body: postForm.value,
    });
    router.push(`/topic/${topic.id}`);
  } catch (e) {
    useMsgError(e.message || e);
    publishing.value = false;
  }
};
</script>

<style lang="scss" scoped>
.publish-form {
  border-radius: var(--border-radius);
  background: var(--bg-color);
  padding: 10px 18px 18px 18px;

  .form-title {
    font-size: 18px;
    font-weight: 500;
    margin-bottom: 10px;
    padding-bottom: 10px;
    border-bottom: 1px solid var(--border-color);
  }

  .form-footer {
    text-align: right;

    .btn-publish {
      width: 130px;
      color: #fff;
    }
  }

  .field {
    margin-bottom: 10px;

    input {
      border: 1px solid var(--border-color);
      background-color: var(--bg-color);
      border-radius: 3px;

      &:focus-visible {
        outline-width: 0;
      }
      &:focus {
        box-shadow: none;
      }
    }
  }

  .topic-tags {
    margin-bottom: 10px;
    display: flex;
    gap: 10px;

    .topic-tag {
      cursor: pointer;
      padding: 0 12px;
      display: flex;
      justify-content: center;
      align-items: center;
      border-radius: 3px;
      background: var(--bg-color3);
      // border: 1px solid var(--border-color);
      color: var(--text-color3);
      font-size: 14px;
      line-height: 24px;

      &:hover {
        color: var(--text-link-color);
        background: var(--bg-color5);
        // border: 1px solid var(--border-hover-color);
      }

      &.selected {
        color: var(--text-color5);
        background: var(--text-link-color);
      }
    }
  }
}
</style>

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
      <div v-else class="topic-create-form">
        <div class="topic-form-title">
          {{ postForm.type === 0 ? "发帖子" : "发动态" }}
        </div>

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
              placeholder="隐藏内容，评论后可见"
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

        <div v-if="postForm.captchaUrl" class="field is-horizontal">
          <div class="field control has-icons-left">
            <input
              v-model="postForm.captchaCode"
              class="input"
              type="text"
              placeholder="验证码"
              style="max-width: 150px; margin-right: 20px"
            />
            <span class="icon is-small is-left">
              <i class="iconfont icon-captcha" />
            </span>
          </div>
          <div class="field">
            <a @click="showCaptcha">
              <img :src="postForm.captchaUrl" style="height: 40px" />
            </a>
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <a
              :class="{ 'is-loading': publishing }"
              class="button is-success"
              @click="createTopic"
              >{{ postForm.type === 1 ? "发表动态" : "发表帖子" }}</a
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

const userStore = useUserStore();
const configStore = useConfigStore();
const route = useRoute();
const router = useRouter();

const nodeId =
  parseInt(route.query.nodeId) || configStore.config.defaultNodeId || 0;

const postForm = ref({
  type: Number.parseInt(route.query.type) || 0,
  nodeId: nodeId,
  title: "",
  tags: [],
  content: "",
  hideContent: "",
  imageList: [],

  captchaId: "",
  captchaUrl: "",
  captchaCode: "",
});
const publishing = ref(false);
const simpleEditorComponent = ref(null);

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
  useMyFetch("/api/topic/nodes")
);

init();

watch(
  () => route.query,
  (newQuery, oldQuery) => {
    // eslint-disable-next-line no-console
    // console.log(newQuery, oldQuery);

    init();
  },
  { deep: true }
);

onMounted(() => {
  showCaptcha();
});

function init() {
  postForm.value.type = Number.parseInt(route.query.type) || 0;
  useHead({
    title: postForm.value.type === 0 ? "发帖子" : "发动态",
  });
}

async function createTopic() {
  if (publishing.value) {
    return;
  }

  if (postForm.value.type === 1) {
    if (simpleEditorComponent.value.loading) {
      useMsgWarning("图片上传中,请稍后重试...");
      return;
    }
  }

  publishing.value = true;
  try {
    const topic = await useHttpPost("/api/topic/create", {
      body: postForm.value,
    });
    router.push(`/topic/${topic.id}`);
  } catch (e) {
    showCaptcha();
    useMsgError(e.message || e);
    publishing.value = false;
  }
}

async function showCaptcha() {
  if (topicCaptchaEnabled.value) {
    try {
      const ret = await useHttpGet("/api/captcha/request", {
        params: {
          captchaId: postForm.value.captchaId || "",
        },
      });
      postForm.value.captchaId = ret.captchaId;
      postForm.value.captchaUrl = ret.captchaUrl;
    } catch (e) {
      useMsgError(e.message || e);
    }
  }
}
</script>

<style lang="scss" scoped></style>

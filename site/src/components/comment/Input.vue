<template>
  <div class="comment-form">
    <div class="comment-create">
      <div ref="commentEditor" class="comment-input-wrapper">
        <div v-if="quote" class="comment-quote-info">
          {{ t('component.comment.input.replyTo') }}
          <label v-text="quote.user.nickname" />
          <i class="iconfont icon-close" :alt="t('component.comment.input.cancelReply')" @click="cancelReply" />
        </div>
        <text-editor
          ref="textEditorRef"
          v-model:content="content"
          v-model:imageList="imageList"
          :height="90"
          :focus-height="120"
          @submit="create"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
const { t } = useI18n();

const props = defineProps({
  entityType: {
    type: String,
    default: "",
    required: true,
  },
  entityId: {
    type: Number,
    default: 0,
    required: true,
  },
});

const emits = defineEmits(["created"]);

const textEditorRef = ref(null);
const content = ref("");
const imageList = ref(null);
const sending = ref(false); // 发送中
const quote = ref(null); // 引用的对象
const commentEditor = ref(null); // 编辑器组件

async function create() {
  if (!content.value) {
    useMsgError(t('component.comment.input.pleaseInput'));
    return;
  }
  if (sending.value) {
    return;
  }
  sending.value = true;
  try {
    const data = await useHttpPost(
      "/api/comment/create",
      useJsonToForm({
        contentType: props.contentType,
        entityType: props.entityType,
        entityId: props.entityId,
        content: content.value,
        imageList:
          imageList.value && imageList.value.length
            ? JSON.stringify(imageList.value)
            : "",
        quoteId: quote.value ? quote.value.id : "",
      })
    );
    emits("created", data);

    textEditorRef.value.reset();
    content.value = "";
    imageList.value = [];
    quote.value = null;
    useMsgSuccess(t('component.comment.input.publishSuccess'));
  } catch (e) {
    console.error(e);
    useMsgError(e.message || e);
  } finally {
    sending.value = false;
  }
}
function reply(quote) {
  quote.value = quote;
  commentEditor.value.scrollIntoView({
    block: "start",
    behavior: "smooth",
  });
}
function cancelReply() {
  quote.value = null;
}
</script>

<style scoped lang="scss">
.comment-form {
  background-color: var(--bg-color);
  margin: 10px 0;

  .comment-create {
    // border-radius: 4px;
    overflow: hidden;
    position: relative;
    padding: 0;
    box-sizing: border-box;

    .comment-quote-info {
      font-size: 13px;
      color: var(--text-color);
      margin-bottom: 3px;
      font-weight: 600;

      i {
        font-size: 12px !important;
        color: var(--text-link-color);
        cursor: pointer;
      }

      i:hover {
        color: red;
      }
    }

    .comment-input-wrapper {
      .text-input {
        outline: none;
        width: 100%;
        height: 85px;
        font-size: 14px;
        padding: 10px 40px 10px 10px;
        color: var(--text-color);
        line-height: 16px;
        max-width: 100%;
        resize: none;
        border: 1px solid var(--border-color);
        box-sizing: border-box;
        border-radius: var(--jinsom-border-radius);
      }
    }

    .comment-button-wrapper {
      user-select: none;
      display: flex;
      float: right;
      height: 30px;
      line-height: 30px;

      span {
        color: var(--text-color4);
        font-size: 13px;
        margin-right: 5px;
      }

      button {
        font-weight: 500;
      }
    }
  }
}
</style>

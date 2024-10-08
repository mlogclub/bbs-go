<script setup>
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
const value = ref({
  content: "", // 内容
  imageList: [],
});
const sending = ref(false); // 发送中
const quote = ref(null); // 引用的对象
const commentEditor = ref(null); // 编辑器组件
const simpleEditor = ref(null); // 编辑器组件

async function create() {
  if (!value.value.content) {
    useMsgError("请输入评论内容");
    return;
  }
  if (sending.value) {
    return;
  }
  if (simpleEditor.value && simpleEditor.value.isOnUpload()) {
    useMsgWarning("正在上传中...请上传完成后提交");
    return;
  }
  sending.value = true;
  try {
    const data = await useHttpPostForm("/api/comment/create", {
      body: {
        contentType: props.contentType,
        entityType: props.entityType,
        entityId: props.entityId,
        content: value.value.content,
        imageList:
          value.value.imageList && value.value.imageList.length
            ? JSON.stringify(value.value.imageList)
            : "",
        quoteId: quote.value ? quote.value.id : "",
      },
    });
    emits("created", data);

    value.value.content = "";
    value.value.imageList = [];
    quote.value = null;
    simpleEditor.value.clear();
    useMsgSuccess("发布成功");
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

<template>
  <div class="comment-form">
    <div class="comment-create">
      <div ref="commentEditor" class="comment-input-wrapper">
        <div v-if="quote" class="comment-quote-info">
          回复：
          <label v-text="quote.user.nickname" />
          <i class="iconfont icon-close" alt="取消回复" @click="cancelReply" />
        </div>
        <text-editor ref="simpleEditor" v-model="value" @submit="create" />
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.comment-form {
  background-color: var(--bg-color);
  padding: 10px;
  margin-bottom: 10px;

  .comment-create {
    border-radius: 4px;
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
      margin-bottom: 8px;

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

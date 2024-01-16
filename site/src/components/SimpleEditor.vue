<script setup>
const props = defineProps({
  maxWordCount: {
    type: Number,
    default: 5000,
  },
  modelValue: {
    type: Object,
    default() {
      return {
        content: '',
        imageList: [],
      }
    },
  },
  height: {
    type: String,
    default: '200px',
  },
})
const emits = defineEmits(['update:modelValue', 'update:content', 'update:imageList'])

const post = ref({
  content: '',
  imageList: ref([]),
})
const showImageUpload = ref(false)
const imageUploaderComponent = ref(null)

const wordCount = computed(() => {
  return post.value.content ? post.value.content.length : 0
})
const loading = computed(() => {
  if (imageUploaderComponent.value) {
    return imageUploaderComponent.value.loading
  }
  return false
})

function switchImageUpload() {
  if (!showImageUpload.value) {
    // 打开文件弹窗
    imageUploaderComponent.value.onClick()
  }
  showImageUpload.value = !showImageUpload.value
}

function onContentChange() {
  emits('update:content', post.value.content)
  emits('update:modelValue', post.value)
}

function onImageListChange() {
  emits('update:imageList', post.value.imageList)
  emits('update:modelValue', post.value)
}

defineExpose({
  loading,
})
</script>

<template>
  <div class="simple-editor">
    <div class="simple-editor-toolbar">
      <div class="act-btn">
        <i class="iconfont icon-image" @click="switchImageUpload" />
      </div>
      <div class="publish-container">
        <span class="tip">{{ wordCount }} / {{ maxWordCount }}</span>
      </div>
    </div>
    <label class="simple-editor-input">
      <textarea
        v-model="post.content"
        placeholder="请输入您要发表的内容 ..."
        :style="{ 'min-height': height, 'height': height }"
        @update:model-value="onContentChange"
        @paste="handleParse"
        @drop="handleDrag"
        @keydown.ctrl.enter="doSubmit"
        @keydown.meta.enter="doSubmit"
      />
    </label>
    <div v-show="showImageUpload" class="simple-editor-image-upload">
      <image-upload
        ref="imageUploaderComponent"
        v-model="post.imageList"
        @update:model-value="onImageListChange"
      />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.simple-editor {
  border: 1px solid var(--border-color);
  border-radius: 3px;
  position: relative;
  width: 100%;

  .simple-editor-toolbar {
    width: 100%;
    height: 45px;
    display: flex;
    padding: 0 10px;
    align-items: center;
    background-color: var(--bg-color);
    top: 65px;
    z-index: 6;
    border-bottom: 1px solid var(--border-color);

    .act-btn {
      display: flex;
      padding: 0 10px;

      i {
        cursor: pointer;
        margin-left: 20px;
        font-size: 24px;

        &:first-child {
          margin-left: 0;
        }
      }
    }

    .publish-container {
      margin-left: auto;

      > .button-publish {
        margin-left: auto;

        ::v-deep span {
          font-size: 14px;
        }
      }

      > .tip {
        font-size: 14px;
        margin-right: 10px;
        color: var(--text-color4);
      }
    }
  }

  .simple-editor-input {
    width: 100%;

    textarea {
      font-family: inherit;
      background: var(--bg-color2);
      width: 100%;
      // min-height: 200px;
      border: 0;
      outline: 0;
      display: block;
      position: relative;
      resize: none;
      line-height: 1.8;
      padding: 15px 15px 20px;
      overflow: auto;
      overscroll-behavior: contain;
      transition: all 100ms linear;
      color: var(--text-color);
    }
  }

  .simple-editor-image-upload {
    background: var(--bg-color2);
    padding: 20px 20px 20px;
  }
}
</style>

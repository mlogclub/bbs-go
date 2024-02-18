<template>
  <!--
    @change="onChange"
    @progress="onProgress"
  -->
  <a-upload
    :action="action"
    :file-list="file ? [file] : []"
    :show-file-list="false"
    accept="images/*"
    name="image"
    with-credentials
    @change="onChange"
    @progress="onProgress"
    @success="onSuccess"
    @error="onError"
  >
    <template #upload-button>
      <div
        :class="`arco-upload-list-item${
          file && file.status === 'error' ? ' arco-upload-list-item-error' : ''
        }`"
      >
        <div
          v-if="modelValue"
          class="arco-upload-list-picture custom-upload-avatar"
        >
          <img :src="modelValue" />
          <div class="arco-upload-list-picture-mask">
            <IconEdit />
          </div>
          <a-progress
            v-if="file && file.status === 'uploading' && file.percent < 100"
            :percent="file.percent"
            type="circle"
            size="mini"
            :style="{
              position: 'absolute',
              left: '50%',
              top: '50%',
              transform: 'translateX(-50%) translateY(-50%)',
            }"
          />
        </div>
        <div v-else class="arco-upload-picture-card">
          <div class="arco-upload-picture-card-text">
            <IconPlus />
            <div style="margin-top: 10px; font-weight: 600">Upload</div>
          </div>
        </div>
      </div>
    </template>
  </a-upload>
</template>

<script setup>
  import { getToken } from '@/utils/auth';

  const emits = defineEmits(['update:modelValue']);

  defineProps({
    modelValue: {
      type: String,
      default: '',
    },
  });

  const file = ref();
  const action = computed(() => {
    const baseURL = import.meta.env.VITE_API_BASE_URL || '';
    return `${baseURL}/api/upload?userToken=${getToken()}`;
  });

  const onChange = (_, currentFile) => {
    file.value = {
      ...currentFile,
    };
  };
  const onProgress = (currentFile) => {
    file.value = currentFile;
  };
  const onError = () => {
    useNotificationError('上传失败');
  };
  const onSuccess = (ret) => {
    const resp = ret.response;
    if (resp.success) {
      emits('update:modelValue', resp.data.url);
      useNotificationSuccess('上传成功');
    } else {
      useHandleError(resp);
    }
  };
</script>

<style lang="scss" scoped></style>

<template>
  <a-form v-if="config" :model="config" auto-label-width>
    <!-- 上传方式选择 -->
    <a-card
      :title="$t('pages.settings.upload.uploadConfig')"
      class="upload-settings-card"
    >
      <a-form-item
        field="enableUploadMethod"
        :label="$t('pages.settings.upload.enableUploadMethod')"
      >
        <a-radio-group v-model="config.enableUploadMethod">
          <a-radio value="AliyunOss">{{
            $t('pages.settings.upload.aliyunOss')
          }}</a-radio>
          <a-radio value="TencentCos">{{
            $t('pages.settings.upload.tencentCos')
          }}</a-radio>
        </a-radio-group>
      </a-form-item>
    </a-card>

    <!-- 阿里云OSS配置 -->
    <a-card
      :title="$t('pages.settings.upload.aliyunOss')"
      class="upload-settings-card"
    >
      <a-form-item
        field="aliyunOss.host"
        :label="$t('pages.settings.upload.host')"
      >
        <a-input
          v-model="config.aliyunOss.host"
          :placeholder="$t('pages.settings.upload.placeholder.host')"
        />
      </a-form-item>

      <a-form-item
        field="aliyunOss.bucket"
        :label="$t('pages.settings.upload.bucket')"
      >
        <a-input
          v-model="config.aliyunOss.bucket"
          :placeholder="$t('pages.settings.upload.placeholder.bucket')"
        />
      </a-form-item>

      <a-form-item
        field="aliyunOss.endpoint"
        :label="$t('pages.settings.upload.endpoint')"
      >
        <a-input
          v-model="config.aliyunOss.endpoint"
          :placeholder="$t('pages.settings.upload.placeholder.endpoint')"
        />
      </a-form-item>

      <a-form-item
        field="aliyunOss.accessKeyId"
        :label="$t('pages.settings.upload.accessKeyId')"
      >
        <a-input
          v-model="config.aliyunOss.accessKeyId"
          :placeholder="$t('pages.settings.upload.placeholder.accessKeyId')"
        />
      </a-form-item>

      <a-form-item
        field="aliyunOss.accessKeySecret"
        :label="$t('pages.settings.upload.accessKeySecret')"
      >
        <a-input-password
          v-model="config.aliyunOss.accessKeySecret"
          :placeholder="$t('pages.settings.upload.placeholder.accessKeySecret')"
        />
      </a-form-item>

      <a-divider>{{ $t('pages.settings.upload.imageStyleConfig') }}</a-divider>

      <a-form-item
        field="aliyunOss.styleSplitter"
        :label="$t('pages.settings.upload.styleSplitter')"
      >
        <a-input
          v-model="config.aliyunOss.styleSplitter"
          :placeholder="$t('pages.settings.upload.placeholder.styleSplitter')"
        />
      </a-form-item>

      <a-form-item
        field="aliyunOss.styleAvatar"
        :label="$t('pages.settings.upload.styleAvatar')"
      >
        <a-input
          v-model="config.aliyunOss.styleAvatar"
          :placeholder="$t('pages.settings.upload.placeholder.styleAvatar')"
        />
      </a-form-item>

      <a-form-item
        field="aliyunOss.stylePreview"
        :label="$t('pages.settings.upload.stylePreview')"
      >
        <a-input
          v-model="config.aliyunOss.stylePreview"
          :placeholder="$t('pages.settings.upload.placeholder.stylePreview')"
        />
      </a-form-item>

      <a-form-item
        field="aliyunOss.styleSmall"
        :label="$t('pages.settings.upload.styleSmall')"
      >
        <a-input
          v-model="config.aliyunOss.styleSmall"
          :placeholder="$t('pages.settings.upload.placeholder.styleSmall')"
        />
      </a-form-item>

      <a-form-item
        field="aliyunOss.styleDetail"
        :label="$t('pages.settings.upload.styleDetail')"
      >
        <a-input
          v-model="config.aliyunOss.styleDetail"
          :placeholder="$t('pages.settings.upload.placeholder.styleDetail')"
        />
      </a-form-item>
    </a-card>

    <!-- 腾讯云COS配置 -->
    <a-card
      :title="$t('pages.settings.upload.tencentCos')"
      class="upload-settings-card"
    >
      <a-form-item
        field="tencentCos.bucket"
        :label="$t('pages.settings.upload.bucket')"
      >
        <a-input
          v-model="config.tencentCos.bucket"
          :placeholder="$t('pages.settings.upload.placeholder.tencentBucket')"
        />
      </a-form-item>

      <a-form-item
        field="tencentCos.region"
        :label="$t('pages.settings.upload.region')"
      >
        <a-input
          v-model="config.tencentCos.region"
          :placeholder="$t('pages.settings.upload.placeholder.region')"
        />
      </a-form-item>

      <a-form-item
        field="tencentCos.secretId"
        :label="$t('pages.settings.upload.secretId')"
      >
        <a-input
          v-model="config.tencentCos.secretId"
          :placeholder="$t('pages.settings.upload.placeholder.secretId')"
        />
      </a-form-item>

      <a-form-item
        field="tencentCos.secretKey"
        :label="$t('pages.settings.upload.secretKey')"
      >
        <a-input-password
          v-model="config.tencentCos.secretKey"
          :placeholder="$t('pages.settings.upload.placeholder.secretKey')"
        />
      </a-form-item>
    </a-card>

    <a-form-item>
      <a-button type="primary" :loading="loading" @click="submit">{{
        $t('pages.settings.upload.submit')
      }}</a-button>
    </a-form-item>
  </a-form>
</template>

<script setup>
  const { t } = useI18n();

  const loading = ref(false);
  const config = ref(null);

  const loadConfig = async () => {
    const ret = await axios.get('/api/admin/sys-config/configs');
    config.value = ret.uploadConfig || {
      enableUploadMethod: 'AliyunOss',
      aliyunOss: {
        host: '',
        bucket: '',
        endpoint: '',
        accessKeyId: '',
        accessKeySecret: '',
        styleSplitter: '@',
        styleAvatar: '',
        stylePreview: '',
        styleSmall: '',
        styleDetail: '',
      },
      tencentCos: {
        bucket: '',
        region: '',
        secretId: '',
        secretKey: '',
      },
    };
  };

  loadConfig();

  const submit = async () => {
    loading.value = true;
    try {
      await axios.post('/api/admin/sys-config/save', {
        uploadConfig: config.value,
      });
      await loadConfig();
      useNotificationSuccess(t('pages.settings.upload.message.submitSuccess'));
    } catch (e) {
      useHandleError(e);
    } finally {
      loading.value = false;
    }
  };
</script>

<style scoped lang="less">
  .upload-settings-card {
    margin-bottom: 20px;
  }
</style>

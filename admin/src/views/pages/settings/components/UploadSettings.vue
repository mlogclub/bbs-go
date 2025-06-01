<template>
  <a-form v-if="config" :model="config" auto-label-width>
    <!-- 上传方式选择 -->
    <a-card title="上传配置" class="upload-settings-card">
      <a-form-item field="enableUploadMethod" label="启用上传方式">
        <a-radio-group v-model="config.enableUploadMethod">
          <a-radio value="AliyunOss">阿里云OSS</a-radio>
          <a-radio value="TencentCos">腾讯云COS</a-radio>
        </a-radio-group>
      </a-form-item>
    </a-card>

    <!-- 阿里云OSS配置 -->
    <a-card title="阿里云OSS" class="upload-settings-card">
      <a-form-item field="aliyunOss.host" label="域名">
        <a-input
          v-model="config.aliyunOss.host"
          placeholder="请输入OSS域名，如：https://xxx.oss-cn-beijing.aliyuncs.com/"
        />
      </a-form-item>

      <a-form-item field="aliyunOss.bucket" label="Bucket">
        <a-input
          v-model="config.aliyunOss.bucket"
          placeholder="请输入Bucket名称"
        />
      </a-form-item>

      <a-form-item field="aliyunOss.endpoint" label="Endpoint">
        <a-input
          v-model="config.aliyunOss.endpoint"
          placeholder="请输入Endpoint，如：oss-cn-beijing.aliyuncs.com"
        />
      </a-form-item>

      <a-form-item field="aliyunOss.accessKeyId" label="AccessKey ID">
        <a-input
          v-model="config.aliyunOss.accessKeyId"
          placeholder="请输入AccessKey ID"
        />
      </a-form-item>

      <a-form-item field="aliyunOss.accessKeySecret" label="AccessKey Secret">
        <a-input-password
          v-model="config.aliyunOss.accessKeySecret"
          placeholder="请输入AccessKey Secret"
        />
      </a-form-item>

      <a-divider>图片样式配置</a-divider>

      <a-form-item field="aliyunOss.styleSplitter" label="样式分隔符">
        <a-input
          v-model="config.aliyunOss.styleSplitter"
          placeholder="样式分隔符，如：!"
        />
      </a-form-item>

      <a-form-item field="aliyunOss.styleAvatar" label="头像样式">
        <a-input
          v-model="config.aliyunOss.styleAvatar"
          placeholder="头像样式，如：100x100"
        />
      </a-form-item>

      <a-form-item field="aliyunOss.stylePreview" label="预览样式">
        <a-input
          v-model="config.aliyunOss.stylePreview"
          placeholder="预览样式，如：400x400"
        />
      </a-form-item>

      <a-form-item field="aliyunOss.styleSmall" label="小图样式">
        <a-input
          v-model="config.aliyunOss.styleSmall"
          placeholder="小图样式，如：200x200"
        />
      </a-form-item>

      <a-form-item field="aliyunOss.styleDetail" label="详情样式">
        <a-input
          v-model="config.aliyunOss.styleDetail"
          placeholder="详情样式，如：800x800"
        />
      </a-form-item>
    </a-card>

    <!-- 腾讯云COS配置 -->
    <a-card title="腾讯云COS" class="upload-settings-card">
      <a-form-item field="tencentCos.bucket" label="Bucket">
        <a-input
          v-model="config.tencentCos.bucket"
          placeholder="请输入Bucket名称，格式：bucket-appid"
        />
      </a-form-item>

      <a-form-item field="tencentCos.region" label="地域">
        <a-input
          v-model="config.tencentCos.region"
          placeholder="请输入地域，如：ap-beijing"
        />
      </a-form-item>

      <a-form-item field="tencentCos.secretId" label="SecretId">
        <a-input
          v-model="config.tencentCos.secretId"
          placeholder="请输入SecretId"
        />
      </a-form-item>

      <a-form-item field="tencentCos.secretKey" label="SecretKey">
        <a-input-password
          v-model="config.tencentCos.secretKey"
          placeholder="请输入SecretKey"
        />
      </a-form-item>
    </a-card>

    <a-form-item>
      <a-button type="primary" :loading="loading" @click="submit"
        >提交</a-button
      >
    </a-form-item>
  </a-form>
</template>

<script setup>
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
      useNotificationSuccess('提交成功');
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

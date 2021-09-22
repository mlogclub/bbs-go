<template>
  <el-upload
    v-loading="loading"
    class="uploader"
    :action="action"
    :with-credentials="true"
    :show-file-list="false"
    :headers="headers"
    name="image"
    accept="image/*"
    :before-upload="startLoad"
    :on-success="handleSuccess"
  >
    <img v-if="value" :src="value" class="upload-image" @load="finishLoad" @error="finishLoad" />
    <i class="el-icon-plus uploader-icon" :class="{ show: !value }" />
  </el-upload>
</template>

<script>
import { getToken } from "@/utils/auth";
export default {
  props: {
    value: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      loading: false,
    };
  },
  computed: {
    action() {
      return process.env.VUE_APP_BASE_API + "/api/upload";
    },
    headers() {
      const userToken = getToken();
      return {
        "X-User-Token": userToken || "",
      };
    },
  },
  methods: {
    handleSuccess(res, file) {
      if (!res.success) {
        this.$message.error(res.message || "上传失败");
        this.loading = false;
        return;
      }
      this.$emit("input", res.data.url);
      this.$message.success("上传成功");
    },
    startLoad() {
      this.loading = true;
    },
    finishLoad() {
      this.loading = false;
    },
  },
};
</script>

<style lang="scss" scoped>
.uploader {
  width: 178px;
  height: 178px;
  border: 2px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;

  &:hover {
    border-color: #409eff;

    .uploader-icon {
      display: block;
    }
  }

  .uploader-icon {
    font-size: 36px;
    font-weight: 700;
    color: #8c939d;
    width: 178px;
    height: 178px;
    line-height: 178px;
    text-align: center;
    position: absolute;
    top: 0;
    left: 0;
    display: none;

    &:hover {
      color: #409eff;
    }
    &.show {
      display: block;
    }
  }
  .upload-image {
    width: 178px;
    height: 178px;
    display: block;
  }
}
</style>

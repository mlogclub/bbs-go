<template>
  <div class="avatar-edit">
    <div class="avatar-view" :style="{ backgroundImage: 'url(' + value + ')' }">
      <div class="upload-view" @click="pickImage">
        <i class="iconfont icon-upload" />
        <span>{{ $t("component.avatarEdit.update") }}</span>
      </div>
    </div>

    <input
      ref="uploadImage"
      accept="image/*"
      type="file"
      @input="uploadAvatar"
    />
  </div>
</template>

<script setup>
const props = defineProps({
  value: {
    type: String,
    default: "",
  },
});

const { t } = useI18n();
const uploadImage = ref(null);

function pickImage() {
  const currentObj = uploadImage.value;
  currentObj.dispatchEvent(new MouseEvent("click"));
}

async function uploadAvatar(e) {
  const files = e.target.files;
  if (files.length <= 0) {
    return;
  }
  try {
    // 上传头像
    const file = files[0];
    const formData = new FormData();
    formData.append("image", file, file.name);
    const ret = await useHttpPost("/api/upload", formData);

    // 设置头像
    await useHttpPost(
      "/api/user/update/avatar",
      useJsonToForm({
        avatar: ret.url,
      })
    );

    // 重新加载数据
    const userStore = useUserStore();
    userStore.fetchCurrent();
    useMsgSuccess(t("component.avatarEdit.updateSuccess"));
  } catch (e) {
    console.error(e);
    useMsgError(t("component.avatarEdit.updateFailed"));
  }
}
</script>

<style lang="scss" scoped>
.avatar-edit {
  .avatar-view {
    width: 120px;
    height: 120px;
    background-size: cover;
    background-color: var(--bg-color2);
    border-radius: 50%;
    position: relative;

    &:hover {
      .upload-view {
        visibility: visible;
      }
    }

    .upload-view {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      color: var(--text-color);
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      border-radius: 50%;
      background-color: var(--bg-color-alpha);
      visibility: hidden;
      cursor: pointer;

      span {
        font-size: 13px;
        font-weight: 500;
      }
    }
  }

  input[type="file"] {
    display: none;
  }
}
</style>

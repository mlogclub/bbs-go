<template>
  <div
    class="profile"
    :style="{ backgroundImage: 'url(' + backgroundImage + ')' }"
  >
    <div v-if="isOwner" class="file is-light is-small change-bg">
      <label class="file-label">
        <input class="file-input" type="file" @change="uploadBackground" />
        <span class="file-cta">
          <span class="file-icon">
            <i class="iconfont icon-upload" />
          </span>
          <span class="file-label">设置背景</span>
        </span>
      </label>
    </div>
    <my-avatar
      :user="localUser"
      :round="true"
      :size="100"
      class="profile-avatar"
    />
    <div class="profile-info">
      <div class="metas">
        <h1 class="nickname">
          <nuxt-link :to="'/user/' + localUser.id">{{
            localUser.nickname
          }}</nuxt-link>
        </h1>
        <div v-if="localUser.description" class="description">
          <p>{{ localUser.description }}</p>
        </div>
      </div>
      <div class="action-btns">
        <follow-btn
          v-if="!currentUser || currentUser.id !== localUser.id"
          :user-id="localUser.id"
          :followed="followed"
          @onFollowed="onFollowed"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import defaultUserBg from "~/assets/images/default-user-bg.jpg";
const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
});
const localUser = ref(props.user);
const userStore = useUserStore();
const currentUser = computed(() => {
  return userStore.user;
});
const isOwner = computed(() => {
  return (
    localUser.value &&
    currentUser.value &&
    localUser.value.id === currentUser.value.id
  );
});
const backgroundImage = computed(() => {
  if (localUser.value.smallBackgroundImage) {
    return localUser.value.smallBackgroundImage;
  }
  return defaultUserBg;
});

const { data: followed } = await useAsyncData("followed", () =>
  useMyFetch(`/api/fans/isfollowed?userId=${localUser.value.id}`)
);

async function uploadBackground(e) {
  const files = e.target.files;
  if (files.length <= 0) {
    return;
  }
  try {
    // 上传头像
    const file = files[0];
    const formData = new FormData();
    formData.append("image", file, file.name);
    const ret = await useHttpPostMultipart("/api/upload", formData);

    // 设置背景
    await useHttpPostForm("/api/user/set_background_image", {
      body: {
        backgroundImage: ret.url,
      },
    });

    // 重新加载数据
    localUser.value = await userStore.fetchCurrent();

    useMsgSuccess("背景设置成功");
  } catch (e) {
    useMsgError(e.message || e);
    console.error(e);
  }
}

function onFollowed(userId, _followed) {
  followed.value = _followed;
}
</script>

<style lang="scss" scoped>
.profile {
  display: flex;
  margin-bottom: 10px;
  border-top-left-radius: 6px;
  border-top-right-radius: 6px;
  background-size: cover;
  background-position: 50%;
  height: 220px;

  // filter: blur(2px) contrast(0.8);

  position: relative;

  .profile-avatar {
    position: absolute;
    top: 90px;
    left: 10px;
  }

  .change-bg {
    position: absolute;
    top: 10px;
    right: 10px;
    opacity: 0.5;
    &:hover {
      opacity: 1;
    }
  }

  .profile-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    margin-top: 150px;
    padding: 10px 10px 10px 120px;
    background-image: linear-gradient(
      90deg,
      #ffffffff,
      rgba(255, 255, 255, 0.5),
      #dce9f200
    );

    .metas {
      display: flex;
      align-items: flex-start;
      flex-direction: column;
      width: 100%;

      .nickname {
        font-size: 18px;
        font-weight: 700;
        a {
          color: var(--text-color);
          &:hover {
            color: var(--text-color);
            text-decoration: underline;
          }
        }
      }

      .description {
        font-size: 14px;
        color: var(--text-color);
      }

      .homepage {
        font-size: 14px;
        a {
          color: var(--color2);
          &:hover {
            color: var(--text-link-color);
            text-decoration: underline;
          }
        }
      }
    }
    .action-btns {
      margin-left: 10px;
    }
  }
}

[data-theme="dark"],
.theme-dark,
.dark-mode {
  .profile-info {
    background-image: linear-gradient(
      90deg,
      #000000dd,
      rgba(255, 255, 255, 0.5),
      #dce9f200
    );
  }
}
</style>

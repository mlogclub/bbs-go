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
    <avatar
      :user="localUser"
      :round="true"
      :has-border="true"
      size="100"
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

<script>
export default {
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      localUser: Object.assign({}, this.user),
      followed: false,
    }
  },
  computed: {
    backgroundImage() {
      if (this.localUser.smallBackgroundImage) {
        return this.localUser.smallBackgroundImage
      }
      return require('~/assets/images/default-user-bg.jpg')
    },
    currentUser() {
      return this.$store.state.user.current
    },
    // 是否是主人态
    isOwner() {
      const current = this.$store.state.user.current
      return this.localUser && current && this.localUser.id === current.id
    },
  },
  mounted() {
    this.loadIsFollowed()
  },
  methods: {
    async uploadBackground(e) {
      const files = e.target.files
      if (files.length <= 0) {
        return
      }
      try {
        // 上传头像
        const file = files[0]
        const formData = new FormData()
        formData.append('image', file, file.name)
        const ret = await this.$axios.post('/api/upload', formData, {
          headers: { 'Content-Type': 'multipart/form-data' },
        })

        // 设置头像
        await this.$axios.post('/api/user/set/background/image', {
          backgroundImage: ret.url,
        })

        // 重新加载数据
        this.localUser = await this.$store.dispatch('user/getCurrentUser')

        this.$message.success('背景设置成功')
      } catch (e) {
        this.$message.error(e.message || e)
        console.error(e)
      }
    },
    async loadIsFollowed() {
      this.followed = await this.$axios.get(
        '/api/fans/isfollowed?userId=' + this.user.id
      )
    },
    onFollowed(userId, followed) {
      this.followed = followed
    },
  },
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
</style>

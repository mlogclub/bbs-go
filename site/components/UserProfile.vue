<template>
  <div
    class="profile"
    :class="{ background: backgroundImage }"
    :style="{ backgroundImage: 'url(' + localUser.smallBackgroundImage + ')' }"
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
    <div class="profile-info">
      <avatar
        v-if="backgroundImage"
        :user="localUser"
        :round="true"
        :has-border="true"
        size="100"
        :extra-style="{ position: 'absolute', top: '50px' }"
      />
      <avatar v-else :user="localUser" :round="true" size="100" />
      <div class="meta">
        <h1>
          <nuxt-link :to="'/user/' + localUser.id">{{
            localUser.nickname
          }}</nuxt-link>
        </h1>
        <div v-if="localUser.description" class="description">
          <p>{{ localUser.description }}</p>
        </div>
        <div v-if="localUser.homePage" class="homepage">
          <i class="iconfont icon-home"></i>
          <a
            :href="localUser.homePage"
            target="_blank"
            rel="external nofollow"
            >{{ localUser.homePage }}</a
          >
        </div>
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
    }
  },
  computed: {
    backgroundImage() {
      return this.localUser.backgroundImage
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
  },
}
</script>

<style lang="scss" scoped>
.profile {
  display: flex;
  margin-bottom: 10px;
  position: relative;
  border-top-left-radius: 6px;
  border-top-right-radius: 6px;

  .change-bg {
    position: absolute;
    top: 10px;
    right: 10px;
    opacity: 0.7;
    &:hover {
      opacity: 1;
    }
  }

  .profile-info {
    display: flex;
    width: 100%;
    padding: 10px 30px;
    background: #fff;

    // .avatar {
    //   max-width: 66px;
    //   max-height: 66px;
    //   min-width: 66px;
    //   min-height: 66px;
    // }

    .meta {
      margin-left: 18px;

      i {
        margin-right: 6px;
      }

      h1 {
        font-size: 28px;
        font-weight: 700;
        margin-bottom: 6px;
        a {
          color: #000;
          &:hover {
            color: #000;
            text-decoration: underline;
          }
        }
      }

      .description {
        font-size: 14px;
        color: 000;
        margin-bottom: 6px;
      }

      .homepage {
        font-size: 14px;
        a {
          color: #555;
          &:hover {
            color: #3273dc;
            text-decoration: underline;
          }
        }
      }
    }
  }

  &.background {
    //background-image: url('http://file.mlog.club/bg1.jpg!768_auto');
    background-size: cover;
    background-position: 50%;
    // filter: blur(2px) contrast(0.8);

    .profile-info {
      margin-top: 100px;
      background-color: unset;
      background-image: linear-gradient(
        90deg,
        // #dce9f25c
        #ffffffff,
        rgba(255, 255, 255, 0.5),
        #dce9f200
      );

      .meta {
        margin-left: 138px;
      }
    }
  }
}
</style>

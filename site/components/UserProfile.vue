<template>
  <div
    class="profile"
    :class="{ background: backgroundImage }"
    :style="{ backgroundImage: 'url(' + user.smallBackgroundImage + ')' }"
  >
    <div v-if="isOwner" class="file is-info is-small change-bg">
      <label class="file-label">
        <input class="file-input" type="file" @change="uploadBackground" />
        <span class="file-cta">
          <span class="file-icon">
            <i class="iconfont icon-upload" />
          </span>
          <span class="file-label">
            更换背景
          </span>
        </span>
      </label>
    </div>
    <div class="profile-info">
      <avatar
        v-if="backgroundImage"
        :user="user"
        :round="true"
        :has-border="true"
        size="100"
        :extra-style="{ position: 'absolute', top: '50px' }"
      />
      <avatar v-else :user="user" :round="true" size="100" />
      <div class="meta">
        <h1>
          <a :href="'/user/' + user.id">{{ user.nickname }}</a>
        </h1>
        <div v-if="user.description" class="description">
          <p>{{ user.description }}</p>
        </div>
        <div v-if="user.homePage" class="homepage">
          <i class="iconfont icon-home"></i>
          <a :href="user.homePage" target="_blank" rel="external nofollow">{{
            user.homePage
          }}</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Avatar from '~/components/Avatar'
export default {
  components: { Avatar },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
  computed: {
    backgroundImage() {
      return this.user.backgroundImage
    },
    currentUser() {
      return this.$store.state.user.current
    },
    // 是否是主人态
    isOwner() {
      const current = this.$store.state.user.current
      return this.user && current && this.user.id === current.id
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
        this.user = await this.$store.dispatch('user/getCurrentUser')

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
    padding: 10px;
    background: #fff;

    .avatar {
      max-width: 66px;
      max-height: 66px;
      min-width: 66px;
      min-height: 66px;
    }

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
        color: #555;
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

    .profile-info {
      margin-top: 100px;
      background-color: unset;
      background-image: linear-gradient(
        90deg,
        #dce9f25c,
        rgba(255, 255, 255, 0.76),
        #dce9f25c
      );

      .meta {
        margin-left: 138px;
      }
    }
  }
}
</style>

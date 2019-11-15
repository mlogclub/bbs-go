<template>
  <div>
    <div v-if="loading" class="loading modal is-active">
      <div class="modal-background" />
      <div class="modal-content">
        <div class="loading-animation" />
        <span class="loading-text">登录中，请稍后...</span>
      </div>
    </div>
  </div>
</template>

<script>
import utils from '~/common/utils'
export default {
  layout: 'no-footer',
  asyncData({ params, query }) {
    return {
      code: query.code,
      state: query.state,
      ref: query.ref
    }
  },
  data() {
    return {
      loading: false
    }
  },
  mounted() {
    this.callback()
  },
  methods: {
    async callback() {
      this.loading = true
      try {
        const user = await this.$store.dispatch('user/signinByGithub', {
          code: this.code,
          state: this.state
        })

        if (this.ref) {
          // 跳到登录前
          utils.linkTo(this.ref)
        } else {
          // 跳到个人主页
          utils.linkTo('/user/' + user.id)
        }
      } catch (e) {
        console.log(e)
        this.$toast.error('登录失败：' + (e.message || e), {
          onComplete() {
            utils.linkTo('/user/signin')
          }
        })
      } finally {
        this.loading = false
      }
    }
  },
  head() {
    return {
      title: this.$siteTitle('登录处理中...')
    }
  }
}
</script>

<style lang="scss" scoped>
.loading {
  .modal-background {
    background-color: rgba(10, 10, 10, 0.6);
  }
  .modal-content {
    text-align: center;
    color: #fdfdfd;
    font-weight: bold;
    font-size: 18px;
  }

  .loading-text {
    margin-left: 10px;
  }
}
</style>

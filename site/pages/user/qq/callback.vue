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
export default {
  layout: 'no-footer',
  asyncData({ params, query }) {
    return {
      code: query.code,
      state: query.state,
      ref: query.ref,
    }
  },
  data() {
    return {
      loading: false,
    }
  },
  mounted() {
    this.callback()
  },
  methods: {
    async callback() {
      this.loading = true
      try {
        const user = await this.$store.dispatch('user/signinByQQ', {
          code: this.code,
          state: this.state,
        })

        if (this.ref) {
          // 跳到登录前
          this.$linkTo(this.ref)
        } else {
          // 跳到个人主页
          this.$linkTo('/user/' + user.id)
        }
      } catch (e) {
        const me = this
        this.$msg({
          message: '登录失败' + (e.message || e),
          onClose() {
            me.$linkTo('/user/signin')
          },
        })
      } finally {
        this.loading = false
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle('登录处理中...'),
    }
  },
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

  .loading-animation {
    width: 20px;
    height: 20px;
    display: inline-block;
    color: red;
    vertical-align: middle;
    pointer-events: none;
    position: relative;
  }
  .loading-animation:before,
  .loading-animation:after {
    content: '';
    width: inherit;
    height: inherit;
    border-radius: 50%;
    background-color: currentcolor;
    opacity: 0.6;
    position: absolute;
    top: 0;
    left: 0;
    -webkit-animation: loading-animation 2s infinite ease-in-out;
    animation: loading-animation 2s infinite ease-in-out;
  }
  .loading-animation:after {
    -webkit-animation-delay: -1s;
    animation-delay: -1s;
  }
  @-webkit-keyframes loading-animation {
    0%,
    100% {
      -webkit-transform: scale(0);
      transform: scale(0);
    }
    50% {
      -webkit-transform: scale(1);
      transform: scale(1);
    }
  }
  @keyframes loading-animation {
    0%,
    100% {
      -webkit-transform: scale(0);
      transform: scale(0);
    }
    50% {
      -webkit-transform: scale(1);
      transform: scale(1);
    }
  }
}
</style>

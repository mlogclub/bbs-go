<template>
  <div>
    <div class="widget">
      <div class="widget-header">
        <span><i class="iconfont icon-setting" />&nbsp;编辑资料</span>
      </div>
      <div class="widget-content stable">
        <div class="str">
          <div class="slabel">用户名</div>
          <div class="svalue">{{ user.username }}</div>
        </div>
        <div class="str">
          <div class="slabel">昵称</div>
          <div class="svalue">{{ user.nickname }}</div>
        </div>
        <div class="str">
          <div class="slabel">个人简介</div>
          <div class="svalue">{{ user.description }}</div>
        </div>
        <div class="str">
          <div class="slabel">邮箱</div>
          <div class="svalue">{{ user.email }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  layout: 'ucenter',
  middleware: 'authenticated',
  async asyncData({ $axios }) {
    const user = await $axios.get('/api/user/current')
    return {
      user,
    }
  },
  data() {
    return {}
  },
  head() {
    return {
      title: this.$siteTitle(this.user.nickname + ' - 编辑资料'),
    }
  },
  methods: {},
}
</script>

<style lang="scss" scoped>
.control {
  a,
  label {
    // padding-top: .375em;
    // font-size: 14px;
    line-height: 32px;
  }
}

.file-cta {
  margin-top: 5px;
}
</style>

<template>
  <div>
    <div class="widget">
      <div class="widget-header">
        <span><i class="iconfont icon-setting" />&nbsp;编辑资料</span>
      </div>
      <div class="widget-content profile-edit-table">
        <div class="edit-row">
          <div class="edit-label">用户名</div>
          <div>
            <div class="svalue">{{ user.username }}</div>
          </div>
        </div>
        <div class="edit-row">
          <div class="edit-label">昵称</div>
          <div>
            <div class="svalue">{{ user.nickname }}</div>
          </div>
        </div>
        <div class="edit-row">
          <div class="edit-label">个人简介</div>
          <div>
            <div class="svalue">{{ user.description }}</div>
          </div>
        </div>
        <div class="edit-row">
          <div class="edit-label">邮箱</div>
          <div>
            <div class="svalue">{{ user.email }}</div>
          </div>
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
.profile-edit-table {
  margin: 30px;
  .edit-row {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .edit-label {
      text-align: right;
      min-width: 70px;
      font-size: 16px;
      font-weight: 500;
      color: #8590a6;
    }
  }
}
</style>

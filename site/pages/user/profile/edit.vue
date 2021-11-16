<template>
  <div>
    <div class="widget">
      <div class="widget-header">
        <span><i class="iconfont icon-setting" />&nbsp;编辑资料</span>
      </div>
      <div class="widget-content profile-edit-table">
        <div class="edit-row">
          <div class="edit-label">用户名</div>
          <div class="edit-content">
            <div>{{ user.username }}</div>
            <div><a>编辑</a></div>
          </div>
        </div>
        <div class="edit-row">
          <div class="edit-label">昵称</div>
          <div class="edit-content">
            <div>{{ user.nickname }}</div>
            <div><a>编辑</a></div>
          </div>
        </div>
        <div class="edit-row">
          <div class="edit-label">个人简介</div>
          <div class="edit-content">
            <div>{{ user.description }}</div>
            <div><a>编辑</a></div>
          </div>
        </div>
        <div class="edit-row">
          <div class="edit-label">邮箱</div>
          <div class="edit-content">
            <div>{{ user.email }}</div>
            <div><a>编辑</a></div>
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
  padding: 30px;
  .edit-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    line-height: 70px;

    .edit-label {
      min-width: 70px;
      text-align: right;
      font-size: 16px;
      margin-right: 42px;
      font-weight: 500;
      color: #8590a6;
    }

    &:not(:last-child) {
      .edit-content {
        border-bottom: 1px solid #eee;
      }
    }
    .edit-content {
      width: 100%;
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0 10px;
      color: #0b0b37;
      font-size: 16px;
      font-weight: 500;

      a {
        color: #2469f6;
        font-weight: 400;
        font-size: 14px;
      }
    }
  }
}
</style>

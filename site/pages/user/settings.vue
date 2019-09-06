<template>
  <section class="main">
    <div class="container">
      <div class="card">
        <header class="card-header">
          <p class="card-header-title">
            修改资料
          </p>
        </header>
        <div class="card-content">
          <div class="field">
            <label class="label"><span style="color:red;">*&nbsp;</span>用户名</label>
            <div class="control has-icons-left">
              <input
                v-model="user.username"
                class="input is-success"
                type="text"
                disabled="disabled"
              >
              <span class="icon is-small is-left"><i class="iconfont icon-username" /></span>
            </div>
          </div>

          <div class="field">
            <label class="label"><span style="color:red;">*&nbsp;</span>邮箱</label>
            <div class="control has-icons-left">
              <input
                v-model="user.email"
                class="input is-success"
                type="text"
                disabled="disabled"
              >
              <span class="icon is-small is-left"><i class="iconfont icon-email" /></span>
            </div>
          </div>

          <div class="field">
            <label class="label"><span style="color:red;">*&nbsp;</span>昵称</label>
            <div class="control has-icons-left">
              <input
                v-model="user.nickname"
                name="nickname"
                class="input is-success"
                type="text"
                placeholder="请输入昵称"
              >
              <span class="icon is-small is-left"><i class="iconfont icon-username" /></span>
            </div>
          </div>

          <div class="field">
            <label class="label"><span style="color:red;">*&nbsp;</span>头像</label>
            <div class="control">
              <img :src="user.avatar" style="width: 150px;height:150px;">
              <div class="file">
                <label class="file-label">
                  <input class="file-input" type="file" @change="uploadAvatar">
                  <span class="file-cta"><span class="file-icon"><i class="iconfont icon-upload" /></span>
                    <span class="file-label">选择头像</span></span>
                </label>
              </div>
              <span style="font-weight: bold; color:red;">*图像必须为正方形，大小不要超过1M。</span>
            </div>
          </div>

          <div class="field">
            <label class="label">简介</label>
            <div class="control">
              <textarea
                v-model="user.description"
                name="description"
                class="textarea"
                rows="2"
                placeholder="一句话介绍你自己"
              />
            </div>
          </div>

          <div class="field">
            <div class="control">
              <a class="button is-success" @click="submitForm">提交修改</a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  middleware: 'authenticated',
  head() {
    return {
      title: this.$siteTitle(this.user.nickname + ' - 编辑资料')
    }
  },
  async asyncData({ $axios, params }) {
    const [user] = await Promise.all([
      $axios.get('/api/user/current')
    ])
    return {
      user: user
    }
  },
  methods: {
    async submitForm() {
      try {
        await this.$axios.post('/api/user/edit/' + this.user.id, {
          nickname: this.user.nickname,
          avatar: this.user.avatar,
          description: this.user.description
        })
        this.$toast.success('修改成功')
      } catch (e) {
        console.error(e)
        this.$toast.error('修改失败：' + (e.message || e))
      }
    },
    async uploadAvatar(e) {
      const files = e.target.files
      if (files.length <= 0) {
        return
      }
      try {
        const file = files[0]
        const formData = new FormData()
        formData.append('image', file, file.name)
        const ret = await this.$axios.post('/api/upload', formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })
        this.user.avatar = ret.url
      } catch (e) {
        console.error(e)
      }
    }
  }
}
</script>

<style lang="scss" scoped>

</style>

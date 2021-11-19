<template>
  <div>
    <div class="widget">
      <div class="widget-header">
        <span>
          <i class="iconfont icon-setting" />
          <span>个人资料</span>
        </span>
      </div>
      <div class="widget-content">
        <!-- 头像 -->
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">头像：</label>
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control">
                <!-- <img
                        v-if="user.avatar"
                        :src="user.avatar"
                        style="width: 150px; height: 150px"
                      /> -->
                <avatar :user="user" :size="150" />
                <div class="file">
                  <label class="file-label">
                    <input
                      class="file-input"
                      type="file"
                      accept="image/png,image/jpeg,image/gif"
                      @change="uploadAvatar"
                    />
                    <span class="file-cta">
                      <span class="file-icon">
                        <i class="iconfont icon-upload" />
                      </span>
                      <span class="file-label">修改头像</span>
                    </span>
                  </label>
                </div>
                <span style="font-weight: bold; color: red"
                  >*图像必须为正方形，大小不要超过1M。</span
                >
              </div>
            </div>
          </div>
        </div>

        <!-- 昵称 -->
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">昵称：</label>
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.nickname"
                  class="input"
                  type="text"
                  autocomplete="off"
                  placeholder="请输入昵称"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-username" />
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- 简介 -->
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">简介：</label>
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control">
                <textarea
                  v-model="form.description"
                  class="textarea"
                  rows="2"
                  placeholder="一句话介绍你自己"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 个人主页 -->
        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">个人主页：</label>
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.homePage"
                  class="input"
                  type="text"
                  autocomplete="off"
                  placeholder="请输入个人主页"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-net" />
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="field is-horizontal">
          <div class="field-label is-normal" />
          <div class="field-body">
            <div class="field">
              <div class="control">
                <a class="button is-success" @click="submitForm">保存修改</a>
              </div>
            </div>
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
    const form = { ...user }
    return {
      user,
      form,
    }
  },
  data() {
    return {
      form: {
        username: '',
        email: '',
        nickname: '',
        avatar: '',
        homePage: '',
        description: '',
        password: '',
        rePassword: '',
        oldPassword: '',
      },

      showSetUsername: false,
      // username: '',

      showSetEmail: false,
      // email: '',

      showSetPassword: false, // 显示设置密码
      showUpdatePassword: false, // 显示修改密码
      // password: '',
      // rePassword: '',
      // oldPassword: ''
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.user.nickname + ' - 编辑资料'),
    }
  },
  methods: {
    async submitForm() {
      try {
        await this.$axios.post('/api/user/edit/' + this.user.id, {
          nickname: this.form.nickname,
          avatar: this.form.avatar,
          homePage: this.form.homePage,
          description: this.form.description,
        })
        await this.reload()
        this.$message.success('资料修改成功')
      } catch (e) {
        console.error(e)
        this.$message.error('资料修改失败：' + (e.message || e))
      }
    },
    async uploadAvatar(e) {
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
        await this.$axios.post('/api/user/update/avatar', {
          avatar: ret.url,
        })

        // 重新加载数据
        await this.reload()

        this.$message.success('头像更新成功')
      } catch (e) {
        console.error(e)
      }
    },
    async setUsername() {
      try {
        const me = this
        await this.$axios.post('/api/user/set/username', {
          username: me.form.username,
        })
        await this.reload()
        this.$message.success('用户名设置成功')
        this.showSetUsername = false
      } catch (err) {
        this.$message.error('用户名设置失败：' + (err.message || err))
      }
    },
    async setEmail() {
      try {
        const me = this
        await this.$axios.post('/api/user/set/email', {
          email: me.form.email,
        })
        await this.reload()
        this.$message.success('邮箱设置成功')
        this.showSetEmail = false
      } catch (err) {
        this.$message.error('邮箱设置失败：' + (err.message || err))
      }
    },
    async setPassword() {
      try {
        const me = this
        await this.$axios.post('/api/user/set/password', {
          password: me.form.password,
          rePassword: me.form.rePassword,
        })
        await this.reload()
        this.$message.success('密码设置成功')
        this.showSetPassword = false
      } catch (err) {
        this.$message.error('密码设置失败：' + (err.message || err))
      }
    },
    async updatePassword() {
      try {
        const me = this
        await this.$axios.post('/api/user/update/password', {
          oldPassword: me.form.oldPassword,
          password: me.form.password,
          rePassword: me.form.rePassword,
        })
        await this.reload()
        this.$message.success('密码修改成功')
        this.showUpdatePassword = false
      } catch (err) {
        this.$message.error('密码修改失败：' + (err.message || err))
      }
    },
    async reload() {
      this.user = await this.$axios.get('/api/user/current')
      this.form = { ...this.user }
    },
    async requestEmailVerify() {
      this.$nuxt.$loading.start()
      try {
        await this.$axios.post('/api/user/email/verify')
        this.$message.success(
          '邮件已经发送到你的邮箱：' + this.user.email + '，请注意查收。'
        )
      } catch (err) {
        this.$message.error('请求验证失败：' + (err.message || err))
      } finally {
        this.$nuxt.$loading.finish()
      }
    },
  },
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

.modal {
  .widget {
    background: #ffffff;
    margin: 0px;
    padding: 10px;
  }
}

.file-cta {
  margin-top: 5px;
}
</style>

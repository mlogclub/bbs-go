<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <div class="widget">
          <div class="widget-header">
            <nav class="breadcrumb">
              <ul>
                <li>
                  <a href="/">首页</a>
                </li>
                <li>
                  <a :href="'/user/' + user.id">{{ user.nickname }}</a>
                </li>
                <li class="is-active">
                  <a href="#" aria-current="page">编辑资料</a>
                </li>
              </ul>
            </nav>
          </div>
          <div class="widget-content">
            <div class="field is-horizontal">
              <div class="field-label is-normal">
                <label class="label">用户名：</label>
              </div>
              <div class="field-body">
                <div class="field">
                  <div class="control has-icons-left">
                    <label v-if="user.username">{{ user.username }}</label>
                    <a v-else @click="showSetUsername = true">点击设置</a>
                  </div>
                </div>
              </div>
            </div>

            <!-- 邮箱 -->
            <div class="field is-horizontal">
              <div class="field-label is-normal">
                <label class="label">邮箱：</label>
              </div>
              <div class="field-body">
                <div class="field">
                  <div class="control has-icons-left">
                    <label v-if="user.email">{{ user.email }}</label>
                    <a v-else @click="showSetEmail = true">点击设置</a>
                  </div>
                </div>
              </div>
            </div>

            <!-- 密码 -->
            <div class="field is-horizontal">
              <div class="field-label is-normal">
                <label class="label">密码：</label>
              </div>
              <div class="field-body">
                <div class="field">
                  <div class="control has-icons-left">
                    <template v-if="user.passwordSet">
                      <label>密码已设置&nbsp;</label>
                      <a @click="showUpdatePassword = true">点击修改</a>
                    </template>
                    <a v-else @click="showSetPassword = true">点击设置</a>
                  </div>
                </div>
              </div>
            </div>

            <!-- 头像 -->
            <div class="field is-horizontal">
              <div class="field-label is-normal">
                <label class="label">
                  <span style="color:red;">*&nbsp;</span>头像：
                </label>
              </div>
              <div class="field-body">
                <div class="field">
                  <div class="control">
                    <img
                      :src="user.avatar"
                      style="width: 150px;height:150px;"
                    />
                    <div class="file">
                      <label class="file-label">
                        <input
                          @change="uploadAvatar"
                          class="file-input"
                          type="file"
                        />
                        <span class="file-cta">
                          <span class="file-icon">
                            <i class="iconfont icon-upload" />
                          </span>
                          <span class="file-label">选择头像</span>
                        </span>
                      </label>
                    </div>
                    <span style="font-weight: bold; color:red;"
                      >*图像必须为正方形，大小不要超过1M。</span
                    >
                  </div>
                </div>
              </div>
            </div>

            <!-- 昵称 -->
            <div class="field is-horizontal">
              <div class="field-label is-normal">
                <label class="label">
                  <span style="color:red;">*&nbsp;</span>昵称：
                </label>
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
                    <a @click="submitForm" class="button is-success"
                      >提交修改</a
                    >
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <user-center-sidebar :user="user" />
    </div>

    <!-- 设置用户名 -->
    <div :class="{ 'is-active': showSetUsername }" class="modal">
      <div class="modal-background" />
      <div class="modal-card">
        <div class="widget">
          <div class="widget-header">
            设置用户名
            <button
              @click="showSetUsername = false"
              class="delete"
              aria-label="close"
            />
          </div>
          <div class="widget-content">
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.username"
                  @keydown.enter="setUsername"
                  class="input is-success"
                  type="text"
                  placeholder="请输入用户名"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-username" />
                </span>
              </div>
            </div>
          </div>
          <div class="widget-footer is-right">
            <a @click="setUsername" class="button is-success">确定</a>
            <a @click="showSetUsername = false" class="button">取消</a>
          </div>
        </div>
      </div>
    </div>

    <!-- 设置邮箱 -->
    <div :class="{ 'is-active': showSetEmail }" class="modal">
      <div class="modal-background" />
      <div class="modal-card">
        <div class="widget">
          <div class="widget-header">
            设置邮箱
            <button
              @click="showSetEmail = false"
              class="delete"
              aria-label="close"
            />
          </div>
          <div class="widget-content">
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.email"
                  @keydown.enter="setEmail"
                  class="input is-success"
                  type="text"
                  placeholder="请输入邮箱"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-email" />
                </span>
              </div>
            </div>
          </div>
          <div class="widget-footer is-right">
            <a @click="setEmail" class="button is-success">确定</a>
            <a @click="showSetEmail = false" class="button">取消</a>
          </div>
        </div>
      </div>
    </div>

    <!-- 设置密码 -->
    <div :class="{ 'is-active': showSetPassword }" class="modal">
      <div class="modal-background" />
      <div class="modal-card">
        <div class="widget">
          <div class="widget-header">
            设置密码
            <button
              @click="showSetPassword = false"
              class="delete"
              aria-label="close"
            />
          </div>
          <div class="widget-content">
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.password"
                  @keydown.enter="setPassword"
                  class="input is-success"
                  type="password"
                  placeholder="请输入密码"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.rePassword"
                  @keydown.enter="setPassword"
                  class="input is-success"
                  type="password"
                  placeholder="请再次确认密码"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>
          </div>
          <div class="widget-footer is-right">
            <a @click="setPassword" class="button is-success">确定</a>
            <a @click="showSetPassword = false" class="button">取消</a>
          </div>
        </div>
      </div>
    </div>

    <!-- 修改密码 -->
    <div :class="{ 'is-active': showUpdatePassword }" class="modal">
      <div class="modal-background" />
      <div class="modal-card">
        <div class="widget">
          <div class="widget-header">
            修改密码
            <button
              @click="showUpdatePassword = false"
              class="delete"
              aria-label="close"
            />
          </div>
          <div class="widget-content">
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.oldPassword"
                  @keydown.enter="updatePassword"
                  class="input is-success"
                  type="password"
                  placeholder="请输入当前密码"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.password"
                  @keydown.enter="updatePassword"
                  class="input is-success"
                  type="password"
                  placeholder="请输入密码"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.rePassword"
                  @keydown.enter="updatePassword"
                  class="input is-success"
                  type="password"
                  placeholder="请再次确认密码"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>
          </div>
          <div class="widget-footer is-right">
            <a @click="updatePassword" class="button is-success">确定</a>
            <a @click="showUpdatePassword = false" class="button">取消</a>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import UserCenterSidebar from '~/components/UserCenterSidebar'
export default {
  middleware: 'authenticated',
  components: {
    UserCenterSidebar
  },
  async asyncData({ $axios, params }) {
    const user = await $axios.get('/api/user/current')
    const form = { ...user }
    return {
      user,
      form
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
        oldPassword: ''
      },

      showSetUsername: false,
      // username: '',

      showSetEmail: false,
      // email: '',

      showSetPassword: false, // 显示设置密码
      showUpdatePassword: false // 显示修改密码
      // password: '',
      // rePassword: '',
      // oldPassword: ''
    }
  },
  methods: {
    async submitForm() {
      try {
        await this.$axios.post('/api/user/edit/' + this.user.id, {
          nickname: this.form.nickname,
          avatar: this.form.avatar,
          homePage: this.form.homePage,
          description: this.form.description
        })
        await this.reload()
        this.$toast.success('资料修改成功')
      } catch (e) {
        console.error(e)
        this.$toast.error('资料修改失败：' + (e.message || e))
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
          headers: { 'Content-Type': 'multipart/form-data' }
        })

        // 设置头像
        await this.$axios.post('/api/user/update/avatar', {
          avatar: ret.url
        })

        // 重新加载数据
        await this.reload()

        this.$toast.success('头像更新成功')
      } catch (e) {
        console.error(e)
      }
    },
    async setUsername() {
      try {
        const me = this
        await this.$axios.post('/api/user/set/username', {
          username: me.form.username
        })
        await this.reload()
        this.$toast.success('用户名设置成功')
        this.showSetUsername = false
      } catch (err) {
        this.$toast.error('用户名设置失败：' + (err.message || err))
      }
    },
    async setEmail() {
      try {
        const me = this
        await this.$axios.post('/api/user/set/email', {
          email: me.form.email
        })
        await this.reload()
        this.$toast.success('邮箱设置成功')
        this.showSetEmail = false
      } catch (err) {
        this.$toast.error('邮箱设置失败：' + (err.message || err))
      }
    },
    async setPassword() {
      try {
        const me = this
        await this.$axios.post('/api/user/set/password', {
          password: me.form.password,
          rePassword: me.form.rePassword
        })
        await this.reload()
        this.$toast.success('密码设置成功')
        this.showSetPassword = false
      } catch (err) {
        this.$toast.error('密码设置失败：' + (err.message || err))
      }
    },
    async updatePassword() {
      try {
        const me = this
        await this.$axios.post('/api/user/update/password', {
          oldPassword: me.form.oldPassword,
          password: me.form.password,
          rePassword: me.form.rePassword
        })
        await this.reload()
        this.$toast.success('密码修改成功')
        this.showUpdatePassword = false
      } catch (err) {
        this.$toast.error('密码修改失败：' + (err.message || err))
      }
    },
    async reload() {
      this.user = await this.$axios.get('/api/user/current')
      this.form = { ...this.user }
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.user.nickname + ' - 编辑资料')
    }
  }
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
</style>

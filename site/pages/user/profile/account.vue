<template>
  <div>
    <div class="widget no-margin">
      <div class="widget-header">
        <div>
          <i class="iconfont icon-setting" />
          <span>编辑资料</span>
        </div>
        <nuxt-link :to="'/user/' + user.id" style="font-size: 13px">
          <i class="iconfont icon-return" />
          <span>返回个人主页</span>
        </nuxt-link>
      </div>
      <div class="widget-content">
        <!-- <div class="my-field">
              <div>用户名</div>
              <div>
                <span v-if="user.username">{{ user.username }}</span>
              </div>
              <div>
                <a @click="showSetUsername = true">点击设置</a>
              </div>
            </div> -->
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

        <div class="field is-horizontal">
          <div class="field-label is-normal">
            <label class="label">邮箱：</label>
          </div>
          <div class="field-body">
            <div class="field">
              <div class="control has-icons-left">
                <template v-if="user.email">
                  <label>{{ user.email }}</label>
                  <a @click="showSetEmail = true">修改</a>
                  <a
                    v-if="!user.emailVerified"
                    class="has-text-danger"
                    style="font-weight: 700"
                    @click="requestEmailVerify"
                    >验证&gt;&gt;</a
                  >
                </template>
                <template v-else>
                  <a @click="showSetEmail = true">点击设置</a>
                </template>
              </div>
            </div>
          </div>
        </div>

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
      </div>
    </div>

    <!-- 设置用户名 -->
    <div :class="{ 'is-active': showSetUsername }" class="modal">
      <div class="modal-background" />
      <div class="modal-card">
        <div class="widget">
          <div class="widget-header">
            设置用户名
            <button
              class="delete"
              aria-label="close"
              @click="showSetUsername = false"
            />
          </div>
          <div class="widget-content">
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.username"
                  class="input is-success"
                  type="text"
                  placeholder="请输入用户名"
                  @keydown.enter="setUsername"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-username" />
                </span>
              </div>
            </div>
          </div>
          <div class="widget-footer is-right">
            <a class="button is-success" @click="setUsername">确定</a>
            <a class="button" @click="showSetUsername = false">取消</a>
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
              class="delete"
              aria-label="close"
              @click="showSetEmail = false"
            />
          </div>
          <div class="widget-content">
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.email"
                  class="input is-success"
                  type="text"
                  placeholder="请输入邮箱"
                  @keydown.enter="setEmail"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-email" />
                </span>
              </div>
            </div>
          </div>
          <div class="widget-footer is-right">
            <a class="button is-success" @click="setEmail">确定</a>
            <a class="button" @click="showSetEmail = false">取消</a>
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
              class="delete"
              aria-label="close"
              @click="showSetPassword = false"
            />
          </div>
          <div class="widget-content">
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.password"
                  class="input is-success"
                  type="password"
                  placeholder="请输入密码"
                  @keydown.enter="setPassword"
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
                  class="input is-success"
                  type="password"
                  placeholder="请再次确认密码"
                  @keydown.enter="setPassword"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>
          </div>
          <div class="widget-footer is-right">
            <a class="button is-success" @click="setPassword">确定</a>
            <a class="button" @click="showSetPassword = false">取消</a>
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
              class="delete"
              aria-label="close"
              @click="showUpdatePassword = false"
            />
          </div>
          <div class="widget-content">
            <div class="field">
              <div class="control has-icons-left">
                <input
                  v-model="form.oldPassword"
                  class="input is-success"
                  type="password"
                  placeholder="请输入当前密码"
                  @keydown.enter="updatePassword"
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
                  class="input is-success"
                  type="password"
                  placeholder="请输入密码"
                  @keydown.enter="updatePassword"
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
                  class="input is-success"
                  type="password"
                  placeholder="请再次确认密码"
                  @keydown.enter="updatePassword"
                />
                <span class="icon is-small is-left">
                  <i class="iconfont icon-password" />
                </span>
              </div>
            </div>
          </div>
          <div class="widget-footer is-right">
            <a class="button is-success" @click="updatePassword">确定</a>
            <a class="button" @click="showUpdatePassword = false">取消</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
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
        password: '',
        rePassword: '',
        oldPassword: '',
      },
      showSetUsername: false,
      showSetEmail: false,
      showSetPassword: false, // 显示设置密码
      showUpdatePassword: false, // 显示修改密码
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.user.nickname + ' - 账号设置'),
    }
  },
  methods: {
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
    async requestEmailVerify() {
      this.$nuxt.$loading.start()
      try {
        await this.$axios.post('/api/user/send_verify_email')
        this.$message.success(
          '邮件已经发送到你的邮箱：' + this.user.email + '，请注意查收。'
        )
      } catch (err) {
        this.$message.error('请求验证失败：' + (err.message || err))
      } finally {
        this.$nuxt.$loading.finish()
      }
    },
    async reload() {
      this.user = await this.$axios.get('/api/user/current')
      this.form = { ...this.user }
    },
  },
}
</script>

<style lang="scss" scoped>
.my-field {
  display: flex;
}
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
    background-color: var(--bg-color);
    margin: 0;
    padding: 10px;
  }
}

.file-cta {
  margin-top: 5px;
}
</style>

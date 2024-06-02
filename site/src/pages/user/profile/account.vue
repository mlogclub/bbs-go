<template>
  <div v-if="user" class="widget no-margin">
    <div class="widget-header">
      <div class="account">
        <i class="iconfont icon-setting" />
        <span>账号设置</span>
      </div>
      <nuxt-link :to="'/user/' + user.id">
        <i class="iconfont icon-return" />
        <span>返回个人主页</span>
      </nuxt-link>
    </div>
    <div class="widget-content">
      <div class="settings">
        <div class="settings-item">
          <div class="settings-item-title">用户名</div>
          <div class="settings-item-input">
            <div class="input-value">{{ user.username }}</div>
            <div class="action-box">
              <a v-if="!user.username" @click="showUsernameDialog">设置</a>
            </div>
          </div>
        </div>

        <div class="settings-item">
          <div class="settings-item-title">邮箱</div>
          <div class="settings-item-input">
            <div class="input-value">
              <span>{{ user.email }}</span>
              <span
                v-if="user.emailVerified"
                style="margin-left: 4px; font-size: 80%"
                >(已验证)</span
              >
            </div>
            <div class="action-box">
              <a v-if="user.email" @click="showEmailDialog">修改</a>
              <a
                v-if="user.email && !user.emailVerified"
                @click="requestEmailVerify"
                >验证</a
              >
              <a v-if="!user.email">设置</a>
            </div>
          </div>
        </div>

        <div class="settings-item">
          <div class="settings-item-title">密码</div>
          <div class="settings-item-input">
            <div class="input-value">
              {{ user.passwordSet ? "已设置" : "未设置" }}
            </div>
            <div class="action-box">
              <a v-if="user.passwordSet" @click="showUpdatePasswordDialog"
                >修改</a
              >
              <a v-else @click="showSetPasswordDialog">设置</a>
            </div>
          </div>
        </div>
      </div>
    </div>

    <my-dialog
      ref="usernameDialog"
      v-model:visible="usernameDialogVisible"
      title="设置用户名"
      :width="320"
      @ok="setUsername"
    >
      <div style="padding: 30px 0">
        <input
          v-model="usernameForm.username"
          class="input is-small"
          type="text"
          placeholder="用户名"
        />
      </div>
    </my-dialog>

    <my-dialog
      ref="emailDialog"
      v-model:visible="emailDialogVisible"
      title="设置邮箱"
      :width="320"
      @ok="setEmail"
    >
      <div style="padding: 30px 0">
        <input
          v-model="emailForm.email"
          class="input is-small"
          type="text"
          placeholder="用户名"
        />
      </div>
    </my-dialog>

    <my-dialog
      ref="updatePasswordDialog"
      v-model:visible="updatePasswordDialogVisible"
      title="修改密码"
      :width="320"
      @ok="updatePassword"
    >
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="updatePasswordForm.oldPassword"
            class="input is-small"
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
            v-model="updatePasswordForm.password"
            class="input is-small"
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
            v-model="updatePasswordForm.rePassword"
            class="input is-small"
            type="password"
            placeholder="请再次确认密码"
            @keydown.enter="updatePassword"
          />
          <span class="icon is-small is-left">
            <i class="iconfont icon-password" />
          </span>
        </div>
      </div>
    </my-dialog>

    <my-dialog
      ref="setPasswordDialog"
      v-model:visible="setPasswordDialogVisible"
      title="设置密码"
      :width="320"
      @ok="setPassword"
    >
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="setPasswordForm.password"
            class="input is-small"
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
            v-model="setPasswordForm.rePassword"
            class="input is-small"
            type="password"
            placeholder="请再次确认密码"
            @keydown.enter="setPassword"
          />
          <span class="icon is-small is-left">
            <i class="iconfont icon-password" />
          </span>
        </div>
      </div>
    </my-dialog>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: ["auth"],
  layout: "profile",
});

useHead({
  title: useSiteTitle("账号设置"),
});

const { data: user, refresh: userRefresh } = await useAsyncData("user", () =>
  useMyFetch("/api/user/current")
);

const usernameDialog = ref(null);
const usernameDialogVisible = ref(false);
const usernameForm = reactive({
  username: user.value ? user.value.username : "",
});
function showUsernameDialog() {
  usernameDialog.value.show();
}
async function setUsername() {
  try {
    await useHttpPostForm("/api/user/set/username", {
      body: {
        username: usernameForm.username,
      },
    });
    await userRefresh();
    useMsgSuccess("用户名设置成功");
    usernameDialog.value.close();
  } catch (err) {
    useMsgError("用户名设置失败：" + (err.message || err));
  }
}

const emailDialog = ref(null);
const emailDialogVisible = ref(false);
const emailForm = reactive({
  email: user.value ? user.value.email : "",
});
function showEmailDialog() {
  emailDialog.value.show();
}
async function setEmail() {
  try {
    await useHttpPostForm("/api/user/set/email", {
      body: {
        email: emailForm.email,
      },
    });
    await userRefresh();
    useMsgSuccess("邮箱设置成功");
    emailDialog.value.close();
  } catch (err) {
    useMsgError("邮箱设置失败：" + (err.message || err));
  }
}

async function requestEmailVerify() {
  const loading = useLoading();
  try {
    await useHttpPost("/api/user/send_verify_email");
    useMsgSuccess(
      "邮件已经发送到你的邮箱：" + user.value.email + "，请注意查收。"
    );
  } catch (err) {
    useMsgError("请求验证失败：" + (err.message || err));
  } finally {
    loading.close();
  }
}

const updatePasswordDialog = ref(null);
const updatePasswordDialogVisible = ref(false);
const updatePasswordForm = reactive({
  password: "",
  rePassword: "",
});
function showUpdatePasswordDialog() {
  updatePasswordDialog.value.show();
}
async function updatePassword() {
  try {
    await useHttpPostForm("/api/user/update/password", {
      body: updatePasswordForm,
    });
    await userRefresh();
    useMsgSuccess("密码修改成功");
    updatePasswordDialog.value.close();
  } catch (err) {
    useMsgError("密码修改失败：" + (err.message || err));
  }
}

const setPasswordDialog = ref(null);
const setPasswordDialogVisible = ref(false);
const setPasswordForm = reactive({
  password: "",
  rePassword: "",
});
function showSetPasswordDialog() {
  setPasswordDialog.value.show();
}
async function setPassword() {
  try {
    await useHttpPostForm("/api/user/set/password", {
      body: setPasswordForm,
    });
    await userRefresh();
    useMsgSuccess("密码修改成功");
    setPasswordDialog.value.close();
  } catch (err) {
    useMsgError("密码修改失败：" + (err.message || err));
  }
}
</script>
<style lang="scss" scoped>
.field {
  margin-bottom: 10px;

  input {
    &:focus-visible {
      outline-width: 0;
    }
  }
}
.widget-header {
  padding: 18px 0;

  & > div {
    display: flex;
    align-items: center;
    span {
      margin-left: 6px;
    }
  }

  & > a {
    display: flex;
    align-items: center;
    font-weight: 500;
    font-size: 12px;
    span {
      margin-left: 6px;
    }
  }
}
.settings {
  .settings-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 18px 0;

    @media (max-width: 600px) {
      flex-direction: column;
      align-items: flex-start;

      .settings-item-title {
        margin-bottom: 6px;
      }
    }

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }
    .settings-item-title {
      width: 100px;
      color: var(--text-color2);
      font-size: 15px;
    }

    .settings-item-input {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
      font-size: 14px;
      .input-value {
        flex: 1;
        color: var(--text-color3);
      }
      .action-box {
        display: flex;
        align-items: center;
        column-gap: 10px;
        a {
          color: var(--text-link-color);
          font-size: 12px;

          &:hover {
            color: var(--text-link-hover-color);
          }
        }
      }
    }
  }
}
</style>

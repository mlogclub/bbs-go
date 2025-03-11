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
              <a v-if="!user.email" @click="showEmailDialog">设置</a>
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

    <AccountSetUsernameDialog ref="setUsernameDialog" @success="userRefresh" />
    <AccountSetEmailDialog ref="setEmailDialog" @success="userRefresh" />
    <AccountSetPasswordDialog ref="setPasswordDialog" @success="userRefresh" />
    <AccountUpdatePasswordDialog
      ref="updatePasswordDialog"
      @success="userRefresh"
    />
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
  useHttpGet("/api/user/current")
);

const setUsernameDialog = ref(null);
const setEmailDialog = ref(null);
const setPasswordDialog = ref(null);
const updatePasswordDialog = ref(null);
const showUsernameDialog = () => setUsernameDialog.value.show();
const showEmailDialog = () => setEmailDialog.value.show();
const showSetPasswordDialog = () => setPasswordDialog.value.show();
const showUpdatePasswordDialog = () => updatePasswordDialog.value.show();

async function requestEmailVerify() {
  const loading = useLoading();
  try {
    await useHttpPost("/api/user/send_verify_email");
    useMsgSuccess(
      "邮件已经发送到你的邮箱：" + user.value.email + "，请注意查收。"
    );
  } catch (err) {
    useMsgError(err.message || err);
  } finally {
    loading.close();
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
      border-bottom: 1px solid var(--border-color4);
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

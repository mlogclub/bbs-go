<template>
  <div v-if="user" class="widget no-margin">
    <div class="widget-header">
      <div class="account">
        <i class="iconfont icon-setting" />
        <span>{{ $t("user.profile.account.title") }}</span>
      </div>
      <nuxt-link :to="'/user/' + user.id">
        <i class="iconfont icon-return" />
        <span>{{ $t("user.profile.backToProfile") }}</span>
      </nuxt-link>
    </div>
    <div class="widget-content">
      <div class="settings">
        <div class="settings-item">
          <div class="settings-item-title">
            {{ $t("user.profile.account.username") }}
          </div>
          <div class="settings-item-input">
            <div class="input-value">{{ user.username }}</div>
            <div class="action-box">
              <a v-if="!user.username" @click="showUsernameDialog">{{
                $t("user.profile.account.set")
              }}</a>
            </div>
          </div>
        </div>

        <div class="settings-item">
          <div class="settings-item-title">
            {{ $t("user.profile.account.email") }}
          </div>
          <div class="settings-item-input">
            <div class="input-value">
              <span>{{ user.email }}</span>
              <span
                v-if="user.emailVerified"
                style="margin-left: 4px; font-size: 80%"
                >({{ $t("user.profile.account.verified") }})</span
              >
            </div>
            <div class="action-box">
              <a v-if="user.email" @click="showEmailDialog">{{
                $t("user.profile.account.modify")
              }}</a>
              <a
                v-if="user.email && !user.emailVerified"
                @click="requestEmailVerify"
                >{{ $t("user.profile.account.verify") }}</a
              >
              <a v-if="!user.email" @click="showEmailDialog">{{
                $t("user.profile.account.set")
              }}</a>
            </div>
          </div>
        </div>

        <div class="settings-item">
          <div class="settings-item-title">
            {{ $t("user.profile.account.password") }}
          </div>
          <div class="settings-item-input">
            <div class="input-value">
              {{
                user.passwordSet
                  ? $t("user.profile.account.passwordSet")
                  : $t("user.profile.account.passwordNotSet")
              }}
            </div>
            <div class="action-box">
              <a v-if="user.passwordSet" @click="showUpdatePasswordDialog">{{
                $t("user.profile.account.modify")
              }}</a>
              <a v-else @click="showSetPasswordDialog">{{
                $t("user.profile.account.set")
              }}</a>
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
    <AccountWxBindDialog ref="wxBindDialog" />
  </div>
</template>

<script setup>
definePageMeta({
  middleware: ["auth"],
  layout: "profile",
});

const { t } = useI18n();

useHead({
  title: useSiteTitle(t("user.profile.account.title")),
});

const { data: user, refresh: userRefresh } = await useMyFetch(
  "/api/user/current"
);

const setUsernameDialog = ref(null);
const setEmailDialog = ref(null);
const setPasswordDialog = ref(null);
const updatePasswordDialog = ref(null);
const wxBindDialog = ref(null);
const showUsernameDialog = () => setUsernameDialog.value.show();
const showEmailDialog = () => setEmailDialog.value.show();
const showSetPasswordDialog = () => setPasswordDialog.value.show();
const showUpdatePasswordDialog = () => updatePasswordDialog.value.show();
const showWxBindDialog = () => wxBindDialog.value.show();

async function requestEmailVerify() {
  const loading = useLoading();
  try {
    await useHttpPost("/api/user/send_verify_email");
    useMsgSuccess(
      t("user.profile.account.emailVerifySuccess", { email: user.value.email })
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

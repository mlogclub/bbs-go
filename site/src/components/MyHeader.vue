<template>
  <nav class="navbar has-shadow" role="navigation" aria-label="main navigation">
    <div class="container">
      <div class="navbar-brand">
        <nuxt-link to="/" class="navbar-item">
          <img :alt="config.siteTitle" src="~/assets/images/logo.png" />
        </nuxt-link>
        <a
          :class="{ 'is-active': navbarActive }"
          class="navbar-burger burger"
          data-target="navbarBasic"
          @click="toggleNav"
        >
          <span aria-hidden="true" />
          <span aria-hidden="true" />
          <span aria-hidden="true" />
        </a>
      </div>
      <div :class="{ 'is-active': navbarActive }" class="navbar-menu">
        <div class="navbar-start">
          <nuxt-link
            v-for="(nav, index) in config.siteNavs"
            :key="index"
            :to="nav.url"
            class="navbar-item"
          >
            {{ nav.title }}
          </nuxt-link>
        </div>

        <div class="navbar-end">
          <div class="navbar-item">
            <search-input />
            <!-- <nuxt-link to="/search">xxxx</nuxt-link> -->
          </div>

          <div class="navbar-item">
            <create-topic-btn />
          </div>

          <msg-notice v-if="user" />

          <div
            v-if="user"
            class="navbar-item has-dropdown is-hoverable user-menus"
          >
            <div class="navbar-link">
              <MyAvatar :user="user" :size="24" />
              <span
                :to="`/user/${user.id}`"
                class="user-menus-nickname ellipsis"
                >{{ user.nickname }}</span
              >
            </div>
            <div class="navbar-dropdown">
              <nuxt-link class="navbar-item" :to="`/user/${user.id}`">
                <i class="iconfont icon-username" />
                <span>个人中心</span>
              </nuxt-link>
              <nuxt-link class="navbar-item" to="/user/favorites">
                <i class="iconfont icon-favorite" />
                <span>我的收藏</span>
              </nuxt-link>
              <nuxt-link class="navbar-item" to="/user/profile">
                <i class="iconfont icon-edit" />
                <span>编辑资料</span>
              </nuxt-link>
              <a class="navbar-item" @click="signout">
                <i class="iconfont icon-log-out" />
                <span>退出登录</span>
              </a>
            </div>
          </div>
          <div v-else class="navbar-item">
            <div class="buttons">
              <nuxt-link class="button login-btn" to="/user/signin">
                登录
              </nuxt-link>
            </div>
          </div>
          <div class="navbar-item">
            <color-mode />
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
const userStore = useUserStore();
const configStore = useConfigStore();

const { user } = storeToRefs(userStore);
const { config } = storeToRefs(configStore);

const navbarActive = ref(false);

function toggleNav() {
  navbarActive.value = !navbarActive.value;
}

async function signout() {
  if (confirm("确定退出登录吗？")) {
    await userStore.signout();
    useLinkTo("/");
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  /*opacity: 0.99;*/
  /*border-bottom: 1px solid #e7edf3;*/
  background-color: var(--bg-color);

  .navbar-item {
    font-weight: 700;
  }

  .publish {
    color: var(--text-color);
    background-color: #3174dc;
    width: 100px;
    &:hover {
      color: var(--text-color);
      background-color: #4d91fa;
    }
  }

  .login-btn {
    height: 32px;
    border-color: #000; // TODO
    &:hover {
      color: var(--text-color3);
      border-color: var(--text-color3);
    }
  }
}

.user-menus {
  .navbar-link {
    display: flex;
    align-items: center;

    .user-menus-nickname {
      margin-left: 5px;
      padding: 0 4px;
      font-size: 14px;
      color: var(--text-color);

      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
  .navbar-dropdown {
    border: 1px solid var(--border-color);

    a {
      display: flex;
      align-items: center;
      // padding: 8px 16px;
      img {
        width: 20px;
        height: 20px;
      }
      span {
        margin-left: 10px;
        width: 56px;
        height: 20px;
        font-size: 14px;
        font-weight: 400;
        line-height: 20px;
      }
    }
  }
}
</style>

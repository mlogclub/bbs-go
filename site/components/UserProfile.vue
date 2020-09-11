<template>
  <div class="profile">
    <img :src="user.avatar" class="avatar" />
    <div class="meta">
      <h1>
        <a :href="'/user/' + user.userId">{{ user.nickname }}</a>
      </h1>
      <div v-if="user.description" class="description">
        <p>{{ user.description }}</p>
      </div>
      <div v-if="user.homePage" class="homepage">
        <i class="iconfont icon-home"></i>
        <a :href="user.homePage" target="_blank" rel="external nofollow">{{
          user.homePage
        }}</a>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
}
</script>

<style lang="scss" scoped>
.profile {
  background: #fff;
  padding: 10px;
  margin-bottom: 10px;
  display: flex;

  .avatar {
    max-width: 66px;
    max-height: 66px;
    min-width: 66px;
    min-height: 66px;
  }

  .meta {
    margin-left: 18px;

    i {
      color: #cccccc;
      margin-right: 6px;
    }

    h1 {
      font-size: 24px;
      margin-bottom: 6px;
      a {
        color: #000;
        &:hover {
          color: #3273dc;
          text-decoration: underline;
        }
      }
    }

    .description {
      font-size: 14px;
      color: #555;
      margin-bottom: 6px;
    }

    .homepage {
      font-size: 14px;
      a {
        color: #666666;
        &:hover {
          color: #004e31;
          text-decoration: underline;
        }
      }
    }
  }
}
</style>

<template>
  <ul v-if="users && users.length" class="users">
    <li v-for="user in users" :key="user.id" class="user">
      <div
        :style="{ backgroundImage: 'url(' + user.avatar + ')' }"
        class="avatar small-avatar has-border is-rounded"
      />
      <div class="user-info">
        <a :href="'/user/' + user.id" :title="user.nickname" class="nickname">{{
          user.nickname
        }}</a>
        <span class="create-time">{{ user.createTime | prettyDate }}</span>
      </div>
    </li>
  </ul>
</template>

<script>
export default {
  props: {
    users: {
      type: Array,
      default() {
        return []
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.users {
  .user {
    display: flex;
    font-size: 14px;
    padding: 5px 0;
    &:not(:last-child) {
      border-bottom: 1px solid #f7f7f7;
    }

    .user-info {
      margin-left: 10px;
      .nickname {
        font-weight: bold;
        font-size: 15px;

        overflow: hidden;
        word-break: break-all;
        text-overflow: ellipsis;
        -webkit-line-clamp: 1;
        -webkit-box-orient: vertical;
        display: -webkit-box;
      }

      .create-time {
        font-size: 13px;
        color: #666;
      }
    }
  }
}
</style>

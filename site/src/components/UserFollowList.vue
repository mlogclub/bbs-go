<template>
  <div>
    <div v-if="users && users.length">
      <div v-for="item in users" :key="item.id" class="user-follow-item">
        <my-avatar :user="item" :size="40" has-border />
        <div class="user-follow-item-info">
          <div class="nickname">
            <nuxt-link :to="'/user/' + item.id">{{ item.nickname }}</nuxt-link>
          </div>
          <div class="description">
            {{ item.description }}{{ item.description }}{{ item.description }}
          </div>
        </div>
        <div>
          <follow-btn
            :user-id="item.id"
            :followed="item.followed"
            @onFollowed="onFollowed"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    users: {
      type: Array,
      default() {
        return [];
      },
    },
  },
  methods: {
    onFollowed(userId, followed) {
      this.$emit("onFollowed", userId, followed);
    },
  },
};
</script>

<style lang="scss" scoped>
.user-follow-item {
  display: flex;
  align-items: center;
  &:not(:last-child) {
    margin: 10px 0;
  }
  .user-follow-item-info {
    width: 100%;
    margin: auto 10px;

    .nickname {
      font-size: 14px;
      a {
        color: var(--text-color);
      }
    }
    .description {
      font-size: 12px;
      color: var(--text-color3);

      overflow: hidden;
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 1;
      text-align: justify;
      word-break: break-all;
      text-overflow: ellipsis;
    }
  }
}
</style>

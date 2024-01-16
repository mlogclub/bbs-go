<template>
  <div class="widget">
    <div class="widget-header">个人成就</div>
    <div class="widget-content extra-info">
      <ul>
        <li>
          <span>积分</span><br />
          <nuxt-link to="/user/scores">
            <b>{{ user.score }}</b>
          </nuxt-link>
        </li>
        <li>
          <span>话题</span><br />
          <b>{{ user.topicCount }}</b>
        </li>
        <li>
          <span>评论</span><br />
          <b>{{ user.commentCount }}</b>
        </li>
        <li>
          <span>注册排名</span><br />
          <b>{{ user.id }}</b>
        </li>
      </ul>
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
};
</script>

<style lang="scss" scoped>
.extra-info {
  ul {
    display: flex;
    li {
      width: 100%;
      text-align: center;
      span {
        font-size: 13px;
        font-weight: 400;
        color: var(--text-color3);
      }
      a,
      b {
        color: var(--text-color);
      }
    }
  }
}
</style>

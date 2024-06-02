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

<template>
  <div class="user-simple">
    <div class="base-info">
      <my-avatar :user="user" :size="80" :extra-style="{ margin: '0 auto' }" />
      <div class="nickname">
        <nuxt-link :to="`/user/${user.id}`" :alt="user.nickname">
          {{ user.nickname }}
        </nuxt-link>
      </div>
      <div class="description">
        {{ user.description }}
      </div>
    </div>
    <div class="extra-info">
      <ul class="extra-data">
        <li>
          <span>积分</span><br />
          <b>{{ user.score }}</b>
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

<style lang="scss" scoped>
.user-simple {
  background-color: var(--bg-color);
  padding: 0;
  margin: 0 0 10px 0;
  border-radius: 3px;

  .user-background {
    text-align: center;
    // background-image: url('https://file.mlog.club/images/2020/10/13/6e7933f5c9b2fe515210a17ea1762105.jpg!768_auto');
    background-size: cover;
    background-position: 50%;
  }

  .base-info {
    padding: 10px;
    text-align: center;

    .nickname {
      font-size: 15px;
      font-weight: 700;
      margin: 10px auto;
      a:hover {
        text-decoration: underline;
      }
    }

    .description {
      text-align: center;
      font-size: 13px;
      margin-top: 5px;
      overflow: hidden;
      word-break: break-all;
      text-overflow: ellipsis;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      display: -webkit-box;
    }
  }

  .extra-info {
    padding: 5px 0;
    background: rgba(0, 0, 0, 0.01);
    border-top: 1px solid var(--border-color);
    ul.extra-data {
      display: flex;
      li {
        width: 100%;
        text-align: center;
        span {
          font-size: 13px;
          font-weight: 400;
          color: var(--text-color3);
        }
      }
    }
  }
}
</style>

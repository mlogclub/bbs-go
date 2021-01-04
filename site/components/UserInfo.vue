<template>
  <div class="user-simple">
    <div class="base-info">
      <avatar :user="user" size="80" :extra-style="{ margin: '0 auto' }" />
      <div class="nickname">
        <a :href="'/user/' + user.id" :alt="user.nickname">{{
          user.nickname
        }}</a>
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
          <span>注册排名</span><br />
          <b>{{ user.id }}</b>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import Avatar from '~/components/Avatar'
export default {
  components: { Avatar },
  props: {
    user: {
      type: Object,
      required: true,
    },
  },
}
</script>
<style lang="scss" scoped>
.user-simple {
  background: #fff;
  padding: 0;
  margin: 0 0 10px 0;

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
    padding: 0 10px;
    background: rgba(0, 0, 0, 0.01);
    border-top: 1px solid #f5f5f5;
    ul.extra-data {
      display: flex;
      li {
        width: 100%;
        text-align: center;
        span {
          font-size: 13px;
          font-weight: 400;
          color: #868e96;
        }
      }
    }
  }
}
</style>

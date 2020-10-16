<template>
  <div class="user-simple">
    <div class="base-info">
      <a :href="'/user/' + user.id" :alt="user.nickname">
        <img :src="user.smallAvatar" :alt="user.nickname" class="avatar" />
      </a>
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
.user-simple {
  background: #fff;
  padding: 0;
  margin: 10px 0;

  .user-background {
    text-align: center;
    background-image: url('https://file.mlog.club/images/2020/10/13/6e7933f5c9b2fe515210a17ea1762105.jpg!768_auto');
    background-size: cover;
    background-position: 50%;
  }

  .base-info {
    padding: 10px;
    text-align: center;

    .avatar {
      min-width: 80px;
      min-height: 80px;
      width: 80px;
      height: 80px;
    }

    .nickname {
      font-size: 15px;
      font-weight: 700;
      a:hover {
        text-decoration: underline;
      }
    }

    .description {
      text-align: left;
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

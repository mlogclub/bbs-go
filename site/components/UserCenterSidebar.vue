<template>
  <div>
    <div v-if="isOwner" class="widget">
      <div class="widget-content">
        <a class="button is-primary" href="/topic/create">
          <i class="iconfont icon-topic" />&nbsp;
          <strong>发表主题</strong>
        </a>
        <a class="button is-success" href="/article/create">
          <i class="iconfont icon-publish" />&nbsp;
          <strong>发表文章</strong>
        </a>
      </div>
    </div>

    <div class="widget">
      <div class="widget-header">
        个人资料
      </div>
      <div class="widget-content">
        <img :src="user.avatar" class="img-avatar" />
        <div class="nickname">
          <a :href="'/user/' + user.id">{{ user.nickname }}</a>
        </div>
        <div v-if="user.description" class="description">
          <p>{{ user.description }}</p>
        </div>
        <div v-if="user.type === 1">
          <img
            :src="
              'https://open.weixin.qq.com/qr/code?username=' + user.username
            "
          />
        </div>
        <ul v-if="isOwner" class="operations">
          <li>
            <i class="iconfont icon-edit" />
            <a href="/user/settings">&nbsp;编辑资料</a>
          </li>
          <li>
            <i class="iconfont icon-message" />
            <a href="/user/messages">&nbsp;消息</a>
          </li>
          <li>
            <i class="iconfont icon-favorites" />
            <a href="/user/favorites">&nbsp;收藏</a>
          </li>
        </ul>
        <!-- 展示广告288x288
        <div style="text-align: center;">
        <ins
          class="adsbygoogle"
          style="display:inline-block;width:288px;height:288px"
          data-ad-client="ca-pub-5683711753850351"
          data-ad-slot="4922900917"
        />
        <script>
          (adsbygoogle = window.adsbygoogle || []).push({});
        </script>
      </div>
      --></div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    user: {
      type: Object,
      required: true
    },
    currentUser: {
      type: Object,
      required: true
    }
  },
  computed: {
    isOwner() {
      return (
        this.user && this.currentUser && this.user.id === this.currentUser.id
      )
    }
  }
}
</script>

<style lang="scss" scoped>
.nickname {
  font-size: 18px;
  font-weight: bold;
  a {
    color: #3273dc;
  }
}
.img-avatar {
  margin-top: 5px;
  border: 1px dotted #eeeeee;
  border-radius: 5%;
  width: 190px;
  height: 190px;
}
.description {
  font-size: 14px;
  padding: 10px 15px;
  border: 1px dotted #eeeeee;
  border-left: 3px solid #eeeeee;
  background-color: #fbfbfb;
}
.operations {
  list-style: none;
  margin-left: 0px;

  li {
    padding-left: 3px;

    font-size: 13px;
    &:hover {
      cursor: pointer;
      background-color: #fcf8e3;
      color: #8a6d3b;
      font-weight: bold;
    }

    a {
      color: #3273dc;
    }
  }
}
</style>

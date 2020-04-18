<template>
  <ul class="tweets">
    <li v-for="tweet in tweets" :key="tweet.tweetId">
      <div class="tweet">
        <div class="pin-header-row">
          <div class="account-group">
            <div>
              <a :href="'/user/' + tweet.user.id" :title="tweet.user.nickname">
                <img v-lazy="tweet.user.avatar" class="avatar size-45" />
              </a>
            </div>
            <div class="pin-header-content">
              <div>
                <a
                  :href="'/user/' + tweet.user.id"
                  :title="tweet.user.nickname"
                  target="_blank"
                  class="nickname"
                  >{{ tweet.user.nickname }}</a
                >
              </div>
              <div class="meta-box">
                <div class="position ellipsis">
                  {{ tweet.user.description }}
                </div>
                <div class="dot">·</div>
                <time
                  datetime="2020-04-17T17:26:43.760Z"
                  title="Sat Apr 18 2020 01:26:43 GMT+0800 (中国标准时间)"
                  class="time"
                  >12小时前
                </time>
              </div>
            </div>
          </div>
        </div>
        <div class="pin-content-row">
          <a :href="'/tweet/' + tweet.tweetId" class="content-box">{{
            tweet.content
          }}</a>
        </div>
        <ul
          v-if="tweet.imageList && tweet.imageList.length > 0"
          class="pin-image-row"
        >
          <li v-for="image in tweet.imageList" :key="image">
            <a :href="'/tweet/' + tweet.tweetId" class="image-item">
              <img v-lazy="image" />
            </a>
          </li>
        </ul>
        <div class="pin-action-row">
          <div class="action-box">
            <div class="like-action action">
              <div class="action-title-box">
                <i class="iconfont icon-like" />
                <span class="action-title">{{
                  tweet.likeCount > 0 ? tweet.likeCount : '赞'
                }}</span>
              </div>
            </div>
            <div class="comment-action action">
              <div class="action-title-box">
                <i class="iconfont icon-comment" />
                <span class="action-title">{{
                  tweet.commentCount > 0 ? tweet.commentCount : '评论'
                }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </li>
  </ul>
</template>

<script>
export default {
  props: {
    tweets: {
      type: Array,
      default() {
        return []
      },
      required: false
    }
  }
}
</script>

<style scoped lang="scss"></style>

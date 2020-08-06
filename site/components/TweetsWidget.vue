<template>
  <div v-if="tweets && tweets.length" class="widget">
    <div class="widget-header">
      最新动态<span class="slot"><a href="/tweets">更多</a></span>
    </div>
    <div class="widget-content">
      <ul class="tweets-widget">
        <li v-for="tweet in tweets" :key="tweet.tweetId">
          <a :href="'/tweet/' + tweet.tweetId" class="tweet-content">
            {{ tweet.content }}
            <i
              v-if="tweet.imageList && tweet.imageList.length"
              class="iconfont icon-image"
            />
          </a>
          <div class="tweet-metas">
            喜欢 {{ tweet.likeCount }}&nbsp;&nbsp;&nbsp; 评论
            {{ tweet.commentCount }}
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    tweets: {
      type: Array,
      default() {
        return []
      },
      required: false,
    },
  },
}
</script>

<style scoped lang="scss">
.tweets-widget {
  li {
    font-size: 14px;
    line-height: 22px;
    margin-bottom: 4px;
    &:not(:last-child) {
      margin-bottom: 10px;
    }

    a.tweet-content {
      color: #2d2d2d;
      text-decoration: none;
      background-color: transparent;
      outline: none;
      cursor: pointer;
      transition: color 0.3s;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;

      &:hover {
        color: #ec7259;
      }
    }
    .tweet-metas {
      font-size: 12px;
      color: #969696;
      user-select: none;
    }
  }
}
</style>

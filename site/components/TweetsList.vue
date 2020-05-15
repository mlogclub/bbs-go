<template>
  <ul class="tweets">
    <li v-for="tweet in tweets" :key="tweet.tweetId">
      <div class="tweet">
        <div class="pin-header-row">
          <div class="account-group">
            <div>
              <a :href="'/user/' + tweet.user.id" :title="tweet.user.nickname">
                <img :src="tweet.user.smallAvatar" class="avatar size-45" />
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
                  :datetime="
                    tweet.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')
                  "
                  itemprop="datePublished"
                  >{{ tweet.createTime | prettyDate }}</time
                >
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
          <li v-for="(image, index) in tweet.imageList" :key="image + index">
            <a :href="'/tweet/' + tweet.tweetId" class="image-item">
              <img v-lazy="image.preview" />
            </a>
          </li>
        </ul>
        <div class="pin-action-row">
          <div class="action-box">
            <div class="like-action action" @click="like(tweet)">
              <div class="action-title-box">
                <i class="iconfont icon-like" />
                <span class="action-title">{{
                  tweet.likeCount > 0 ? tweet.likeCount : '赞'
                }}</span>
              </div>
            </div>
            <a :href="'/tweet/' + tweet.tweetId" class="comment-action action">
              <div class="action-title-box">
                <i class="iconfont icon-comments" />
                <span class="action-title">{{
                  tweet.commentCount > 0 ? tweet.commentCount : '评论'
                }}</span>
              </div>
            </a>
          </div>
        </div>
      </div>
    </li>
  </ul>
</template>

<script>
import utils from '~/common/utils'
export default {
  props: {
    tweets: {
      type: Array,
      default() {
        return []
      },
      required: false
    }
  },
  methods: {
    async like(tweet) {
      try {
        await this.$axios.post('/api/tweet/like/' + tweet.tweetId)
        tweet.liked = true
        tweet.likeCount++
      } catch (e) {
        if (e.errorCode === 1) {
          this.$toast.info('请登录后点赞！！！', {
            action: {
              text: '去登录',
              onClick: (e, toastObject) => {
                utils.toSignin()
              }
            }
          })
        } else {
          this.$toast.error(e.message || e)
        }
      }
    }
  }
}
</script>

<style scoped lang="scss"></style>

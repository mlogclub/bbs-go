<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <div class="tweet">
          <div class="pin-header-row">
            <div class="account-group">
              <div>
                <a
                  :href="'/user/' + tweet.user.id"
                  :title="tweet.user.nickname"
                >
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
            <div class="content-box">{{ tweet.content }}</div>
          </div>
          <ul
            v-if="tweet.imageList && tweet.imageList.length > 0"
            v-viewer
            class="pin-image-row"
          >
            <li v-for="image in tweet.imageList" :key="image">
              <div class="image-item">
                <img :src="image.preview" :data-src="image.url" />
              </div>
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
              <div class="comment-action action">
                <div class="action-title-box">
                  <i class="iconfont icon-comments" />
                  <span class="action-title">{{
                    tweet.commentCount > 0 ? tweet.commentCount : '评论'
                  }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <comment
          :entity-id="tweet.tweetId"
          :comments-page="commentsPage"
          mode="text"
          entity-type="tweet"
        />
      </div>
      <div class="right-container">
        <site-notice />
        <score-rank :score-rank="scoreRank" />
      </div>
    </div>
  </section>
</template>

<script>
import Vue from 'vue'
import Viewer from 'v-viewer'
import SiteNotice from '~/components/SiteNotice'
import ScoreRank from '~/components/ScoreRank'
import Comment from '~/components/Comment'
import utils from '~/common/utils'
import 'viewerjs/dist/viewer.css'

Vue.use(Viewer, {
  defaultOptions: {
    zIndex: 9999,
    navbar: false,
    title: false,
    tooltip: false,
    movable: false,
    scalable: false,
    url: 'data-src',
  },
})

export default {
  components: {
    SiteNotice,
    ScoreRank,
    Comment,
  },
  async asyncData({ $axios, params }) {
    const [tweet, commentsPage, scoreRank] = await Promise.all([
      $axios.get('/api/tweet/' + params.id),
      $axios.get('/api/comment/list', {
        params: {
          entityType: 'tweet',
          entityId: params.id,
        },
      }),
      $axios.get('/api/user/score/rank'),
    ])
    return { tweet, commentsPage, scoreRank }
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
              },
            },
          })
        } else {
          this.$toast.error(e.message || e)
        }
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle('动态'),
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: this.$siteDescription(),
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() },
      ],
    }
  },
}
</script>

<style lang="scss" scoped></style>

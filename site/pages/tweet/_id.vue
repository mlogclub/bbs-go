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
            class="pin-image-row"
          >
            <li v-for="image in tweet.imageList" :key="image">
              <div class="image-item">
                <img v-lazy="image" />
              </div>
            </li>
          </ul>
          <div class="pin-action-row">
            <div class="action-box">
              <div class="like-action action">
                <div @click="like(tweet)" class="action-title-box">
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

        <!-- 评论 -->
        <comment
          :entity-id="tweet.tweetId"
          :comments-page="commentsPage"
          entity-type="tweet"
        />
      </div>
      <topic-side :score-rank="scoreRank" :links="links" />
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'
import TopicSide from '~/components/TopicSide'
import Comment from '~/components/Comment'

export default {
  components: {
    TopicSide,
    Comment
  },
  async asyncData({ $axios, params, error }) {
    const [tweet, commentsPage, scoreRank, links] = await Promise.all([
      $axios.get('/api/tweet/' + params.id),
      $axios.get('/api/comment/list', {
        params: {
          entityType: 'tweet',
          entityId: params.id
        }
      }),
      $axios.get('/api/user/score/rank'),
      $axios.get('/api/link/toplinks')
    ])
    return { tweet, commentsPage, scoreRank, links }
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
  },
  head() {
    return {
      title: this.$siteTitle('动态'),
      meta: [
        {
          hid: 'description',
          name: 'description',
          content: this.$siteDescription()
        },
        { hid: 'keywords', name: 'keywords', content: this.$siteKeywords() }
      ]
    }
  }
}
</script>

<style lang="scss" scoped></style>

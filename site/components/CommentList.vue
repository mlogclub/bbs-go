<template>
  <div class="comments">
    <load-more
      v-if="commentsPage"
      ref="commentsLoadMore"
      v-slot="{ results }"
      :init-data="commentsPage"
      :params="{ entityType: entityType, entityId: entityId }"
      url="/api/comment/list"
    >
      <ul>
        <li
          v-for="(comment, index) in results"
          :key="comment.commentId"
          class="comment"
          itemprop="comment"
          itemscope
          itemtype="http://schema.org/Comment"
        >
          <adsbygoogle
            v-if="showAd && (index + 1) % 3 === 0 && index !== 0"
            ad-slot="4980294904"
            ad-format="fluid"
            ad-layout-key="-ht-19-1m-3j+mu"
          />
          <div class="comment-avatar">
            <avatar :user="comment.user" size="35" />
          </div>
          <div class="comment-meta">
            <span
              class="comment-nickname"
              itemprop="creator"
              itemscope
              itemtype="http://schema.org/Person"
            >
              <a :href="'/user/' + comment.user.id" itemprop="name">
                {{ comment.user.nickname }}
              </a>
            </span>
            <span class="comment-time">
              <time
                :datetime="
                  comment.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')
                "
                itemprop="datePublished"
                >{{ comment.createTime | prettyDate }}</time
              >
            </span>
            <span class="comment-reply">
              <a @click="reply(comment)">回复</a>
            </span>
          </div>
          <div class="comment-content content">
            <blockquote v-if="comment.quote" class="comment-quote">
              <div class="comment-quote-user">
                <avatar :user="comment.quote.user" size="20" />
                <a class="quote-nickname">{{ comment.quote.user.nickname }}</a>
                <span class="quote-time">
                  {{ comment.quote.createTime | prettyDate }}
                </span>
              </div>
              <div
                v-lazy-container="{ selector: 'img' }"
                itemprop="text"
                v-html="comment.quote.content"
              />
            </blockquote>
            <p
              v-lazy-container="{ selector: 'img' }"
              v-html="comment.content"
            />
          </div>
        </li>
      </ul>
    </load-more>
  </div>
</template>

<script>
import Avatar from '~/components/Avatar'
import LoadMore from '~/components/LoadMore'

export default {
  components: {
    Avatar,
    LoadMore,
  },
  props: {
    entityType: {
      type: String,
      default: '',
      required: true,
    },
    entityId: {
      type: Number,
      default: 0,
      required: true,
    },
    commentsPage: {
      type: Object,
      default() {
        return {}
      },
    },
    showAd: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    user() {
      return this.$store.state.user.current
    },
    isLogin() {
      return this.$store.state.user.current != null
    },
  },
  methods: {
    append(data) {
      if (!data) return

      this.$refs.commentsLoadMore.unshiftResults(data)
    },
    reply(quote) {
      if (!this.isLogin) {
        this.$toSignin()
      }
      this.$emit('reply', quote)
    },
    cancelReply() {
      this.quote = null
    },
  },
}
</script>

<style scoped lang="scss"></style>

<template>
  <div>
    <div class="comment-header">
      评论<span v-if="commentCount > 0">({{ commentCount }})</span>
    </div>
    <comment-input
      v-if="mode === 'markdown'"
      ref="input"
      :entity-id="entityId"
      :entity-type="entityType"
      @created="commentCreated"
    />
    <comment-text-input
      v-else
      ref="input"
      :entity-id="entityId"
      :entity-type="entityType"
      @created="commentCreated"
    />

    <comment-list
      ref="list"
      :entity-id="entityId"
      :entity-type="entityType"
      :comments-page="commentsPage"
      @reply="reply"
    />
  </div>
</template>

<script>
import CommentList from '~/components/CommentList'
import CommentInput from '~/components/CommentInput'
import CommentTextInput from '~/components/CommentTextInput'
export default {
  name: 'Comment',
  components: {
    CommentList,
    CommentInput,
    CommentTextInput,
  },
  props: {
    mode: {
      type: String,
      default: 'markdown',
    },
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
    commentCount: {
      type: Number,
      default: 0,
    },
    showAd: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    commentCreated(data) {
      this.$refs.list.append(data)
    },
    reply(quote) {
      this.$refs.input.reply(quote)
    },
  },
}
</script>
<style lang="scss" scoped>
.comment-header {
  display: flex;
  padding-top: 20px;
  margin: 0 10px;
  border-top: 1px solid rgba(228, 228, 228, 0.6);
  color: #6d6d6d;
  font-size: 16px;
}
</style>

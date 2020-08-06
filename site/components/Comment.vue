<template>
  <div>
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

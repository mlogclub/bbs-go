<template>
  <topic-list
    v-if="topics && topics.length"
    class="sticky-topics"
    show-sticky
    :topics="topics"
  />
</template>

<script>
export default {
  props: {
    nodeId: {
      type: Number,
      default: 0,
    },
  },
  data() {
    return {
      topics: [],
    }
  },
  mounted() {
    this.loadStickyTopics()
  },
  methods: {
    async loadStickyTopics() {
      try {
        this.topics = await this.$axios.get('/api/topic/sticky_topics', {
          params: {
            nodeId: this.nodeId,
          },
        })
      } catch (e) {
        console.error(e)
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.sticky-topics {
  margin-bottom: 10px;
}
</style>

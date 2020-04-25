<template>
  <div class="right-container">
    <site-notice />
    <score-rank :score-rank="scoreRank" />
    <friend-links :links="links" />
  </div>
</template>

<script>
import SiteNotice from '~/components/SiteNotice'
import ScoreRank from '~/components/ScoreRank'
import FriendLinks from '~/components/FriendLinks'

export default {
  components: { SiteNotice, ScoreRank, FriendLinks },
  props: {
    currentNodeId: {
      type: Number,
      default: 0
    },
    links: {
      type: Array,
      default() {
        return null
      }
    },
    scoreRank: {
      type: Array,
      default() {
        return null
      }
    }
  }
}
</script>

<style lang="scss" scoped></style>

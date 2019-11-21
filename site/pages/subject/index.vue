<template>
  <section class="main">
    <div class="container-wrapper main-container right-main">
      <subject-bar />
      <div class="right-container">
        <load-more
          v-slot="{ results }"
          :init-data="subjectContentPage"
          url="/api/subject/contents"
        >
          <subject-content-list :subject-contents="results" />
        </load-more>
      </div>
    </div>
  </section>
</template>

<script>
import SubjectBar from '~/components/SubjectBar'
import SubjectContentList from '~/components/SubjectContentList'
import LoadMore from '~/components/LoadMore'
export default {
  components: {
    SubjectBar,
    SubjectContentList,
    LoadMore
  },
  async asyncData({ $axios, params, query }) {
    const [subjectContentPage] = await Promise.all([
      $axios.get('/api/subject/contents')
    ])
    return {
      subjectContentPage
    }
  },
  head() {
    return {
      title: this.$siteTitle('专栏')
    }
  }
}
</script>

<style lang="scss" scoped></style>

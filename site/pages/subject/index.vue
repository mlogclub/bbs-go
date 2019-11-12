<template>
  <section class="main">
    <div class="container main-container right-main">
      <subject-bar />
      <div class="m-right">
        <subject-content-list :subject-contents="subjectContentPage.results" />
        <pagination url-prefix="/subject?page=" :page="subjectContentPage.page" />
      </div>
    </div>
  </section>
</template>

<script>
import SubjectBar from '~/components/SubjectBar'
import SubjectContentList from '~/components/SubjectContentList'
import Pagination from '~/components/Pagination'
export default {
  components: {
    SubjectBar, SubjectContentList, Pagination
  },
  head() {
    return {
      title: this.$siteTitle('专栏')
    }
  },
  async asyncData({ $axios, params, query }) {
    const [subjectContentPage] = await Promise.all([
      $axios.get('/api/subject/contents', {
        params: {
          page: query.page || 1
        }
      })
    ])
    return {
      subjectContentPage: subjectContentPage
    }
  }
}
</script>

<style lang="scss" scoped>

</style>

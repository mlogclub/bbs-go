<template>
  <section class="main">
    <div class="container main-container">
      <div class="right-main-container">
        <subject-bar :subject="subject" />
        <div class="m-right">
          <subject-content-list :subject-contents="subjectContentPage.results" />
          <pagination :url-prefix="'/subject/' + subject.id + '?page='" :page="subjectContentPage.page" />
        </div>
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
    const title = this.subject.title + ' 专栏'
    const description = this.subject.description || ''
    return {
      title: this.$siteTitle(title),
      meta: [
        { hid: 'description', name: 'description', content: description }
      ]
    }
  },
  async asyncData({ $axios, params, query }) {
    const [subject, subjectContentPage] = await Promise.all([
      $axios.get('/api/subject/' + params.id),
      $axios.get('/api/subject/contents', {
        params: {
          subjectId: params.id,
          page: query.page || 1
        }
      })
    ])
    return {
      subject: subject,
      subjectContentPage: subjectContentPage
    }
  }
}
</script>

<style lang="scss" scoped>

</style>

<template>
  <section class="main">
    <div class="container-wrapper main-container right-main">
      <subject-bar :subject="subject" />
      <div class="right-container">
        <load-more
          v-slot="{ results }"
          :init-data="subjectContentPage"
          :params="{ subjectId: subject.id }"
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
    const [subject, subjectContentPage] = await Promise.all([
      $axios.get('/api/subject/' + params.id),
      $axios.get('/api/subject/contents', {
        params: {
          subjectId: params.id
        }
      })
    ])
    return {
      subject,
      subjectContentPage
    }
  },
  head() {
    const title = this.subject.title + ' 专栏'
    const description = this.subject.description || ''
    return {
      title: this.$siteTitle(title),
      meta: [{ hid: 'description', name: 'description', content: description }]
    }
  }
}
</script>

<style lang="scss" scoped></style>

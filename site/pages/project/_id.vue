<template>
  <section class="main">
    <div class="container main-container is-white left-main">
      <div class="left-container">
        <div v-if="project" class="project">
          <div class="project-header">
            <h1>
              <span class="project-name">{{ project.name }}</span>
              <span v-if="project.title" class="project-title"
                >&nbsp;-&nbsp;{{ project.title }}</span
              >
            </h1>
            <div class="project-meta">
              <span>
                <nuxt-link :to="'/user/' + project.user.id">{{
                  project.user.nickname
                }}</nuxt-link>
              </span>
              <span>{{ project.createTime | prettyDate }}</span>
            </div>
          </div>

          <div class="ad">
            <!-- 展示广告 -->
            <adsbygoogle ad-slot="1742173616" />
          </div>

          <div
            v-lazy-container="{ selector: 'img' }"
            class="content"
            v-html="project.content"
          ></div>

          <div class="footer">
            <a
              v-if="projectUrl"
              :href="projectUrl"
              class="homepage"
              target="_blank"
              >项目主页</a
            >
            <a v-if="docUrl" :href="docUrl" class="homepage" target="_blank"
              >文档地址</a
            >
            <a
              v-if="downloadUrl"
              :href="downloadUrl"
              class="homepage"
              target="_blank"
              >下载地址</a
            >
          </div>
        </div>

        <!-- 评论 -->
        <comment
          :entity-id="project.projectId"
          :comments-page="commentsPage"
          :show-ad="true"
          entity-type="project"
        />
      </div>
      <div class="right-container">
        <site-notice />

        <div class="ad">
          <!-- 展示广告 -->
          <adsbygoogle ad-slot="1742173616" />
        </div>

        <div class="ad">
          <!-- 展示广告 -->
          <adsbygoogle ad-slot="1742173616" />
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  async asyncData({ $axios, params, store }) {
    const [project, commentsPage] = await Promise.all([
      $axios.get('/api/project/' + params.id),
      $axios.get('/api/comment/comments', {
        params: {
          entityType: 'project',
          entityId: params.id,
        },
      }),
    ])
    // 构建url，如果登录了直接跳转到原地址，如果没登陆那么跳转到登录
    function buildUrl(url) {
      if (!url || !project) {
        return ''
      }
      if (store.state.user.current) {
        // 如果用户登录了
        return '/redirect?url=' + encodeURI(url)
      } else {
        // 没登陆，引导跳转登录
        return '/user/signin?ref=' + encodeURI('/project/' + project.projectId)
      }
    }
    return {
      project,
      commentsPage,
      projectUrl: buildUrl(project.url),
      docUrl: buildUrl(project.docUrl),
      downloadUrl: buildUrl(project.downloadUrl),
    }
  },
  head() {
    let siteTitle = this.project.name
    if (this.project.title) {
      siteTitle += ' - ' + this.project.title
    }
    return {
      title: this.$siteTitle(siteTitle),
      meta: [{ hid: 'description', name: 'description', content: siteTitle }],
    }
  },
}
</script>

<style lang="scss" scoped>
.project {
  margin-bottom: 10px;

  .project-header {
    padding: 10px 0;
    border-bottom: 1px solid var(--border-color);
    .project-name {
      font-size: 18px;
      font-weight: 700;
      color: rgba(0, 0, 0, 0.75);
    }

    .project-title {
      font-size: 16px;
      font-weight: 500;
      color: rgba(0, 0, 0, 0.6);
    }

    .project-meta {
      span {
        font-size: 14px;
        color: var(--text-color3);
        padding-top: 6px;

        &:not(:first-child) {
          margin-left: 10px;
        }

        a {
          color: var(--text-color3);
          cursor: pointer;

          &:hover {
            color: var(--text-link-color);
            font-weight: 500;
          }
        }
      }
    }
  }

  .footer {
    .homepage {
      border: 1px solid var(--border-color);
      padding: 10px;
      border-radius: 3px;
      font-weight: 900;
    }
  }
}
</style>

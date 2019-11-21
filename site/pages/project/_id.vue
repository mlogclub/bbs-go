<template>
  <section class="main">
    <div class="container-wrapper main-container left-main">
      <div class="left-container">
        <div v-if="project" class="project">
          <div class="project-header">
            <span class="project-name">{{ project.name }}</span>
            <span v-if="project.title" class="project-title"
              >&nbsp;-&nbsp;{{ project.title }}</span
            >
          </div>
          <div class="meta">
            <span>
              <a :href="'/user/' + project.user.id">{{
                project.user.nickname
              }}</a>
            </span>
            <span>{{ project.createTime | prettyDate }}</span>
          </div>
          <div class="content">
            <ins
              class="adsbygoogle"
              style="display:block"
              data-ad-format="fluid"
              data-ad-layout-key="-ig-s+1x-t-q"
              data-ad-client="ca-pub-5683711753850351"
              data-ad-slot="4728140043"
            />
            <script>
              ;(adsbygoogle = window.adsbygoogle || []).push({})
            </script>
            <p v-highlight v-html="project.content" />
          </div>
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
        <div style="max-height:60px;">
          <!-- 展示广告190x90 -->
          <ins
            class="adsbygoogle"
            style="display:inline-block;width:190px;height:90px"
            data-ad-client="ca-pub-5683711753850351"
            data-ad-slot="9345305153"
          />
          <script>
            ;(adsbygoogle = window.adsbygoogle || []).push({})
          </script>

          <!-- 展示广告190x190 -->
          <ins
            class="adsbygoogle"
            style="display:inline-block;width:190px;height:190px"
            data-ad-client="ca-pub-5683711753850351"
            data-ad-slot="5685455263"
          />
          <script>
            ;(adsbygoogle = window.adsbygoogle || []).push({})
          </script>

          <!-- 展示广告190x480 -->
          <ins
            class="adsbygoogle"
            style="display:inline-block;width:190px;height:480px"
            data-ad-client="ca-pub-5683711753850351"
            data-ad-slot="3438372357"
          />
          <script>
            ;(adsbygoogle = window.adsbygoogle || []).push({})
          </script>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import Comment from '~/components/Comment'
export default {
  components: {
    Comment
  },
  async asyncData({ $axios, params, store }) {
    const [project, commentsPage] = await Promise.all([
      $axios.get('/api/project/' + params.id),
      $axios.get('/api/comment/list', {
        params: {
          entityType: 'project',
          entityId: params.id
        }
      })
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
      downloadUrl: buildUrl(project.downloadUrl)
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.project.name),
      meta: [
        { hid: 'description', name: 'description', content: this.project.title }
      ]
    }
  }
}
</script>

<style lang="scss" scoped>
.project {
  margin-bottom: 10px;

  .project-header {
    .project-name {
      font-size: 18px;
      font-weight: 700;
      color: rgba(0, 0, 0, 0.75);
      margin-top: 5px;
      margin-bottom: 5px;
    }

    .project-title {
      font-size: 16px;
      font-weight: 400;
      color: rgba(0, 0, 0, 0.6);
    }
  }

  .meta {
    margin-bottom: 15px;
    border-bottom: 1px solid #f4f4f5;
    span {
      display: inline-block;
      font-size: 14px;
      color: #999;
      padding-top: 6px;

      a {
        color: #999;
        cursor: pointer;

        &:hover {
          color: #3273dc;
          font-weight: 500;
        }
      }
    }
  }

  .footer {
    .homepage {
      border: 1px solid #f4f4f5;
      padding: 10px;
      border-radius: 3px;
      font-weight: 900;
    }
  }
}
</style>

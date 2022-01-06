<template>
  <div class="projects">
    <ul>
      <li v-for="(p, index) in projects" :key="p.projectId" class="project">
        <article itemscope itemtype="http://schema.org/BlogPosting">
          <div
            v-if="index === 2 || index === 6 || index === 12 || index === 18"
          >
            <!-- 信息流广告 -->
            <adsbygoogle
              ad-slot="4980294904"
              ad-format="fluid"
              ad-layout-key="-ht-19-1m-3j+mu"
            />
          </div>
          <div class="project-header">
            <h1 itemprop="headline">
              <nuxt-link :to="'/project/' + p.projectId">
                <span class="project-name">{{ p.name }}</span>
                <span v-if="p.title" class="project-title"
                  >&nbsp;-&nbsp;{{ p.title }}</span
                >
              </nuxt-link>
            </h1>
          </div>
          <div class="summary" itemprop="description">
            {{ p.summary }}
          </div>
          <span class="meta">
            <span
              itemprop="author"
              itemscope
              itemtype="http://schema.org/Person"
            >
              <nuxt-link :to="'/user/' + p.user.id" itemprop="name">{{
                p.user.nickname
              }}</nuxt-link>
            </span>
            <span>
              <time
                :datetime="p.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')"
                itemprop="datePublished"
                >{{ p.createTime | prettyDate }}</time
              >
            </span>
          </span>
        </article>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  props: {
    projects: {
      type: Array,
      default() {
        return []
      },
    },
  },
}
</script>

<style lang="scss" scoped>
.projects {
  .project {
    &:not(:last-child) {
      border-bottom: 1px dashed var(--border-color);
    }
    padding-top: 5px;
    padding-bottom: 5px;

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

    .summary {
      font-size: 14px;
      color: rgba(0, 0, 0, 0.7);
      margin-top: 10px;
      margin-bottom: 10px;
    }

    .meta {
      span {
        display: inline-block;
        font-size: 13px;
        color: var(--text-color3);
        padding-top: 6px;

        &:not(:first-child) {
          margin-left: 8px;
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
}
</style>

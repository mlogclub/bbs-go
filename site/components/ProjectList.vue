<template>
  <div class="projects">
    <ul>
      <li v-for="p in projects" :key="p.projectId" class="project">
        <div class="project-header">
          <a :href="'/project/' + p.projectId">
            <span class="project-name">{{ p.name }}</span>
            <span v-if="p.title" class="project-title"
              >&nbsp;-&nbsp;{{ p.title }}</span
            >
          </a>
        </div>
        <div class="summary">
          {{ p.summary }}
        </div>
        <div class="meta">
          <span>
            <a :href="'/user/' + p.user.id">{{ p.user.nickname }}</a>
          </span>
          <span>{{ p.createTime | prettyDate }}</span>
        </div>
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
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.projects {
  .project {
    &:not(:last-child) {
      border-bottom: 1px dashed #f4f4f5;
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
        color: #999;
        padding-top: 6px;

        &:not(:first-child) {
          margin-left: 8px;
        }

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
  }
}
</style>

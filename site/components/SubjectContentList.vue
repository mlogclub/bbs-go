<template>
  <div class="subject-contents">
    <ul>
      <li
        v-for="(subjectContent, index) in subjectContents"
        :key="subjectContent.subjectContentId"
        class="subject-content"
      >
        <div class="content-title">
          <a :href="subjectContent.url" target="_blank">{{
            subjectContent.title
          }}</a>
        </div>
        <div class="content-summary">
          {{ subjectContent.summary }}
        </div>
        <div class="content-meta">
          于 {{ subjectContent.createTime | prettyDate }} 收录入专栏
          <a
            :href="'/subject/' + subjectContent.subjectId"
            class="subject-name"
            >{{ subjectContent.subjectTitle }}</a
          >
        </div>

        <div v-if="index === 1">
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
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  props: {
    subjectContents: {
      type: Array,
      required: true,
      default() {
        return []
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.subject-contents {
  .subject-content {
    padding: 8px 0;
    zoom: 1;
    position: relative;
    overflow: hidden;

    &:not(:last-child) {
      border-bottom: 1px solid #f2f2f2;
    }

    .content-title {
      a {
        font-size: 18px;
        line-height: 30px;
        font-weight: 500;
        color: #17181a;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }

    .content-summary {
      font-size: 14px;
      color: rgba(0, 0, 0, 0.7);
      overflow: hidden;
      padding-top: 6px;
    }

    .content-meta {
      font-size: 13px;
      color: #8a93a0;
      margin-top: 10px;
      a {
        font-weight: bold;
      }

      .subject-name {
        margin-left: 5px;
      }
    }
  }
}
</style>

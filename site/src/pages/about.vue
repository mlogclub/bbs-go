<template>
  <section class="main">
    <div class="container">
      <div class="widget">
        <div class="widget-header">{{ t("pages.about.title") }}</div>
        <div class="widget-content content">
          <p>
            {{ t("pages.about.desc") }}
          </p>
        </div>
      </div>

      <div class="widget">
        <div class="widget-header">{{ t("pages.about.featuresTitle") }}</div>
        <div class="widget-content content">
          <ul>
            <li>
              <strong>{{ t("pages.about.feature1Title") }}</strong
              >{{ t("pages.about.feature1") }}
            </li>
            <li>
              <strong>{{ t("pages.about.feature2Title") }}</strong
              >{{ t("pages.about.feature2") }}
            </li>
            <li>
              <strong>{{ t("pages.about.feature3Title") }}</strong
              >{{ t("pages.about.feature3") }}
            </li>
            <li>
              <strong>{{ t("pages.about.feature4Title") }}</strong
              >{{ t("pages.about.feature4") }}
            </li>
            <li>
              <strong>{{ t("pages.about.feature5Title") }}</strong
              >{{ t("pages.about.feature5") }}
            </li>
          </ul>
        </div>
      </div>

      <div class="widget">
        <div class="widget-header">{{ t("pages.about.joinTitle") }}</div>
        <div class="widget-content content">
          <p>
            {{ t("pages.about.joinDesc") }}
          </p>
          <ul>
            <li>
              <span>{{ t("pages.about.github") }}</span>
              <a href="https://github.com/mlogclub/bbs-go" target="_blank"
                >mlogclub/bbs-go</a
              >
            </li>
            <li>
              <span>{{ t("pages.about.gitee") }}</span>
              <a href="https://gitee.com/mlogclub/bbs-go" target="_blank"
                >mlogclub/bbs-go</a
              >
            </li>
            <li>
              <span>{{ t("pages.about.contributors") }}</span>
              <a
                href="https://github.com/mlogclub/bbs-go/graphs/contributors"
                target="_blank"
                >{{ t("pages.about.contributorsList") }}</a
              >
            </li>
          </ul>
        </div>
      </div>

      <div class="widget">
        <div class="widget-header">{{ t("pages.about.licenseTitle") }}</div>
        <div class="widget-content content">
          <p>
            {{ t("pages.about.licenseDesc") }}
            <a
              href="https://github.com/mlogclub/bbs-go/blob/master/LICENSE"
              target="_blank"
              >GNU General Public License v3.0</a
            >
            {{ t("pages.about.licenseDesc2") }}
          </p>
        </div>
      </div>

      <div class="widget">
        <div class="widget-header">{{ t("pages.about.contactTitle") }}</div>
        <div class="widget-content content">
          <p>
            {{ t("pages.about.contactDesc") }}
          </p>
          <ul>
            <li>
              <span>{{ t("pages.about.community") }}</span>
              <a href="https://bbs.bbs-go.com" target="_blank">BBS-GO社区</a>
            </li>
            <li>
              <span>{{ t("pages.about.docs") }}</span>
              <a href="https://bbs-go.com" target="_blank">bbs-go.com</a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { useI18n } from "vue-i18n";
const { t } = useI18n();
useHead({
  title: t("pages.about.title"),
});
</script>

<style lang="scss" scoped>
.widget {
  .widget-content {
    padding: 16px;
  }
}
</style>

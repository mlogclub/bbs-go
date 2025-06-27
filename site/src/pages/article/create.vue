<template>
  <section class="main">
    <div class="container">
      <article v-if="isNeedEmailVerify" class="message is-warning">
        <div class="message-header">
          <p>{{ $t("pages.article.create.needEmailTitle") }}</p>
        </div>
        <div class="message-body">
          {{ $t("pages.article.create.needEmailBody") }}
          <strong>
            <nuxt-link
              to="/user/profile/account"
              style="color: var(--text-link-color)"
              >{{ $t("pages.article.create.goVerify") }}</nuxt-link
            >
          </strong>
        </div>
      </article>
      <div v-else class="publish-form">
        <div class="form-title">
          <div class="form-title-name">
            {{ $t("pages.article.create.title") }}
          </div>
        </div>

        <div class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input"
              type="text"
              :placeholder="$t('pages.article.create.titlePlaceholder')"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <markdown-editor
              v-model="postForm.content"
              :placeholder="$t('pages.article.create.contentPlaceholder')"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <image-upload v-model="postForm.cover" :limit="1" size="120px">
              <template #add-image-button>
                <div class="cover-add-btn">
                  <i class="iconfont icon-add" />
                  <span>{{ $t("pages.article.create.cover") }}</span>
                </div>
              </template>
            </image-upload>
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <a
              v-if="publishing"
              :class="{ 'is-loading': publishing }"
              disabled
              class="button is-primary"
              >{{ $t("pages.article.create.publishBtn") }}</a
            >
            <a
              v-else
              :class="{ 'is-loading': publishing }"
              class="button is-primary"
              @click="submitCreate"
              >{{ $t("pages.article.create.publishBtn") }}</a
            >
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const { t } = useI18n();
const publishing = ref(false);
const postForm = ref({
  title: "",
  tags: [],
  content: "",
});

const userStore = useUserStore();
const configStore = useConfigStore();
const isNeedEmailVerify = computed(() => {
  return (
    configStore.config.createArticleEmailVerified &&
    !userStore.user.emailVerified
  );
});

if (!configStore.config.modules.article) {
  showError(t("pages.article.create.featureClosed"));
}

useHead({
  title: useSiteTitle(t("pages.article.create.title")),
});

definePageMeta({
  middleware: "auth",
});

if (!configStore.isEnabledArticle) {
  throw createError({
    statusCode: 403,
    message: t("pages.article.create.featureForbidden"),
  });
}

async function submitCreate() {
  if (publishing.value) {
    return;
  }
  publishing.value = true;
  try {
    const article = await useHttpPost(
      "/api/article/create",
      useJsonToForm({
        title: postForm.value.title,
        content: postForm.value.content,
        tags: postForm.value.tags ? postForm.value.tags.join(",") : "",
        cover:
          postForm.value.cover && postForm.value.cover.length
            ? JSON.stringify(postForm.value.cover[0])
            : null,
      })
    );
    useMsg({
      message: t("pages.article.create.success"),
      onClose() {
        useLinkTo(`/article/${article.id}`);
      },
    });
  } catch (e) {
    useMsgError(e.message || t("pages.article.create.error"));
    publishing.value = false;
  }
}
</script>

<style lang="scss" scoped>
.cover-add-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  i {
    font-size: 18px;
    color: var(--text-color3);
  }

  span {
    font-size: 14px;
    color: var(--text-color3);
    font-weight: 500;
  }
}
</style>

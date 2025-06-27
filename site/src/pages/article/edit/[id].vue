<template>
  <section class="main">
    <div class="container" v-if="postForm">
      <div class="publish-form">
        <div class="form-title">
          <div class="form-title-name">
            {{ $t("pages.article.edit.title") }}
          </div>
        </div>

        <div class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input article-title"
              type="text"
              :placeholder="$t('pages.article.edit.titlePlaceholder')"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <markdown-editor
              v-model="postForm.content"
              :placeholder="$t('pages.article.edit.contentPlaceholder')"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <a
              v-if="publishing"
              :class="{ 'is-loading': publishing }"
              disabled
              class="button is-primary"
              >{{ $t("pages.article.edit.submitBtn") }}</a
            >
            <a
              v-else
              :class="{ 'is-loading': publishing }"
              class="button is-primary"
              @click="submitCreate"
              >{{ $t("pages.article.edit.submitBtn") }}</a
            >
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const { t } = useI18n();
definePageMeta({
  middleware: "auth",
});
useHead({
  title: useSiteTitle(t("pages.article.edit.title")),
});

const route = useRoute();

const { data: postForm } = await useMyFetch(
  `/api/article/edit/${route.params.id}`
);
const publishing = ref(false);

async function submitCreate() {
  if (publishing.value) {
    return;
  }
  publishing.value = true;

  try {
    useHttpPost(
      `/api/article/edit/${postForm.value.id}`,
      useJsonToForm({
        title: postForm.value.title,
        content: postForm.value.content,
        cover: postForm.value.cover,
        tags: postForm.value.tags ? postForm.value.tags.join(",") : "",
      })
    );
    useMsg({
      message: $t("pages.article.edit.editSuccess"),
      onClose() {
        useLinkTo(`/article/${postForm.value.id}`);
      },
    });
  } catch (e) {
    publishing.value = false;
    useMsgError($t("pages.article.edit.editFailed") + (e.message || e));
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
    font-size: 24px;
    color: #1878f3;
  }

  span {
    font-size: 14px;
    color: var(--text-color3);
  }
}
</style>

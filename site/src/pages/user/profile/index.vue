<template>
  <div v-if="user" class="widget no-margin">
    <div class="widget-header">
      <div>
        <i class="iconfont icon-username" />
        <span>{{ $t("user.profile.title") }}</span>
      </div>
      <nuxt-link :to="'/user/' + user.id">
        <i class="iconfont icon-return" />
        <span>{{ $t("user.profile.backToProfile") }}</span>
      </nuxt-link>
    </div>
    <div class="widget-content">
      <!-- 头像 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">{{ $t("user.profile.avatar") }}</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <avatar-edit :value="user.avatar" />
            </div>
          </div>
        </div>
      </div>

      <!-- 昵称 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">{{ $t("user.profile.nickname") }}</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <input
                v-model="form.nickname"
                class="input"
                type="text"
                autocomplete="off"
                :placeholder="$t('user.profile.nicknamePlaceholder')"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 个性签名 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">{{ $t("user.profile.signature") }}</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <textarea
                v-model="form.description"
                class="textarea"
                rows="2"
                :placeholder="$t('user.profile.signaturePlaceholder')"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 个人主页 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">{{ $t("user.profile.homepage") }}</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <input
                v-model="form.homePage"
                class="input"
                type="text"
                autocomplete="off"
                :placeholder="$t('user.profile.homepagePlaceholder')"
              />
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label is-normal" />
        <div class="field-body">
          <div class="field">
            <div class="control">
              <a class="button is-primary" @click="submitForm">{{
                $t("user.profile.save")
              }}</a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: ["auth"],
  layout: "profile",
});

const { t } = useI18n();

useHead({
  title: useSiteTitle(t("user.profile.title")),
});

const userStore = useUserStore();
const user = computed(() => {
  return userStore.user;
});

const form = ref({
  nickname: "",
  avatar: "",
  homePage: "",
  description: "",
});

if (user.value != null) {
  form.value.nickname = user.value.nickname;
  form.value.avatar = user.value.avatar;
  form.value.homePage = user.value.homePage;
  form.value.description = user.value.description;
}

async function submitForm() {
  try {
    await useHttpPost(
      `/api/user/edit/${user.value.id}`,
      useJsonToForm(form.value)
    );
    useMsgSuccess(t("user.profile.editSuccess"));
  } catch (e) {
    console.error(e);
    useMsgError(t("user.profile.editFailed") + "：" + (e.message || e));
  }
}
</script>
<style lang="scss" scoped>
.widget-header {
  padding: 18px 0;

  & > div {
    display: flex;
    align-items: center;
    span {
      margin-left: 6px;
    }
  }

  & > a {
    display: flex;
    align-items: center;
    font-weight: 500;
    font-size: 12px;
    span {
      margin-left: 6px;
    }
  }
}

.field {
  margin-bottom: 10px;

  input,
  textarea {
    &:focus-visible {
      outline-width: 0;
    }
  }
}
</style>

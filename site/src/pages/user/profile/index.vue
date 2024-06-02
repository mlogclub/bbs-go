<template>
  <div v-if="user" class="widget no-margin">
    <div class="widget-header">
      <div>
        <i class="iconfont icon-setting" />
        <span>个人资料</span>
      </div>
      <nuxt-link :to="'/user/' + user.id">
        <i class="iconfont icon-return" />
        <span>返回个人主页</span>
      </nuxt-link>
    </div>
    <div class="widget-content">
      <!-- 头像 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">头像</label>
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
          <label class="label">昵称</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <input
                v-model="form.nickname"
                class="input"
                type="text"
                autocomplete="off"
                placeholder="请输入昵称"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 个性签名 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">个性签名</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <textarea
                v-model="form.description"
                class="textarea"
                rows="2"
                placeholder="一句话介绍你自己"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- 个人主页 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">个人主页</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <input
                v-model="form.homePage"
                class="input"
                type="text"
                autocomplete="off"
                placeholder="请输入个人主页"
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
              <a class="button is-success" @click="submitForm">保存</a>
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

useHead({
  title: useSiteTitle("个人资料"),
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
    await useHttpPostForm(`/api/user/edit/${user.value.id}`, {
      body: form.value,
    });
    await reload();
    useMsgSuccess("资料修改成功");
  } catch (e) {
    console.error(e);
    useMsgError("资料修改失败：" + (e.message || e));
  }
}
async function reload() {
  user.value = await useHttpGet("/api/user/current");
  form.value = { ...user.value };
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

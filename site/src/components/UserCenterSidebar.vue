<template>
  <div class="left-container">
    <my-counts :user="localUser" />
    <my-profile :user="localUser" />
    <fans-widget :user="localUser" />
    <follow-widget :user="localUser" />

    <div v-if="isAdmin" class="widget">
      <div class="widget-header">操作</div>
      <div class="widget-content">
        <ul class="operations">
          <li v-if="localUser.forbidden">
            <i class="iconfont icon-forbidden" />
            <a @click="removeForbidden">&nbsp;取消禁言</a>
          </li>
          <template v-else>
            <li>
              <i class="iconfont icon-forbidden" />
              <a @click="forbidden(7)">&nbsp;禁言7天</a>
            </li>
            <li>
              <i v-if="isSiteOwner" class="iconfont icon-forbidden" />
              <a @click="forbidden(-1)">&nbsp;永久禁言</a>
            </li>
          </template>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ElMessageBox } from "element-plus";
const userStore = useUserStore();
const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
});
const localUser = ref(props.user);

const isSiteOwner = computed(() => {
  return userIsOwner(userStore.user);
});

const isAdmin = computed(() => {
  return userIsOwner(userStore.user) || userIsAdmin(userStore.user);
});

function forbidden(days) {
  const msg = days > 0 ? "是否禁言该用户？" : "是否永久禁言该用户？";
  ElMessageBox.confirm(msg, "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(() => {
      doForbidden(days);
    })
    .catch(() => {});
}

async function doForbidden(days) {
  try {
    await useHttpPostForm("/api/user/forbidden", {
      body: {
        userId: localUser.value.id,
        days,
      },
    });
    localUser.value.forbidden = true;
    useMsgSuccess("禁言成功");
  } catch (e) {
    useMsgError("禁言失败");
  }
}

async function removeForbidden() {
  try {
    await useHttpPostForm("/api/user/forbidden", {
      body: {
        userId: localUser.value.id,
        days: 0,
      },
    });
    localUser.value.forbidden = false;
    useMsgSuccess("取消禁言成功");
  } catch (e) {
    useMsgError("取消禁言失败");
  }
}
</script>

<style lang="scss" scoped>
.img-avatar {
  margin-top: 5px;
  border: 1px dotted var(--border-color);
  border-radius: 5%;
}

.operations {
  list-style: none;

  li {
    padding-left: 3px;

    font-size: 13px;
    &:hover {
      cursor: pointer;
      background-color: #fcf8e3;
      color: #8a6d3b;
      font-weight: bold;
    }

    a {
      color: var(--text-link-color);
    }
  }
}
</style>

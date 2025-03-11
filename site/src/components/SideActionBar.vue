<template>
  <div class="action-list">
    <!-- <div class="action btn-share" @click="handleActionClick('share')">
      <i class="i-share" />
    </div> -->
    <div
      class="action btn-like"
      :class="{ active: _liked }"
      @click="handleActionClick('like')"
    >
      <span v-if="_likeCount > 0" class="act-num">{{ _likeCount }}</span>
      <i class="i-like" :class="{ active: _liked }" />
    </div>
    <div class="action btn-comment" @click="handleActionClick('comment')">
      <span v-if="commentCount > 0" class="act-num">{{ commentCount }}</span>
      <i class="i-comment" />
    </div>
    <div
      class="action btn-favorite"
      :class="{ active: _favorited }"
      @click="handleActionClick('favorite')"
    >
      <i class="i-favorite" :class="{ active: _favorited }" />
    </div>
    <div
      v-show="showTopAction"
      class="action"
      @click="handleActionClick('top')"
    >
      <i class="i-top" />
    </div>
  </div>
</template>

<script setup>
import throttle from "lodash.throttle";
const showTopAction = ref(false);
const emits = defineEmits(["handleClick"]);

const props = defineProps({
  parentEleId: {
    type: String,
    default: "",
  },
  entityType: {
    type: String,
    required: true,
  },
  entityId: {
    type: Number,
    required: true,
  },
  liked: {
    type: Boolean,
    default: false,
  },
  likeCount: {
    type: Number,
    default: 0,
  },
  commentCount: {
    type: Number,
    default: 0,
  },
  favorited: {
    type: Boolean,
    default: false,
  },
});

const _liked = ref(props.liked);
const _likeCount = ref(props.likeCount);
const _favorited = ref(props.favorited);

onMounted(() => {
  (document.getElementById(props.parentEleId) || window)?.addEventListener(
    "scroll",
    throttle(onPageScroll, 200)
  );
});

/**
 * 处理页面滚动事件
 */
const onPageScroll = () => {
  const element = document.getElementById(props.parentEleId);
  const winHeight = element?.clientHeight || window.innerHeight;
  const scrollTopVal = element
    ? element.scrollTop
    : window.pageYOffset ||
      document.documentElement.scrollTop ||
      document.body.scrollTop;
  showTopAction.value = scrollTopVal > winHeight / 2;
};

const handleActionClick = (action) => {
  if (action === "top") {
    useScrollToTop(500, props.parentEleId);
  } else if (action === "comment") {
    window.scrollTo(0, document.getElementById("JComment").offsetTop);
  } else if (action === "like") {
    handleLike();
  } else if (action === "favorite") {
    handleFavorite();
  } else {
    emits("handleClick", action);
  }
};

const handleLike = async () => {
  if (_liked.value) {
    try {
      await useHttpPostForm("/api/like/unlike", {
        body: {
          entityType: props.entityType,
          entityId: props.entityId,
        },
      });
    } catch (e) {
      console.log(e);
    }
    _liked.value = false;
    _likeCount.value = _likeCount.value > 0 ? _likeCount.value - 1 : 0;

    useMsgSuccess("已取消点赞");
  } else {
    try {
      await useHttpPostForm("/api/like/like", {
        body: {
          entityType: props.entityType,
          entityId: props.entityId,
        },
      });
    } catch (e) {
      console.log(e);
    }
    _liked.value = true;
    _likeCount.value++;

    useMsgSuccess("点赞成功");
  }
};

const handleFavorite = async () => {
  if (_favorited.value) {
    try {
      await useHttpPostForm("/api/favorite/delete", {
        body: {
          entityType: props.entityType,
          entityId: props.entityId,
        },
      });
    } catch (e) {
      console.log(e);
    }
    _favorited.value = false;

    useMsgSuccess("已取消收藏");
  } else {
    try {
      await useHttpPostForm("/api/favorite/add", {
        body: {
          entityType: props.entityType,
          entityId: props.entityId,
        },
      });
    } catch (e) {
      console.log(e);
    }
    _favorited.value = true;

    useMsgSuccess("收藏成功");
  }
};
</script>

<style lang="scss" scoped>
.action-list {
  .action {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 48px;
    height: 48px;
    border-radius: 50%;
    margin: 0 0 16px;
    text-align: center;
    background: var(--bg-color);
    box-shadow: 0 4px 4px rgba(163, 172, 180, 0.13);
    cursor: pointer;
    transition: all 0.2s linear;
    user-select: none;

    &:hover {
      background: var(--bg-color4);
    }

    &.active {
      background: linear-gradient(336.45deg, #ff5543 42.98%, #ff7827 86.07%);
      box-shadow: 0 4px 4px rgba(255, 92, 63, 0.15);
    }

    i {
      width: 24px;
      height: 24px;
      background-size: 100% 100%;
      background-color: #fff;
      &.i-share {
        background: url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHBhdGggZD0iTTkuNjUzIDIuNUg0Ljc1YTMgMyAwIDAwLTMgM3Y5YTMgMyAwIDAwMyAzaDEwLjM3YTMgMyAwIDAwMy0zdi0zTTYuODMgMTMuNWMwLTUgNC41MTctOC41IDExLjg1NS04LjUiIHN0cm9rZT0iIzg0ODk5NiIgc3Ryb2tlLXdpZHRoPSIxLjg4IiBzdHJva2UtbGluZWNhcD0icm91bmQiLz48cGF0aCBkPSJNMTYuNDI3IDIuNWwyLjQgMi4xMjZhLjUuNSAwIDAxMCAuNzQ4bC0yLjQgMi4xMjYiIHN0cm9rZT0iIzg0ODk5NiIgc3Ryb2tlLXdpZHRoPSIxLjg4IiBzdHJva2UtbGluZWNhcD0icm91bmQiLz48L3N2Zz4=)
          no-repeat 50%;
      }

      &.i-like {
        background: url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHBhdGggZD0iTTYuMDYzIDIuNUE0LjgxIDQuODEgMCAwMDEuMjUgNy4zMDdjMCA0LjQyNyA0LjgyNSA4LjQ4NCA3Ljk4IDkuODkxLjQ5Mi4yMiAxLjA0OC4yMiAxLjU0IDAgMy4xNTUtMS40MDcgNy45OC01LjQ2NCA3Ljk4LTkuODkxQTQuODEgNC44MSAwIDAwMTMuOTM3IDIuNWMtLjk5IDAtMS45MS4yOTktMi42NzUuODEtLjcxOC40ODItMS44MDYuNDgyLTIuNTI0IDBhNC43OTQgNC43OTQgMCAwMC0yLjY3NS0uODF6IiBzdHJva2U9IiM4NDg5OTYiIHN0cm9rZS13aWR0aD0iMS44NzUiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCIvPjwvc3ZnPg==)
          no-repeat 50%;
        &.active {
          background: url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PHBhdGggZD0iTTYuMDYzIDIuNUE0LjgxIDQuODEgMCAwMDEuMjUgNy4zMDdjMCA0LjQyNyA0LjgyNSA4LjQ4NCA3Ljk4IDkuODkxLjQ5Mi4yMiAxLjA0OC4yMiAxLjU0IDAgMy4xNTUtMS40MDcgNy45OC01LjQ2NCA3Ljk4LTkuODkxQTQuODEgNC44MSAwIDAwMTMuOTM3IDIuNWMtLjk5IDAtMS45MS4yOTktMi42NzUuODEtLjcxOC40ODItMS44MDYuNDgyLTIuNTI0IDBhNC43OTQgNC43OTQgMCAwMC0yLjY3NS0uODF6IiBmaWxsPSIjZmZmIiBzdHJva2U9IiNmZmYiIHN0cm9rZS13aWR0aD0iMS44NzUiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCIvPjwvc3ZnPg==)
            no-repeat 50%;
        }
      }

      &.i-comment {
        background: url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAiIGhlaWdodD0iMjAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGcgc3Ryb2tlPSIjODQ4OTk2IiBzdHJva2Utd2lkdGg9IjEuODc1IiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiPjxwYXRoIGQ9Ik0xOC44NTEgNWEyLjUgMi41IDAgMDAtMi41LTIuNUgzLjc1QTIuNSAyLjUgMCAwMDEuMjUgNXY4LjQ0NEEyLjI1NyAyLjI1NyAwIDAwMy41MDcgMTUuN2guNjY3Yy42MTMgMCAxLjE1My40MDcgMS4zMjEuOTk3djBhMS4zNzQgMS4zNzQgMCAwMDIuMjI2LjY1NmwxLjIxMi0xLjA2YTIuMzk1IDIuMzk1IDAgMDExLjU3OC0uNTkzaDUuODRhMi41IDIuNSAwIDAwMi41LTIuNVY1ek02LjkgOS4xaDYuMjUiLz48L2c+PC9zdmc+)
          no-repeat 50%;
      }

      &.i-favorite {
        background: url(data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGcgY2xpcC1wYXRoPSJ1cmwoI2NsaXAwXzU2OTNfMjU2NikiPjxwYXRoIGQ9Ik03LjA4NCAxLjU4OGMuMzUtLjc5OCAxLjQ4Mi0uNzk4IDEuODMyIDBsMS4yMjkgMi44MDJhMSAxIDAgMDAuNzkyLjU5bDIuOTk3LjM3NGExIDEgMCAwMS41NzMgMS43MWwtMi4yODIgMi4yMThhMSAxIDAgMDAtLjI4NS45MDJsLjU5NSAzLjE1Yy4xNi44NDUtLjc1IDEuNDg0LTEuNDkgMS4wNDhMOC41MDcgMTIuODlhMSAxIDAgMDAtMS4wMTQgMGwtMi41MzggMS40OTNjLS43NC40MzYtMS42NS0uMjAzLTEuNDktMS4wNDdsLjU5NS0zLjE1YTEgMSAwIDAwLS4yODUtLjkwM0wxLjQ5MyA3LjA2M2ExIDEgMCAwMS41NzMtMS43MWwyLjk5Ny0uMzczYTEgMSAwIDAwLjc5Mi0uNTlsMS4yMy0yLjgwMnoiIHN0cm9rZT0iIzczNzc4MiIgc3Ryb2tlLXdpZHRoPSIxLjg3NSIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIi8+PC9nPjxkZWZzPjxjbGlwUGF0aCBpZD0iY2xpcDBfNTY5M18yNTY2Ij48cGF0aCBmaWxsPSIjZmZmIiBkPSJNMCAwaDE2djE2SDB6Ii8+PC9jbGlwUGF0aD48L2RlZnM+PC9zdmc+)
          no-repeat 50%;

        &.active {
          background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAAAUCAYAAACNiR0NAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAADpSURBVHgBrZTREYIwDIYDx7uM0A1khG7iCI4gI+gE6ASOwChlA9ggJhJODLTawneX6/Wa/A1JCsAPENGQPcUMbIVEHH5oYQskYHGJhVQo+L4imJal1M6H8cUVSqSkha0iOwXuu5Dvg9aBrMuybNCZtGQ9ptOLhvHVKZUmpyQPsB/vDCv8nrVUHE7NwrGODtNxqDu/QdShb4xENKbbvRbL5xuapw4i0TG5ypAHuoT/KYMZArc9HhsStBBPNd8U6vC4EsDv9CbrGZZfYcAH1aNWXWxw/GFM5zwF+qleIYSIspAN+BjxqecXMi9Co2eLeEDpiAAAAABJRU5ErkJggg==)
            no-repeat 50%;
        }
      }

      &.i-top {
        width: 24px;
        height: 24px;
        background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAAACXBIWXMAABYlAAAWJQFJUiTwAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAEcSURBVHgB7ZcxCsJAEEVnPYE3VFutrEQUUVDEzkZb4808Qi5gxk3IgISNTrKzpPkPUmRhkv/YYvcTAQDAkDhKxPFyn5DjK7PL/U/229X8SQlIIlCFJ8oay9MUEuYCLeEFcwlTgT/hBVMJMwFleMFMYkQGtIUvCp45/wRGsnommugd+BV+t15U66fzbcoj9wiMR+9ElIAmvJBKordAl/BCColeAn3CC9YSnQViwguWEp0ELMILVhJqAcvwgoWE+hwoL2TNtZjwJRs/Gz4n+EpK1ALsePz9HhteCEkwuVw7r9+BNy39Lrz8FTm3Ci+IhHyfCz4QAEAFOnEIdOIOoBML6MQ9QScmdOIKdGI16MQ16MQ16MQUCToxAAAMyge7rm/YCAKydgAAAABJRU5ErkJggg==)
          no-repeat 50%;
        background-size: 100% 100%;
      }
    }

    .act-num {
      position: absolute;
      top: -4px;
      left: 24px;
      padding: 1px 5px;
      background: #ff5543;
      border-radius: 11px;
      font-style: normal;
      font-weight: 400;
      font-size: 14px;
      line-height: 14px;
      color: #fff;
      white-space: nowrap;
      border: 1px solid #fff;
    }
  }
}
</style>

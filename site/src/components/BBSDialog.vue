<template>
  <div class="dialog-component">
    <div
      class="dialog-mask"
      :class="{ visible: visible }"
      :style="{ zIndex: zIndex }"
    ></div>
    <transition :name="transition">
      <div
        v-if="visible"
        class="dialog-wrapper"
        :style="{ zIndex: zIndex + 1 }"
      >
        <div
          class="dialog-content"
          :style="{
            width: dialogContentWidth,
            maxWidth: dialogContentMaxWidth,
          }"
        >
          <div class="dialog-header">
            <div class="dialog-title">{{ title }}</div>
            <div class="dialog-close">
              <img src="~/assets/images/close2.png" @click="close" />
            </div>
          </div>
          <div class="dialog-main">
            <slot></slot>
          </div>
          <div
            v-if="footerVisible"
            class="dialog-footer"
            :style="{ justifyContent: btnsCenter ? 'center' : 'flex-end' }"
          >
            <button
              v-if="cancelBtnVisible"
              class="button is-small is-light is-danger"
              @click="cancel"
            >
              {{ cancelBtnText }}
            </button>
            <button
              v-if="okBtnVisible"
              class="button is-small is-light is-success"
              @click="ok"
            >
              {{ okBtnText }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  // 宽度
  width: {
    type: Number,
    default: 520, // 0表示100%
  },
  // 最大宽度
  maxWidth: {
    type: Number,
    default: 0, // 0表示不设置
  },
  zIndex: {
    type: Number,
    default: 1001,
  },
  title: {
    type: String,
    default: "",
  },
  transition: {
    type: String,
    default: "bounce", // bounce, fade
  },
  btnsCenter: {
    type: Boolean,
    default: false,
  },
  footerVisible: {
    type: Boolean,
    default: true,
  },
  cancelBtnVisible: {
    type: Boolean,
    default: true,
  },
  cancelBtnText: {
    type: String,
    default: "取消",
  },
  okBtnVisible: {
    type: Boolean,
    default: true,
  },
  okBtnText: {
    type: String,
    default: "确定",
  },
});

const emit = defineEmits(["update:visible", "close", "ok"]);

const dialogContentWidth = computed(() => {
  if (props.width > 0) {
    return `${props.width}px`;
  }
  // return 'auto'
  return "100%";
});

const dialogContentMaxWidth = computed(() => {
  if (props.maxWidth > 0) {
    return `${props.maxWidth}px`;
  }
  return "";
});

onMounted(() => {
  window.addEventListener("keydown", handleEscKey);
});

onBeforeUnmount(() => {
  window.removeEventListener("keydown", handleEscKey);
});

const handleEscKey = (event) => {
  if (event.key === "Escape" || event.keyCode === 27) {
    close();
  }
};
const show = () => {
  emit("update:visible", true);
};
const close = () => {
  emit("update:visible", false);
};
const ok = () => {
  emit("ok");
};
const cancel = () => {
  emit("cancel");
  close();
};

defineExpose({
  show,
  close,
});
</script>

<style lang="scss" scoped>
.dialog-mask {
  // transition: all 2s;
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  overflow: auto;
  background: #000000;
  opacity: 0.45;
  display: none;

  &.visible {
    display: block;
  }
}
.dialog-wrapper {
  // transition: all 2s;
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;

  display: flex;
  align-items: center;
  justify-content: center;

  margin-bottom: 50px;

  .dialog-content {
    position: relative;
    margin: 0 auto;
    // margin-top: 15vh;
    background: var(--bg-color);
    border-radius: 8px;
    padding: 12px 18px;
    .dialog-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      .dialog-title {
        font-size: 16px;
        font-weight: 500;
      }
      .dialog-close {
        cursor: pointer;
        padding: 2px;
        display: flex;
        border-radius: 50%;
        background-color: var(--bg-color3);
        img {
          width: 20px;
          height: 20px;
        }
      }
    }
    .dialog-main {
      padding: 12px 0;
    }
    .dialog-footer {
      display: flex;
      align-items: center;
      // justify-content: flex-end;
      column-gap: 24px;

      .chaitin-btn {
        width: 78px;
      }
    }
  }
}
</style>

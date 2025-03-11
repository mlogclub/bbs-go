<template>
  <div class="dialog-component">
    <transition name="fade">
      <div v-if="visible" class="dialog-mask" :style="{ zIndex: zIndex }"></div>
    </transition>
    <transition :name="transition">
      <div
        v-if="visible"
        class="dialog-wrapper"
        :style="{ zIndex: zIndex + 1 }"
      >
        <div class="dialog-content" style="width: max-content">
          <Rotate
            v-if="visible"
            :config="{}"
            :data="{
              image: captcha.imageBase64,
              thumb: captcha.thumbBase64,
            }"
            :events="{
              refresh: captchaRefresh,
              close: captchaClose,
              confirm: captchaConfirm,
            }"
          />
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import "go-captcha-vue/dist/style.css";
import { Rotate } from "go-captcha-vue";
const emits = defineEmits(["confirm"]);

const props = defineProps({
  transition: {
    type: String,
    default: "bounce", // bounce, fade
  },
  zIndex: {
    type: Number,
    default: 1001,
  },
});

const visible = ref(false);
const captcha = ref(null);

const show = async () => {
  try {
    captcha.value = await useHttpGet("/api/captcha/request_angle");
    visible.value = true;
  } catch (error) {
    console.error(error);
  }
};

const captchaRefresh = async () => {
  try {
    captcha.value = await useHttpGet("/api/captcha/request_angle");
  } catch (e) {
    useCatchError(e);
  }
};

const captchaClose = () => {
  visible.value = false;
};

const captchaConfirm = (angle, reset) => {
  emits(
    "confirm",
    // 参数1：验证码信息
    {
      captchaId: captcha.value.id,
      captchaCode: angle,
    },
    // 参数2：回调
    (success) => {
      captchaClose();
      //   if (success) {
      //     captchaClose();
      //   } else {
      //     captchaRefresh();
      //   }
    }
  );
};

defineExpose({
  show,
});
</script>

<style lang="scss" scoped>
.dialog-mask {
  transition: all 0s;
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  overflow: auto;
  background: #00000066;
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

  .dialog-content {
    position: relative;
    margin: 0 auto;
    // margin-top: 15vh;
    background: #ffffff;
    border-radius: 8px;

    padding: 0;
  }
}
</style>

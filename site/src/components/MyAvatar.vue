<template>
  <nuxt-link
    v-if="user"
    class="avatar-a"
    target="_blank"
    :to="`/user/${user.id}`"
    :class="[sizeClass]"
    :style="extraStyle"
  >
    <img
      v-if="hasAvatarUrl && !loadError"
      :src="avatarUrl"
      class="avatar"
      :class="[sizeClass, roundClass, borderClass]"
      :alt="user.nickname"
      @error="error"
    />
    <span
      v-else-if="styleText"
      class="avatar"
      :class="[sizeClass, roundClass, borderClass]"
      :style="styleText"
    >
      {{ usernameAt }}
    </span>
  </nuxt-link>
</template>

<script>
export default {
  name: "Avatar",
  props: {
    user: {
      type: [Object, String],
      default: () => {},
    },
    size: {
      type: [Number, String],
      default: 50,
    },
    round: {
      type: Boolean,
      default: true,
    },
    hasBorder: {
      type: Boolean,
      default: false,
    },
    extraStyle: {
      type: [Object],
      default: () => {},
    },
  },
  data() {
    return {
      loadError: false,
      sizes: {
        150: "font-size: 60px;line-height: 150px;border-radius: 2px",
        100: "font-size: 40px;line-height: 100px;border-radius: 2px",
        80: "font-size: 32px;line-height: 80px;border-radius: 2px",
        70: "font-size: 28px;line-height: 70px;border-radius: 2px",
        60: "font-size: 26px;line-height: 60px;border-radius: 2px",
        50: "font-size: 24px;line-height: 50px;border-radius: 2px",
        45: "font-size: 22px;line-height: 45px;border-radius: 2px",
        40: "font-size: 20px;line-height: 40px;border-radius: 2px",
        35: "font-size: 18px;line-height: 30px;border-radius: 2px",
        30: "font-size: 18px;line-height: 30px;border-radius: 2px",
        24: "font-size: 12px;line-height: 24px;border-radius: 2px",
        20: "font-size: 10px;line-height: 20px;border-radius: 2px",
      },
    };
  },
  computed: {
    hasAvatarUrl() {
      return this.avatarUrl;
    },
    avatarUrl() {
      return this.user.smallAvatar || this.user.avatar;
    },
    usernameAt() {
      let c = this.user.nickname
        ? this.user.nickname.charAt(0).toUpperCase()
        : "";
      if (!c) {
        c = this.user.username
          ? this.user.username.charAt(0).toUpperCase()
          : "";
      }
      return c;
    },
    styleText() {
      return `background-color: #${useStringToColor(this.usernameAt)};${
        this.sizes[this.size]
      }`;
    },
    sizeClass() {
      return [`avatar-size-${this.size}`];
    },
    roundClass() {
      return this.round ? "round" : "";
    },
    borderClass() {
      return this.hasBorder ? "has-border" : "";
    },
  },
  methods: {
    error() {
      this.loadError = true;
    },
  },
};
</script>

<style lang="scss" scoped>
.avatar-size-150 {
  width: 150px;
  height: 150px;
  min-width: 150px;
  min-height: 150px;
  border-radius: 2px;
}
.avatar-size-100 {
  width: 100px;
  height: 100px;
  min-width: 100px;
  min-height: 100px;
  border-radius: 2px;
}
.avatar-size-80 {
  width: 80px;
  height: 80px;
  min-width: 80px;
  min-height: 80px;
  border-radius: 2px;
}

.avatar-size-70 {
  width: 70px;
  height: 70px;
  min-width: 70px;
  min-height: 70px;
  border-radius: 2px;
}

.avatar-size-60 {
  width: 60px;
  height: 60px;
  min-width: 60px;
  min-height: 60px;
  border-radius: 2px;
}
.avatar-size-50 {
  width: 50px;
  height: 50px;
  min-width: 50px;
  min-height: 50px;
  border-radius: 2px;
}
.avatar-size-45 {
  width: 45px;
  height: 45px;
  min-width: 45px;
  min-height: 45px;
  border-radius: 2px;
}
.avatar-size-40 {
  width: 40px;
  height: 40px;
  min-width: 40px;
  min-height: 40px;
  border-radius: 2px;
}
.avatar-size-35 {
  width: 35px;
  height: 35px;
  min-width: 35px;
  min-height: 35px;
  border-radius: 2px;
}
.avatar-size-30 {
  width: 30px;
  height: 30px;
  min-width: 30px;
  min-height: 30px;
  border-radius: 2px;
}
.avatar-size-24 {
  width: 24px;
  height: 24px;
  min-width: 24px;
  min-height: 24px;
  border-radius: 2px;
}
.avatar-size-20 {
  width: 20px;
  height: 20px;
  min-width: 20px;
  min-height: 20px;
  border-radius: 2px;
}

.round {
  border-radius: 50% !important;
}

.avatar-a {
  display: block;
  position: relative;
}

.avatar {
  &.has-border {
    border: 1px solid var(--border-color);
  }
}

img.avatar {
  object-fit: cover;
  transition: all 0.5s ease-out 0.1s;
  background-color: var(--bg-color);

  &:hover {
    transform: matrix(1.04, 0, 0, 1.04, 0, 0);
    backface-visibility: hidden;
  }
}
span.avatar {
  color: var(--text-color5);
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
}
</style>

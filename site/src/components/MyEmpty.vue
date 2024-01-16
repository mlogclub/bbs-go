<script>
export default {
  props: {
    showLogo: {
      type: Boolean,
      default: true,
    },
    description: {
      type: String,
      default: "暂无数据",
    },
  },
};
</script>

<template>
  <div class="empty-data">
    <svg v-if="showLogo" viewBox="0 0 66 68">
      <g fill="none" fill-rule="evenodd" transform="translate(4 3)">
        <g fill="#F7F7F7">
          <path
            d="M9 10h23.751v3.221H9zM9 16.494h41.083v4.026H9zM9 26.104h23.751v3.221H9zM9 42.208h23.751v3.221H9zM9 33.351h41.083v4.026H9zM9 49.455h41.083v4.026H9z"
          />
        </g>
        <rect
          width="56"
          height="60"
          x="1.139"
          y="1.338"
          stroke="#EBEBEB"
          stroke-width="2"
          rx="6"
        />
      </g>
    </svg>
    <span>
      {{ description }}
    </span>
  </div>
</template>

<style lang="scss" scoped>
.empty-data {
  background-color: var(--bg-color);
  padding: 30px 0;
  border-radius: 0 0 0.2rem 0.2rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  font-weight: 400;
  margin-bottom: 0.67rem;

  svg {
    height: 100px;
  }
  span {
    margin-top: 20px;
    font-size: 15px;
    font-weight: 400;
    color: var(--text-color4);
  }
}
</style>

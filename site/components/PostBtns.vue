<template>
  <div class="aside-actions">
    <a :href="createTopicUrl" class="aside-actions-btn">
      <svg
        viewBox="0 0 1065 1024"
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        class="icon ic"
      >
        <path
          d="M331.329 566.55a7.454 7.454 0 00-1.649 2.903l-53.561 194.985a32.981 32.981 0 0032.057 41.754 31.398 31.398 0 008.443-1.187l194.787-52.77h.726a8.575 8.575 0 006.134-2.507l520.773-517.672a82.519 82.519 0 0024.01-59.366 100.725 100.725 0 00-30.342-70.844l-49.208-49.01a102.242 102.242 0 00-71.371-30.21 83.904 83.904 0 00-59.366 23.746L332.055 564.77c-.528.528 0 1.253-.726 1.78m656.985-384.824l-51.648 51.385-84.036-84.696 51.055-50.725c8.113-8.113 23.68-6.596 32.98 2.375l49.209 49.01a26.385 26.385 0 018.113 18.601 19.13 19.13 0 01-5.475 14.05M423.94 574.664l375.986-373.61 83.97 84.695-375.392 373.017zM355.405 725.85l27.177-98.944 72.558 72.56zm665.89-339.047a35.95 35.95 0 00-35.949 35.62v482.845a46.174 46.174 0 01-46.174 46.173H117.941a46.174 46.174 0 01-46.174-46.173V117.545a46.174 46.174 0 0146.174-46.174h528.886a35.752 35.752 0 100-71.107H112.136A112.136 112.136 0 000 112.4v798.145a112.136 112.136 0 00112.136 112.136h832.511a112.136 112.136 0 00112.532-112.136V422.423a35.95 35.95 0 00-35.95-35.62"
        ></path>
      </svg>
      <span>发布话题</span>
    </a>
    <a href="/article/create" class="aside-actions-btn">
      <svg
        viewBox="0 0 1034 1024"
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        class="icon ic"
      >
        <path
          d="M978.15 494.17c-20.379 0-40.757 15.283-40.757 40.755v412.657H71.323V81.512H483.98c20.378 0 40.756-15.283 40.756-40.756 5.095-20.378-15.283-35.661-35.661-35.661H86.607C40.757 5.095 0 45.85 0 91.7v850.787C0 983.244 40.756 1024 86.607 1024h845.692c45.85 0 86.606-40.756 86.606-86.607V534.925c0-20.378-15.283-40.756-40.756-40.756z"
        ></path>
        <path
          d="M560.398 234.348c0 20.378 15.284 40.756 40.756 40.756H753.99V427.94c0 20.378 15.284 40.757 40.756 40.757s40.756-15.284 40.756-40.757V275.104h152.836c20.378 0 40.757-15.283 40.757-40.756s-15.284-40.756-40.757-40.756H835.502V40.756C835.502 20.378 820.22 0 794.746 0S753.99 15.284 753.99 40.756v152.836H601.154c-25.472 0-40.756 20.378-40.756 40.756z"
        ></path>
      </svg>
      <span>发表文章</span>
    </a>
  </div>
</template>

<script>
export default {
  props: {
    currentNodeId: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    createTopicUrl() {
      let url = '/topic/create'
      if (this.currentNodeId) {
        url += '?nodeId=' + this.currentNodeId
      }
      return url
    },
  },
}
</script>

<style lang="scss" scoped>
.aside-actions {
  background: #fff;
  display: flex;
  text-align: center;
  /*padding: 10px;*/

  .aside-actions-btn {
    display: flex;
    color: #000;
    position: relative;
    margin: 0 auto;
    padding: 0.375rem 0.75rem;
    font-weight: 500;
    font-size: 1rem;
    /*line-height: 2;*/
    border-radius: 0.25rem;
    border: 1px solid transparent;
    background-color: #fff;
    text-align: center;
    vertical-align: middle;
    cursor: pointer;
    flex-direction: column;
    align-items: center;
    white-space: nowrap;
    user-select: none;
    transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out,
      border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;

    span {
      margin-top: 5px;
    }
    &:hover {
      color: #535d72;
      font-weight: 600;
    }
  }
}
</style>

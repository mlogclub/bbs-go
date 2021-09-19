<template>
  <div class="dashboard-container" />
</template>

<script>
import { mapGetters } from "vuex";

export default {
  components: {},
  data() {
    return {};
  },
  computed: {
    ...mapGetters(["roles"]),
  },
};
</script>

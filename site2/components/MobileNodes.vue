<template>
  <div class="mobile-nodes">
    <transition name="fadeDown">
      <div v-show="show" class="nodes">
        <div class="nodes-row first">
          <nuxt-link to="/topics/node/newest">
            <div class="node-item" :class="{ active: currentNodeId === 0 }">
              <span>最新</span>
            </div>
          </nuxt-link>
          <nuxt-link to="/topics/node/recommend">
            <div class="node-item" :class="{ active: currentNodeId === -1 }">
              <span>推荐</span>
            </div>
          </nuxt-link>
        </div>
        <div v-for="(row, index) in rows" :key="index" class="nodes-row">
          <div v-if="row && row.length" class="nodes-row">
            <nuxt-link
              v-for="node in row"
              :key="node.nodeId"
              :to="'/topics/node/' + node.nodeId"
            >
              <div
                class="node-item"
                :class="{ active: currentNodeId === node.nodeId }"
              >
                <span>{{ node.name }}</span>
              </div>
            </nuxt-link>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  data() {
    return {
      nodes: [],
    }
  },
  computed: {
    show() {
      return this.$store.state.env.showMobileNodes
    },
    rows() {
      const rowCount = 3
      const arrTemp = []
      const nodes = this.nodes
      let index = 0
      if (nodes && nodes.length) {
        for (let i = 0; i < nodes.length; i++) {
          index = parseInt(i / rowCount)
          if (arrTemp.length <= index) {
            arrTemp.push([])
          }
          arrTemp[index].push(nodes[i])
        }
      }
      return arrTemp
    },
    currentNodeId() {
      return this.$store.state.env.currentNodeId
    },
  },
  mounted() {
    this.fetch()
  },
  methods: {
    async fetch() {
      this.nodes = await this.$axios.get('/api/topic/nodes')
    },
  },
}
</script>

<style lang="scss" scoped></style>

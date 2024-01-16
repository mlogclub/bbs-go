import { defineStore } from 'pinia'

export const useEnvStore = defineStore('env', {
    state: () => ({
        // 当前所在节点编号，0：最新、-1：recommend
        currentNodeId: 0,
    }),
    getters: {

    },
    actions: {
        setCurrentNodeId(nodeId: number) {
            this.currentNodeId = nodeId
        },
    },
});

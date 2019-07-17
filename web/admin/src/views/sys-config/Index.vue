<template>
  <section>

    <div class="config">

      <el-form label-width="100px">
        <el-form-item v-for="config in configs" :key="config.key" :label="config.name">
          <el-input type="textarea" autosize auto-complete="off" :placeholder="config.description"
                    v-model="config.value"></el-input>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="save">保存</el-button>
        </el-form-item>
      </el-form>

    </div>


  </section>
</template>

<script>
  import HttpClient from '../../apis/HttpClient'

  export default {
    name: "List",
    data() {
      return {
        configs: [],
        loading: false
      }
    },
    mounted() {
      this.load();
    },
    methods: {
      load() {
        let me = this
        HttpClient.get('/api/admin/sys-config/all')
          .then(data => {
            me.configs = data
          })
          .catch(rsp => {
            me.$notify.error({title: '错误', message: rsp.message})
          })
      },
      save() {
        let me = this
        let configParam = {}
        for (let i = 0; i < me.configs.length; i++) {
          let item = me.configs[i]
          configParam[item.key] = item.value
        }
        me.loading = true
        HttpClient.post('/api/admin/sys-config/save', {
          config: JSON.stringify(configParam)
        }).then(() => {
          me.loading = false
          me.$message({message: '提交成功', type: 'success'});
        }).catch(rsp => {
          me.loading = false
          me.$notify.error({title: '错误', message: rsp.message})
        })
      }
    }
  }
</script>

<style scoped lang="scss">
  .config {
    margin-top: 20px;
  }
</style>


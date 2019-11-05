<template>
  <section class="page-container">
    <div class="config">
      <el-form v-model="config" label-width="100px">
        <el-form-item label="网站名称">
          <el-input v-model="config['site.title']" type="text" placeholder="网站名称"></el-input>
        </el-form-item>

        <el-form-item label="网站描述">
          <el-input
            v-model="config['site.description']"
            type="textarea"
            autosize
            placeholder="网站描述"
          ></el-input>
        </el-form-item>

        <el-form-item label="网站关键字">
          <el-select
            v-model="config['site.keywords']"
            style="width:100%"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="网站关键字"
          ></el-select>
        </el-form-item>

        <el-form-item label="推荐标签">
          <el-select
            v-model="config['recommend.tags']"
            style="width:100%"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="推荐标签"
          ></el-select>
        </el-form-item>

        <el-form-item label="论坛导航">
          <el-select
            v-model="config['bbs.nav.tags']"
            style="width:100%"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="论坛导航标签，用于显示在讨论区侧边栏"
          ></el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="save">保存</el-button>
        </el-form-item>
      </el-form>
    </div>
  </section>
</template>

<script>
import HttpClient from "@/apis/HttpClient";

export default {
  name: "List",
  data() {
    return {
      config: {
        "site.title": "",
        "site.description": "",
        "site.keywords": [],
        "recommend.tags": [],
        "bbs.nav.tags": []
      },
      loading: false,
      complateTags: [],
      complateLoading: false
    };
  },
  mounted() {
    this.load();
  },
  methods: {
    load() {
      const me = this;
      HttpClient.get("/api/admin/sys-config/all")
        .then(data => {
          for (let i = 0; i < data.length; i++) {
            const item = data[i];
            if (!me.config.hasOwnProperty(item.key)) {
              continue;
            }
            switch (item.key) {
              case "site.keywords":
              case "recommend.tags":
              case "bbs.nav.tags":
                this.config[item.key] = me.splitBy(item.value, ',');
                break;
              default:
                this.config[item.key] = item.value;
                break;
            }
          }
        })
        .catch(rsp => {
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
    save() {
      console.log(this.config);
      const me = this;
      me.loading = true;
      HttpClient.post("/api/admin/sys-config/save", {
        config: JSON.stringify(this.config)
      })
        .then(() => {
          me.loading = false;
          me.$message({ message: "提交成功", type: "success" });
        })
        .catch(rsp => {
          me.loading = false;
          me.$notify.error({ title: "错误", message: rsp.message });
        });
    },
    splitBy(str, separator){
      if (!str || !separator) {
        return []
      }
      return str.split(separator)
    }
  }
};
</script>

<style scoped lang="scss">
.config {
  padding: 10px 0;
}
</style>

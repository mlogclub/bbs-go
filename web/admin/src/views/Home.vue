<template>
  <el-row class="container">

    <frame-header/>

    <el-col :span="24" class="main">

      <side-menu/>

      <section class="content-container">
        <div class="grid-content bg-purple-light">
          <el-col :span="24" class="breadcrumb-container">
            <strong class="title">{{$route.name}}</strong>
            <el-breadcrumb separator="/" class="breadcrumb-inner">
              <el-breadcrumb-item v-for="item in $route.matched" :key="item.path">
                {{ item.name }}
              </el-breadcrumb-item>
            </el-breadcrumb>
          </el-col>
          <el-col :span="24" class="content-wrapper">
            <transition name="fade" mode="out-in">
              <router-view></router-view>
            </transition>
          </el-col>
        </div>
      </section>
    </el-col>
  </el-row>
</template>

<script>
  import FrameHeader from '../components/FrameHeader';
  import SideMenu from '../components/SideMenu';

  export default {
    components: {
      FrameHeader,
      SideMenu
    },
    data() {
      return {
        form: {
          name: '',
          region: '',
          date1: '',
          date2: '',
          delivery: false,
          type: [],
          resource: '',
          desc: ''
        }
      };
    },
    methods: {},
    computed: {
      collapsed() {
        return this.$store.state.Default.collapsed;
      }
    }
  };

</script>

<style scoped lang="scss">
  @import '../styles/vars.scss';

  .container {
    position: absolute;
    top: 0px;
    bottom: 0px;
    width: 100%;

    .main {
      display: flex;
      position: absolute;
      top: 60px;
      bottom: 0px;
      overflow: hidden;

      .content-container {
        flex: 1;
        overflow-y: scroll;
        padding: 20px;

        .breadcrumb-container {
          .title {
            width: 200px;
            float: left;
            color: #475669;
          }

          .breadcrumb-inner {
            float: right;
          }
        }

        .content-wrapper {
          background-color: #fff;
          box-sizing: border-box;
        }
      }
    }
  }
</style>

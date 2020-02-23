<template>
  <div class="post-twitter-box">
    <ul class="ugc-tab-list">
      <li class="ugc-tab-item current">发表推文</li>
    </ul>
    <div class="imgBox upload-box">
      <textarea placeholder="有什么新鲜事想告诉大家" class="title-input" />
      <p class="words-number">0/2000字</p>
      <div class="upload-footer">
        <div class="bui-left">
          <span
            ga_event="user_ugc_img_open"
            class="show-image-uploader show-uploader"
            ><img
              src="http://sf1-ttcdn-tos.pstatp.com/obj/ttfe/toutiao.com/image-icon.png"
              class="image-icon"
            />
            <span>图片</span></span
          >
          <span class="show-emoji show-uploader"
            ><img
              src="http://sf1-ttcdn-tos.pstatp.com/obj/ttfe/toutiao.com/emoji-icon.png"
              class="image-icon"
            />
            <span>表情</span></span
          >
        </div>
        <div class="bui-right">
          <span class="msg-tip"></span> <a class="upload-publish">发布</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {}
</script>

<style lang="scss" scoped>
.post-twitter-box {
  position: relative;
  border: 1px solid #e8e8e8;
  width: 100%;
  margin: 10px 0;

  .ugc-tab-list {
    height: 44px;
    border-bottom: 1px solid #e8e8e8;
    display: block;
    zoom: 1;

    .ugc-tab-item {
      margin-left: 19px;
      font-size: 15px;
      color: #222;
      line-height: 42px;
      border-bottom: 2px solid transparent;
      cursor: pointer;
      float: left;

      &.current {
        border-bottom-color: #ed4040;
        color: #222;
      }
    }
  }

  .upload-box {
    padding: 0;
    margin: 0;
    box-sizing: border-box;
    position: relative;

    .title-input {
      width: 100%;
      height: 100px;
      display: block;
      font-size: 14px;
      line-height: 1.4;
      padding: 13px 19px;
      border: 0;
      outline: 0;
      resize: none;
      background-color: #f4f5f6;
      transition: all 0.5s;
      animation-duration: 0.8s;
      animation-fill-mode: both;
    }

    .words-number {
      position: absolute;
      z-index: 3;
      bottom: 50px;
      right: 10px;
      display: inline-block;
      background-color: rgba(0, 0, 0, 0.5);
      border-radius: 50px;
      padding: 0 8px;
      color: #fff;
      font-size: 13px;
    }

    .upload-footer {
      border-top: 1px solid #e8e8e8;
      height: 40px;
      display: block;
      zoom: 1;

      .bui-left {
        float: left;

        .show-uploader {
          color: #222;
          font-size: 14px;
          line-height: 39px;
          display: inline-block;
          vertical-align: middle;
          margin: 0 0 0 20px;
          cursor: pointer;

          .image-icon {
            width: 20px;
            height: 20px;
            position: relative;
            top: 6px;
          }

          span {
            display: inline-block;
            vertical-align: middle;
          }
        }
      }

      .bui-right {
        float: right;
        text-align: right;

        .msg-tip {
          color: #ed4040;
          font-size: 12px;
          margin-right: 10px;
        }

        .upload-publish {
          display: inline-block;
          width: 120px;
          line-height: 40px;
          text-align: center;
          font-size: 14px;
          background-color: #ed4040;
          color: #fff;
          opacity: 0.6;

          &.active {
            opacity: 1;
          }
        }
      }
    }
  }
}
</style>

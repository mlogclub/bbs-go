<template>
  <ul class="tweets-list">
    <li v-for="item in tweets" :key="item.tweetsId" class="tweets">
      <div class="pin">
        <div class="pin-header-row">
          <div class="account-group">
            <div class="user-popover-box">
              <a
                href="/user/5a2b9a3a6fb9a045055e0daf"
                target="_blank"
                rel=""
                class="user-link"
              >
                <div
                  data-src="https://user-gold-cdn.xitu.io/2017/12/9/1603a537e43e66ee?imageView2/1/w/100/h/100/q/85/format/webp/interlace/1"
                  class="lazy avatar avatar loaded"
                  style='background-image: url("https://user-gold-cdn.xitu.io/2017/12/9/1603a537e43e66ee?imageView2/1/w/100/h/100/q/85/format/webp/interlace/1");'
                ></div>
              </a>
            </div>
            <div class="pin-header-content">
              <div class="user-popover-box">
                <a
                  href="/user/5a2b9a3a6fb9a045055e0daf"
                  target="_blank"
                  rel=""
                  class="username"
                  >林_小猫</a
                >
              </div>
              <div class="meta-box">
                <div class="position ellipsis">
                  测试工程师@深圳
                </div>
                <div class="dot">·</div>
                <a
                  href="/pin/5e99e6d3e51d451a6532f1b5"
                  target="_blank"
                  rel=""
                  class="time-box"
                >
                  <time
                    datetime="2020-04-17T17:26:43.760Z"
                    title="Sat Apr 18 2020 01:26:43 GMT+0800 (中国标准时间)"
                    class="time"
                    >12小时前
                  </time>
                </a>
              </div>
            </div>
          </div>
          <div class="header-action">
            <button class="subscribe-btn follow-button">
              关注
            </button>
            <div class="pin-header-more header-menu">
              <div class="more-button">
                <svg
                  t="1529034629100"
                  viewBox="0 0 1024 1024"
                  version="1.1"
                  xmlns="http://www.w3.org/2000/svg"
                  p-id="1948"
                  xmlns:xlink="http://www.w3.org/1999/xlink"
                  width="24"
                  height="24"
                  class="icon"
                >
                  <path
                    d="M804.606221 432.282401c120.691803 0 119.469975 187.388854-1.465374 187.388854C682.449044 619.671255 683.426301 432.282401 804.606221 432.282401z"
                    p-id="1949"
                    fill="#b8c1cc"
                  ></path>
                  <path
                    d="M511.428995 432.282401c120.691803 0 119.469975 187.388854-1.465374 187.388854C389.271818 619.671255 390.249075 432.282401 511.428995 432.282401z"
                    p-id="1950"
                    fill="#b8c1cc"
                  ></path>
                  <path
                    d="M218.251769 432.282401c120.691803 0 119.469975 187.388854-1.465374 187.388854C96.094592 619.671255 97.071849 432.282401 218.251769 432.282401z"
                    p-id="1951"
                    fill="#b8c1cc"
                  ></path>
                </svg>
              </div>
              <div class="dropdown" style="display: none;">
                <div class="dropdown-caret"></div>
                <ul class="dropdown-menu">
                  <li>举报</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
        <div class="pin-content-row">
          <div class="content-box content-box">
            <span
              >关于熬夜这件事 或许真是越熬越上头
              三年又三年？三年又三年？？</span
            >
            <div class="limit-ctl-box"><!----><!----></div>
          </div>
        </div>
        <div class="pin-action-row">
          <div class="action-box">
            <div class="like-action action">
              <div class="action-title-box">
                <svg
                  aria-hidden="true"
                  width="20"
                  height="20"
                  viewBox="0 0 20 20"
                  class="icon like-icon"
                >
                  <g fill="none" fill-rule="evenodd">
                    <path d="M0 0h20v20H0z"></path>
                    <path
                      stroke="#8A93A0"
                      stroke-linejoin="round"
                      d="M4.58 8.25V17h-1.4C2.53 17 2 16.382 2 15.624V9.735c0-.79.552-1.485 1.18-1.485h1.4zM11.322 2c1.011.019 1.614.833 1.823 1.235.382.735.392 1.946.13 2.724-.236.704-.785 1.629-.785 1.629h4.11c.434 0 .838.206 1.107.563.273.365.363.84.24 1.272l-1.86 6.513A1.425 1.425 0 0 1 14.724 17H6.645V7.898C8.502 7.51 9.643 4.59 9.852 3.249A1.47 1.47 0 0 1 11.322 2z"
                    ></path>
                  </g>
                </svg>
                <span class="action-title">7</span>
              </div>
            </div>
            <div class="comment-action action">
              <div class="action-title-box">
                <svg
                  aria-hidden="true"
                  width="20"
                  height="20"
                  viewBox="0 0 20 20"
                  class="icon comment-icon"
                >
                  <g fill="none" fill-rule="evenodd">
                    <path d="M0 0h20v20H0z"></path>
                    <path
                      stroke="#8A93A0"
                      stroke-linejoin="round"
                      d="M10 17c-4.142 0-7.5-2.91-7.5-6.5S5.858 4 10 4c4.142 0 7.5 2.91 7.5 6.5 0 1.416-.522 2.726-1.41 3.794-.129.156.41 3.206.41 3.206l-3.265-1.134c-.998.369-2.077.634-3.235.634z"
                    ></path>
                  </g>
                </svg>
                <span class="action-title">21</span>
              </div>
            </div>
            <div class="share-action action">
              <div class="action-title-box">
                <svg
                  aria-hidden="true"
                  width="20"
                  height="20"
                  viewBox="0 0 20 20"
                  class="icon share-icon"
                >
                  <g fill="none" fill-rule="evenodd">
                    <path d="M0 0h20v20H0z"></path>
                    <g
                      stroke="#8A93A0"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    >
                      <path
                        d="M10 2.5v10M13.5 7h2.001c.552 0 .999.437.999.96v8.621c0 1.919-1.447 1.919-1.999 1.919H5.5c-.552 0-1.999 0-1.999-1.919v-8.62c0-.53.443-.961.999-.961H6.5M7.519 4.538L10.019 2l2.5 2.513"
                      ></path>
                    </g>
                  </g>
                </svg>
                <span class="action-title">分享</span>
              </div>
              <!---->
            </div>
          </div>
        </div>
      </div>
    </li>
  </ul>
</template>

<script>
export default {
  props: {
    tweets: {
      type: Array,
      default() {
        return []
      },
      required: false
    }
  }
}
</script>

<style scoped lang="scss">
.tweets-list {
  .tweets {
    margin-bottom: 8px;
    background: #fff;
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);

    .pin {
      background-color: #fff;
      border-radius: 0.2rem;

      .pin-header-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 16px 2rem 0 20px;
      }

      .pin-content-row {
        margin: 5px 0;
        position: relative;
      }
      .action-box {
        display: flex;
        position: relative;
        margin-top: 1.333rem;
        height: 34px;
        border-top: 1px solid #ebebeb;

        .action {
          flex: 1 1 33.333%;
          display: flex;
          align-items: center;
          justify-content: center;
          position: relative;
          height: 100%;
          cursor: pointer;
          user-select: none;

          &:not(:last-child):after {
            content: '';
            position: absolute;
            top: 50%;
            right: 0;
            margin-top: -1rem;
            width: 1px;
            height: 2rem;
            background-color: #ebebeb;
          }

          .action-title-box {
            display: flex;
            align-items: center;
            justify-content: center;

            .icon {
              width: 18px;
              height: 18px;
            }

            .action-title {
              margin-left: 0.3em;
              font-size: 13px;
              font-weight: 500;
              color: #8a93a0;
            }
          }
        }
      }
    }
  }
}
</style>

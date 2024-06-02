<template>
  <div class="widget">
    <div class="widget-header">签到</div>
    <div class="widget-content checkin">
      <div class="checkedin">
        <div class="gold-icon-box">
          <div>
            <span
              ><svg
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 8l3 5m0 0l3-5m-3 5v4m-3-5h6m-6 3h6m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                /></svg
            ></span>
          </div>
        </div>
        <div class="gold-info-box">
          <div v-if="isLogin">
            <span class="gold-info">{{ user.score || 0 }}</span>
          </div>
          <div v-else>
            <span class="gold-info">0</span>
          </div>
        </div>

        <div v-if="checkIn && checkIn.checkIn" class="checkedin-btn-box">
          <a class="button checkedin-btn" disabled>
            <span class="checkedin-btn-icon">
              <svg
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
            </span>
            <span>今日已签到</span>
          </a>
        </div>
        <div v-else class="checkedin-btn-box">
          <a class="checkedin-btn" @click="doCheckIn">
            <span class="checkedin-btn-icon">
              <svg
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
            </span>
            <span>立即签到</span>
          </a>
        </div>
      </div>
      <div v-if="checkIn && checkIn.checkIn" class="checkedin-tips-box">
        <span class="checkedin-tips-icon"
          ><svg
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"
            /></svg
        ></span>
        <span class="checkedin-tips-info"
          >你已经连续签到&nbsp;<b class="checkedin-tips-day">{{
            checkIn.consecutiveDays
          }}</b
          >&nbsp;天啦 !</span
        >
      </div>

      <div v-if="checkInRank && checkInRank.length" class="rank">
        <div class="rank-title">今日排行</div>
        <ul>
          <li v-for="rank in checkInRank" :key="rank.id" class="rounded">
            <my-avatar :user="rank.user" :size="30" class="rank-user-avatar" />
            <div class="rank-user-info">
              <nuxt-link :to="`/user/${rank.user.id}`">
                {{ rank.user.nickname }}
              </nuxt-link>
              <p>@{{ usePrettyDate(rank.updateTime) }}</p>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
const userStore = useUserStore();

const isLogin = computed(() => {
  return userStore.user !== null;
});

const user = computed(() => {
  return userStore.user;
});

const { data: checkIn, refresh: refreshCheckIn } = await useAsyncData(() =>
  useMyFetch(`/api/checkin/checkin`)
);
const { data: checkInRank, refresh: refreshCheckInRank } = await useAsyncData(
  () => useMyFetch(`/api/checkin/rank`)
);

async function doCheckIn() {
  try {
    checkIn.value = await useHttpPostForm("/api/checkin/checkin");
    useMsgSuccess("签到成功");
    refreshCheckIn();
    refreshCheckInRank();
  } catch (e) {
    useCatchError(e);
  }
}
</script>

<style lang="scss" scoped>
.checkin {
  .checkedin {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .gold-icon-box {
      width: 25%;
      span {
        color: var(--text-color4);
        font-size: 0.875rem;
        line-height: 1.25rem;
        display: flex;
        svg {
          width: 1.5rem;
          height: 1.5rem;
        }
      }
    }
    .gold-info-box {
      width: 25%;
      .gold-info {
        color: var(--text-color4);
        font-size: 0.875rem;
        line-height: 1.25rem;
        display: flex;
      }
    }
    .checkedin-btn-box {
      display: flex;
      flex-direction: column;
      width: 50%;
      .checkedin-btn {
        display: flex;
        outline: none;
        align-items: center;
        justify-content: center;
        padding: 6px 0.5rem;
        font-size: 0.75rem;
        line-height: 1rem;
        font-weight: 500;
        white-space: nowrap;
        border-radius: 0.25rem;
        background-color: var(--bg-color);
        border: 1px solid var(--border-color2);
        color: var(--text-color3);
        &:hover {
          background-color: var(--bg-color2);
          border-color: var(--border-color2);
        }
        .checkedin-btn-icon {
          margin-right: 0.5rem;
          display: flex;
          svg {
            width: 1rem;
            height: 1rem;
          }
        }
      }
    }
  }
  .checkedin-tips-box {
    display: flex;
    flex-direction: row;
    justify-content: center;
    text-align: center;
    align-items: center;
    color: var(--text-color3);
    padding-top: 0.5rem;
    padding-bottom: 0.5rem;
    margin-top: 0.5rem;
    margin-bottom: 0.5rem;
    border-radius: 0.25rem;
    background: var(--bg-color2);
    .checkedin-tips-icon {
      margin-right: 0.5rem;
      display: flex;
      svg {
        width: 1.25rem;
        height: 1.25rem;
      }
    }
    .checkedin-tips-info {
      font-size: 0.875rem;
      line-height: 1.25rem;
      display: flex;
      .checkedin-tips-day {
        color: #00a1d6; // TODO
      }
    }
  }
  .rank {
    border-top: 1px solid var(--border-color);
    margin-top: 10px;
    padding-top: 10px;
    .rank-title {
      font-size: 14px;
      font-weight: 600;
    }
    li {
      display: flex;
      list-style: none;
      margin: 8px 0;
      font-size: 13px;
      position: relative;
      background: var(--bg-color);
      padding: 0.5rem;
      border-radius: 0.25rem;
      cursor: pointer;

      &:hover {
        background: var(--bg-color2);
      }

      &:not(:last-child) {
        border-bottom: 1px solid var(--border-color);
      }

      .rank-user-avatar {
        min-width: 30px;
        margin-top: 0.5rem;
      }

      .rank-user-info {
        width: 100%;
        margin-left: 10px;
        line-height: 1.4;
        font-size: 12px;
        a {
          color: var(--text-color2);
          font-weight: 700;
          &:hover {
            color: var(--text-color1);
            text-decoration: none;
          }
        }
        p {
          margin-top: 0.5rem;
        }
      }
    }
  }
}
</style>

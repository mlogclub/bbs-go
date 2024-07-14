<template>
  <div class="banner">
    <div class="banner-inner">
      <a-carousel
        class="carousel"
        :auto-play="true"
        show-arrow="hover"
        animation-name="slide"
      >
        <a-carousel-item v-for="item in carouselItem" :key="item.slogan">
          <div :key="item.slogan" class="carousel-item">
            <div class="carousel-content">
              <div class="carousel-content-title">{{ item.slogan }}</div>
              <div class="carousel-content-sub-title">{{ item.subSlogan }}</div>
            </div>
            <!-- <img class="carousel-image" :src="item.image" /> -->
          </div>
        </a-carousel-item>
      </a-carousel>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { computed } from 'vue';
  import { useI18n } from 'vue-i18n';
  import bannerImage from '@/assets/images/logo364x364.png';

  const { t } = useI18n();
  const carouselItem = computed(() => [
    {
      slogan: t('login.banner.slogan1'),
      subSlogan: t('login.banner.subSlogan1'),
      image: bannerImage,
    },
    {
      slogan: t('login.banner.slogan2'),
      subSlogan: t('login.banner.subSlogan2'),
      image: bannerImage,
    },
    {
      slogan: t('login.banner.slogan3'),
      subSlogan: t('login.banner.subSlogan3'),
      image: bannerImage,
    },
  ]);
</script>

<style lang="less" scoped>
  .banner {
    display: flex;
    align-items: center;
    justify-content: center;

    &-inner {
      flex: 1;
      height: 100%;
    }
  }

  .carousel {
    height: 100%;

    &-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
    }

    &-content {
      height: 150px;
      padding: 0 100px;
      &-title {
        color: var(--color-fill-1);
        font-weight: 500;
        font-size: 20px;
        line-height: 28px;
      }

      &-sub-title {
        margin-top: 18px;
        color: var(--color-text-3);
        font-size: 14px;
        line-height: 22px;
      }
    }

    &-image {
      width: 320px;
      margin-top: 30px;
    }
  }
</style>

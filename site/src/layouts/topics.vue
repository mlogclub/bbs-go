<template>
  <div>
    <MyHeader />

    <section class="main">
      <div class="container main-container left-main size-320">
        <div class="left-container">
          <div class="main-content no-padding no-bg topics-wrapper">
            <div class="topics-nav">
              <topics-nav />
            </div>
            <slot />
          </div>
        </div>
        <div class="right-container">
          <check-in />
          <site-notice />
          <score-rank />
          <friend-links />
        </div>
      </div>
    </section>

    <MyFooter />
  </div>
</template>

<script setup></script>

<style lang="scss" scoped></style>

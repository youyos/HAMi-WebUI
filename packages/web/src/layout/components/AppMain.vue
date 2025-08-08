<template>
  <div id="content" :class="hasParentWindow && 'impose'">
    <!-- <transition name="fade-transform" mode="out-in"> -->
    <!-- <keep-alive> -->
    <router-view />
    <!-- </keep-alive> -->
    <!-- </transition> -->
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import useParentAction from '~/vgpu/hooks/useParentAction';


const { hasParentWindow } = useParentAction();

const route = useRoute();

const isHome = computed(() => route.fullPath === '/admin/home');
</script>

<style lang="scss" scoped>
#content {
  margin: 0;
  flex: 1;
  overflow: auto;
  height: 100%;
  padding-left: 20px;

  &.impose {
    margin: 0 auto;
    overflow: inherit;
    max-width: 1200px;
    padding-left: 0px;

  }
}
</style>

<template>
  <TopBar v-if="!hasParentWindow" />
  <div class="page" :style="{ padding: hasParentWindow ? '20px 30px' : '76px 20px 20px 20px' }">
    <sidebar v-if="!hasParentWindow && !isNoSidebar()" />
    <app-main />
  </div>
</template>

<script>
import '@tabler/core/dist/css/tabler.min.css';
import { AppMain, Sidebar, TopBar } from './components';
import { useRouter } from 'vue-router'
import useParentAction from '~/vgpu/hooks/useParentAction';
import { parseUrl } from '@/utils';


export default {
  name: 'Layout',
  components: {
    AppMain,
    Sidebar,
    TopBar,
  },
  setup() {
    const router = useRouter();
    const { hasParentWindow } = useParentAction();

    return { router, hasParentWindow };
  },
  created() {
    this.isNoSidebar();
    this.setupMessageListener();
  },
  beforeDestroy() {
    window.removeEventListener('message', this.handleMessage);
  },
  mounted() {
    this.sendLoadedMessage();
  },
  methods: {
    isNoSidebar() {
      const noSidebarPaths = ['/admin/home', '/admin/message-center', "/admin/about-system", "/admin/settings/config-map"];
      return noSidebarPaths.includes(this.$route.fullPath);
    },
    setupMessageListener() {
      window.addEventListener('message', this.handleMessage);
    },
    handleMessage(event) {
      try {
        const messageData = JSON.parse(event.data);
        if (messageData.type === "ChangeTheRoute") {
          // const { pathname, query } = parseUrl(messageData.data)
          // console.log(messageData.data, 'messageData.data')
          // this.router.push({
          //   path: pathname,
          //   query
          // });
          this.router.replace(messageData.data)
        }
      } catch (e) { }
    },
    sendLoadedMessage() {
      if (window.parent !== window) {
        const message = {
          type: 'ComponentLoaded',
          data: window.location.pathname
        };
        window.parent.postMessage(JSON.stringify(message), '*');
      }
    }
  },
};
</script>

<style lang="scss" scoped>
@import '~@/styles/mixin.scss';
@import '~@/styles/variables.scss';

.page {
  display: flex;
  flex-direction: row;
  // padding: 20px 30px;
  // padding: 20px;
  // padding-top: 76px;
  height: calc(100vh);
}

.app-wrapper {
  @include clearfix;
  position: relative;
  height: 100%;
  width: 100%;

  /* &.mobile.openSidebar {
    position: fixed;
    top: 0;
  } */
}

.drawer-bg {
  background: #000;
  opacity: 0.3;
  width: 100%;
  top: 0;
  height: 100%;
  position: absolute;
  z-index: 999;
}

.fixed-header {
  position: fixed;
  top: 0;
  right: 0;
  z-index: 9;
  width: calc(100% - #{$sideBarWidth});
  transition: width 0.28s;
}

.hideSidebar .fixed-header {
  width: calc(100% - 54px);
}

.mobile .fixed-header {
  width: 100%;
}
</style>
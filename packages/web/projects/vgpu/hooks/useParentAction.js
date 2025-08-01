import { useRouter } from 'vue-router';

export default function useParentAction() {
  const router = useRouter();

  const sendRouteChange = (url) => {
    if (window.parent !== window) {
      // 在 iframe 中，通知父窗口
      const message = {
        type: 'ChangeThePath',
        data: url
      };
      window.parent.postMessage(JSON.stringify(message), '*');
    } else {
      // 不在 iframe 中，直接使用 router 跳转
      router.push(url);
    }
  };
  const hasParentWindow = window.parent !== window;

  return { sendRouteChange, hasParentWindow };
}
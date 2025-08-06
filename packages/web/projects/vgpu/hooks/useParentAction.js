import { useRouter } from 'vue-router';

export default function useParentAction() {
  const router = useRouter();

  const sendRouteChange = (url, type) => {
    if (window.parent !== window) {
      // 在 iframe 中，通知父窗口
      let message = {
        type: 'ChangeThePath',
        data: url,
      };
      if (type === 'open') {
        message = {
          type: 'OpenThePath',
          data: url,
        };
      }
      if (type === 'back') {
        message = {
          type: 'BackThePath',
          data: url,
        };
      }
      window.parent.postMessage(JSON.stringify(message), '*');
    } else {
      // 不在 iframe 中，直接使用 router 跳转
      if (type === 'open') {
        window.open(url);
      } else if (type === 'back') {
        router.go(-1);
      } else {
        router.push(url);
      }
    }
  };
  const hasParentWindow = window.parent !== window;

  return { sendRouteChange, hasParentWindow };
}

<template>
  <el-date-picker
    ref="pickerRef"
    v-model="value"
    :type="type"
    range-separator="至"
    start-placeholder="开始时间"
    end-placeholder="结束时间"
    unlink-panels
    :shortcuts="type.includes('range') && shortcuts"
    class="date-picker"
    :disabled-date="disabledDate"
    @visible-change="onVisibleChange"
    v-bind="$attrs"
  />
</template>

<script setup lang="jsx">
import { computed, defineProps, defineEmits, ref } from 'vue';

// Props 和 Emits
const props = defineProps({
  modelValue: {},
  type: { type: String, default: 'date' },
  parse: { type: Function, default: (times) => times },
});
const emits = defineEmits(['update:modelValue']);

// picker 引用 和 当前面板是否显示
const pickerRef = ref();
const visible = ref(false);

// 显示状态变化监听
const onVisibleChange = (val) => {
  visible.value = val;
};

// 双向绑定逻辑
const value = computed({
  get() {
    return props.modelValue;
  },
  set(val) {
    emits('update:modelValue', props.parse(val));
  },
});

// 快捷时间生成函数
function genShortcut(text, hoursAgo) {
  return {
    text,
    value: () => {
      const end = new Date();
      const start = new Date(end.getTime() - hoursAgo * 3600 * 1000);

      // 只关闭当前打开的面板，避免影响其他 time-picker
      setTimeout(() => {
        if (visible.value) {
          pickerRef.value?.handleClose?.();
        }
      }, 0);

      return [start, end];
    },
  };
}

// 快捷选项列表
const shortcuts = [
  genShortcut('前 1 小时', 1),
  genShortcut('前 6 小时', 6),
  genShortcut('前 12 小时', 12),
  genShortcut('前 1 天', 24),
  genShortcut('前 2 天', 48),
  genShortcut('前 3 天', 72),
  genShortcut('前 1 周', 168),
];

// 禁用未来时间
const disabledDate = (time) => {
  return time.getTime() >= Date.now();
};
</script>

<style scoped>
.date-picker {
  max-width: 450px;
}
</style>

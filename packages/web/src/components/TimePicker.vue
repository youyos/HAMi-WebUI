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
    v-bind="$attrs"
  />
</template>

<script setup lang="jsx">
import { computed, defineProps, defineEmits, ref } from 'vue';
import { ElDatePicker } from 'element-plus';

// props & emits
const props = defineProps({
  modelValue: {},
  type: { type: String, default: 'date' },
  parse: { type: Function, default: (times) => times },
});
const emits = defineEmits(['update:modelValue']);

// picker ref
const pickerRef = ref();

// 双向绑定
const value = computed({
  get() {
    return props.modelValue;
  },
  set(val) {
    emits('update:modelValue', props.parse(val));
  },
});

// 生成快捷选项的封装函数
function genShortcut(text, hoursAgo) {
  return {
    text,
    value: () => {
      const end = new Date();
      const start = new Date(end.getTime() - hoursAgo * 3600 * 1000);
      // 延迟关闭弹窗，确保值更新后再关闭
      setTimeout(() => {
        pickerRef.value?.handleClose?.();
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

// 禁止选择未来时间
const disabledDate = (time) => {
  return time.getTime() >= Date.now();
};
</script>

<style>
.date-picker {
  max-width: 450px;
}
</style>

<template>

  <block-box>
    <ul class="card-gauges">
      <li v-for="(item, index) in gaugeConfig.slice(0, 4)" :key="index">
        <Gauge v-bind="item" />
      </li>
    </ul>
    <ul class="card-gauges" style="margin-top: 20px;">
      <li v-for="(item, index) in gaugeConfig.slice(4, 7)" :key="index">
        <Gauge v-bind="item" />
      </li>
    </ul>
  </block-box>

  <div class="line-box">
    <block-box title="资源分配趋势（%）">
      <template #extra>
        <time-picker v-model="times" type="datetimerange" size="small" />
      </template>
      <div style="height: 200px">
        <echarts-plus :options="getRangeOptions({
          cpu: gaugeConfig[7].data,
          internal: gaugeConfig[8].data,
          core: gaugeConfig[4].data,
          memory: gaugeConfig[5].data,
        })
          " />
      </div>
    </block-box>
    <block-box title="资源使用趋势（%）">
      <template #extra>
        <time-picker v-model="times" type="datetimerange" size="small" />
      </template>
      <div style="height: 200px">
        <echarts-plus :options="getRangeOptions({
          cpu: gaugeConfig[0].data,
          internal: gaugeConfig[1].data,
          core: gaugeConfig[6].data,
          memory: gaugeConfig[3].data,
        })
          " />
      </div>
    </block-box>
  </div>

  <block-box title="显卡列表">
    <CardList :hideTitle="true" :filters="data" />
  </block-box>

  <block-box title="任务列表">
    <TaskList :hideTitle="true" :filters="data" />
  </block-box>
</template>

<script setup lang="jsx">
import BlockBox from '@/components/BlockBox.vue';
import { ref, defineProps, watchEffect } from 'vue';
import CardList from '~/vgpu/views/card/admin/index.vue';
import TaskList from '~/vgpu/views/task/admin/index.vue';
import Gauge from '~/vgpu/components/gauge.vue';
import useInstantVector from '~/vgpu/hooks/useInstantVector';
import EchartsPlus from '@/components/Echarts-plus.vue';
import { getRangeOptions } from './getOptions';

const props = defineProps(['data'])

watchEffect(() => {
  console.log(props.data, 'data2')
})

const end = new Date();
const start = new Date();
start.setTime(start.getTime() - 3600 * 1000);

const times = ref([start, end]);

const gaugeConfig = useInstantVector(
  [
    {
      title: 'CPU 使用率',
      percent: 0,
      query: `count(node_cpu_seconds_total{mode="idle", instance=~"$node"}) by (instance)*(1 - avg(rate(node_cpu_seconds_total{mode="idle", instance=~"$node"}[5m])) by (instance))`,
      totalQuery: `count(node_cpu_seconds_total{mode="idle", instance=~"$node"}) by (instance)`,
      percentQuery: `100 * (1 - avg by(instance)(irate(node_cpu_seconds_total{mode="idle", instance=~"$node"}[1m])))`,
      total: 0,
      used: 0,
      unit: '核',
    },
    {
      title: '内存 使用率',
      percent: 0,
      query: `avg(node_memory_MemTotal_bytes{instance=~"$node"} - node_memory_MemAvailable_bytes{instance=~"$node"}) by (instance) / 1024 / 1024 / 1024`,
      totalQuery: `avg(node_memory_MemTotal_bytes{instance=~"$node"}) by (instance) / 1024 / 1024 / 1024`,
      percentQuery: `100 * (1 - node_memory_MemAvailable_bytes{instance=~"$node"} / node_memory_MemTotal_bytes{instance=~"$node"})`,
      total: 0,
      used: 0,
      unit: 'GiB',
    },
    {
      title: '磁盘 使用率',
      percent: 0,
      query: `sum(node_filesystem_size_bytes{instance=~"$node", fstype=~"ext4|xfs", mountpoint!~"/var/lib/kubelet/pods.*"} - node_filesystem_free_bytes{instance=~"$node", fstype=~"ext4|xfs", mountpoint!~"/var/lib/kubelet/pods.*"}) by (instance) / 1024 / 1024 / 1024`,
      totalQuery: `sum(node_filesystem_size_bytes{instance=~"$node", fstype=~"ext4|xfs", mountpoint!~"/var/lib/kubelet/pods.*"}) by (instance) / 1024 / 1024 / 1024`,
      percentQuery: ``,
      total: 0,
      used: 0,
      unit: 'GiB',
    },
    {
      title: '显存使用率',
      percent: 0,
      query: `avg(sum(hami_memory_used{node=~"$node"}) by (instance)) / 1024`,
      totalQuery: `avg(sum(hami_memory_size{node=~"$node"}) by (instance))/1024`,
      percentQuery: `(avg(sum(hami_memory_used{node=~"$node"}) by (instance)) / 1024)/(avg(sum(hami_memory_size{node=~"$node"}) by (instance))/1024)*100`,
      total: 0,
      used: 0,
      unit: 'GiB',
    },
    {
      title: '算力分配率',
      percent: 0,
      query: `avg(sum(hami_container_vcore_allocated{node=~"$node"}) by (instance))`,
      totalQuery: `avg(sum(hami_core_size{node=~"$node"}) by (instance))`,
      percentQuery: `avg(sum(hami_container_vcore_allocated{node=~"$node"}) by (instance)) / avg(sum(hami_core_size{node=~"$node"}) by (instance)) *100`,
      total: 0,
      used: 0,
      unit: ' ',
    },
    {
      title: '显存分配率',
      percent: 0,
      query: `avg(sum(hami_container_vmemory_allocated{node=~"$node"}) by (instance)) / 1024`,
      totalQuery: `avg(sum(hami_memory_size{node=~"$node"}) by (instance)) / 1024`,
      percentQuery: `(avg(sum(hami_container_vmemory_allocated{node=~"$node"}) by (instance)) / 1024) /(avg(sum(hami_memory_size{node=~"$node"}) by (instance)) / 1024) *100`,
      total: 0,
      used: 0,
      unit: 'GiB',
    },
    {
      title: '算力使用率',
      percent: 0,
      query: `avg(sum(hami_core_util{node=~"$node"}) by (instance))`,
      percentQuery: `avg(sum(hami_core_util_avg{node=~"$node"}) by (instance))`,
      totalQuery: `avg(sum(hami_core_size{node=~"$node"}) by (instance))`,
      total: 100,
      used: 0,
      unit: ' ',
    },
    {
      title: 'CPU 分配率',
      percent: 0,
      query: ``,
      totalQuery: ``,
      percentQuery: `avg(sum(hami_container_vcore_allocated{node=~"$node"}) by (instance) / sum(hami_core_size{node=~"$node"}) by (instance) * 100)`,
      total: 0,
      used: 0,
      unit: '核',
    },
    {
      title: '内存 分配率',
      percent: 0,
      query: ``,
      totalQuery: ``,
      percentQuery: `avg(sum(hami_container_vmemory_allocated{node=~"$node"}) by (instance) / sum(hami_memory_size{node=~"$node"}) by (instance) * 100)`,
      total: 0,
      used: 0,
      unit: 'GiB',
    },
  ],
  (query) => query.replaceAll(`$node`, props?.data?.nodeName),
  times,
);
</script>

<style scoped lang="scss">
.card-gauges {
  margin: 0;
  padding: 0;
  list-style: none;
  display: flex;
  height: 200px;

  li {
    flex: 0.25;
  }
}

.line-box {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  column-gap: 20px;
}

.node-block {
  display: flex;
  flex-direction: column;

  .home-block-content {
    flex: 1;
  }
}
</style>

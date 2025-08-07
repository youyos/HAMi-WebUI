<template>
  <div>
    <back-header>
      资源池详情
    </back-header>
    <table-plus :api="pollApi.getDetailNodeList({ pool_id: route.params.uid })" :columns="columns"
      :rowAction="rowAction" :hasPagination="false" style="margin-bottom: 15px;" hideTag ref="table" staticPage
      :hasActionBar="false">
    </table-plus>
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
      <CardList :hideTitle="true" :filters="{ nodeUid: detail.uid }" />
    </block-box>

    <block-box title="任务列表">
      <TaskList :hideTitle="true" :filters="{ nodeUid: detail.uid }" />
    </block-box>
  </div>
</template>

<script setup lang="jsx">
import BackHeader from '@/components/BackHeader.vue';
import { useRoute } from 'vue-router';
import BlockBox from '@/components/BlockBox.vue';
import { computed, onMounted, ref, watch, watchEffect } from 'vue';
import CardList from '~/vgpu/views/card/admin/index.vue';
import TaskList from '~/vgpu/views/task/admin/index.vue';
import Gauge from '~/vgpu/components/gauge.vue';
import useInstantVector from '~/vgpu/hooks/useInstantVector';
import EchartsPlus from '@/components/Echarts-plus.vue';
import pollApi from '~/vgpu/api/poll';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getRangeOptions } from './getOptions';
import { bytesToGB, roundToDecimal } from "@/utils";
import useParentAction from '~/vgpu/hooks/useParentAction';

const { sendRouteChange } = useParentAction();

const route = useRoute();

const table = ref();

const detail = ref({});

const end = new Date();
const start = new Date();
start.setTime(start.getTime() - 3600 * 1000);

const times = ref([start, end]);


const columns = [
  {
    title: '节点名称',
    dataIndex: 'name',
    render: ({ uid, name }) => (
      <text-plus text={name} to={`/admin/vgpu/node/admin/${uid}?nodeName=${name}`} />
    ),
  },
  {
    title: '节点 IP',
    dataIndex: 'ip',
  },
  {
    title: '节点状态',
    dataIndex: 'isSchedulable',
    render: ({ isSchedulable, isExternal }) => (
      <el-tag disable-transitions type={isExternal ? 'warning' : (isSchedulable ? 'success' : 'danger')}>
        {isExternal ? '未纳管' : (isSchedulable ? '可调度' : '禁止调度')}
      </el-tag>
    )
  },
  {
    title: '显卡型号',
    dataIndex: 'type',
  },
  {
    title: 'CPU',
    dataIndex: 'coreTotal',
    render: ({ coreTotal }) => `${coreTotal}核`,
  },
  {
    title: '内存',
    dataIndex: 'memoryTotal',
    render: ({ memoryTotal }) => `${bytesToGB(memoryTotal)}GiB`,
  },
  {
    title: '磁盘',
    dataIndex: 'diskSize',
    render: ({ diskSize }) => `${bytesToGB(diskSize)}GiB`,
  },
  {
    title: '显卡数量',
    dataIndex: 'cardCnt',
  },
  {
    title: 'vGPU',
    dataIndex: 'used',
    render: ({ vgpuTotal, vgpuUsed, isExternal }) => (
      <span>
        {isExternal ? '--' : vgpuUsed}/{isExternal ? '--' : vgpuTotal}
      </span>
    ),
  },
  {
    title: '算力(已分配/总量)',
    width: 120,
    dataIndex: 'used',
    render: ({ coreTotal, coreUsed, isExternal }) => (
      <span>
        {isExternal ? '--' : coreUsed}/{coreTotal}
      </span>
    ),
  },
  {
    title: '显存(已分配/总量)',
    dataIndex: 'w',
    width: 120,
    render: ({ memoryTotal, memoryUsed, isExternal }) => (
      <span>
        {isExternal ? '--' : roundToDecimal(memoryUsed / 1024, 1)}/
        {roundToDecimal(memoryTotal / 1024, 1)} GiB
      </span>
    ),
  },
];

const rowAction = [
  {
    title: '查看详情',
    onClick: (row) => {
      sendRouteChange(`/admin/vgpu/node/admin/${row.uid}?nodeName=${row.name}`);
    },
  },
  {
    title: '移除',
    onClick: async (row) => {
      ElMessageBox.confirm(`确定要移除当前节点吗？`, '操作确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
        .then(async () => {
          try {
            await pollApi.remove(
              {
                node_id: row.nodeId,
              }
            ).then(
              () => {
                setTimeout(() => {
                  ElMessage.success('节点移除成功');
                  table.value.fetchData();
                }, 500);
              }
            )
          } catch (error) {
            ElMessage.error(error.message);
          }
        })
        .catch(() => { });
    },
  },
];

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
  (query) => query.replaceAll(`$node`, detail.value.name),
  times,
);

onMounted(async () => {
});
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

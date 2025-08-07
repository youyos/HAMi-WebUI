<template>
  <back-header>
    资源池管理 > {{ route.query?.name || '' }}
  </back-header>
  <table-plus v-loading="loading" :dataSource="list" :columns="columns" :rowAction="rowAction" :hasPagination="false"
    style="margin-bottom: 15px; height: auto;" hideTag ref="table" static :hasActionBar="false">
  </table-plus>
  <DetailBox v-if="!!data.nodeName" :key="data.nodeName" :data="data" />
</template>

<script setup lang="jsx">
import BackHeader from '@/components/BackHeader.vue';
import DetailBox from './DetailBox.vue';
import { useRoute } from 'vue-router';
import { ref, computed, watchEffect, onMounted } from 'vue';
import pollApi from '~/vgpu/api/poll';
import { ElMessage, ElMessageBox } from 'element-plus';
import { bytesToGB, roundToDecimal } from "@/utils";
import useParentAction from '~/vgpu/hooks/useParentAction';

const { sendRouteChange } = useParentAction();

const route = useRoute();

const table = ref();

const list = ref([]);
const loading = ref(true);

const data = computed(() => {
  const result = {
    nodeName: '',
    nodeUid: ''
  };
  list.value.forEach(node => {
    result.nodeName = result.nodeName + (result.nodeName ? '|' : '') + node.name
    result.nodeUid = result.nodeUid + (result.nodeUid ? '|' : '') + node.uid
  });
  return result;
})

watchEffect(() => {
  console.log(data.value, 'data')
})

const getList = async () => {
  loading.value = true
  const res = await pollApi.getDetailNodeList({ pool_id: route.params.uid })
  loading.value = false
  list.value = res.list
}

onMounted(() => {
  getList();
});

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
    hidden: () => list.value.length < 2,
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
                  getList();
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
</script>

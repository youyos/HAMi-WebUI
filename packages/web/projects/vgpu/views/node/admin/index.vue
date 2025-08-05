<template>
  <list-header description="节点管理用于管理和监控计算节点的状态。它可以启用或禁用节点，查看节点上的物理GPU卡，以及监控节点上运行的所有任务。">
    <template #actions>
      <el-button @click="handleAdd" style="margin-right: 24px;" type="primary" round>发现节点</el-button>
    </template>
  </list-header>

  <preview-bar :handle-click=handleClick :key="componentKey" />

  <table-plus :api="nodeApi.getNodeList()" :columns="columns" :rowAction="rowAction" :searchSchema="searchSchema"
    :hasPagination="false" style="height: auto" hideTag ref="table" staticPage>
  </table-plus>

  <el-dialog @close="nodeSelect = []" v-model="dialogVisible" title="添加节点" width="600" :before-close="handleClose">
    <div v-loading="loading">
      <template v-if="nodeList && nodeList.length > 0">
        <div style="display: flex; align-items: center;" v-for="{ nodeIp }, index in nodeList" :key="nodeIp">
          <el-checkbox :model-value="nodeSelect.includes(nodeIp)" @change="handleCheckboxChange(nodeIp)" />
          <span style="color: #0B1524; margin-left: 8px;">{{ nodeIp }}</span>
        </div>
      </template>
      <template v-else>
        <div>
          <el-empty description="暂无节点数据" />
        </div>
      </template>
    </div>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button :loading="btnLoading" type="primary" @click="handleOk">确认</el-button>
    </template>
  </el-dialog>

</template>

<script setup lang="jsx">
import nodeApi from '~/vgpu/api/node';
import searchSchema from '~/vgpu/views/node/admin/searchSchema';
import PreviewBar from '~/vgpu/components/previewBar.vue';
import { bytesToGB, roundToDecimal } from '@/utils';
import { ElMessage, ElMessageBox } from 'element-plus';
import { ref } from 'vue';
import useParentAction from '~/vgpu/hooks/useParentAction';


const { sendRouteChange } = useParentAction();

const table = ref();

const componentKey = ref(0);

// 节点选择相关
const dialogVisible = ref(false)
const nodeList = ref([])
const nodeSelect = ref([])
const loading = ref(true)
const btnLoading = ref(false)



const handleClick = async (params) => {
  const name = params.data.name;
  const { list } = await nodeApi.getNodes({ filters: {} })
  const node = list.find(node => node.name === name);
  if (node) {
    const uuid = node.uid;
    sendRouteChange(`/admin/vgpu/node/admin/${uuid}?nodeName=${name}`);
  } else {
    ElMessage.error('节点未找到');
  }
};

// 确认操作
const handleOk = async () => {
  if (!nodeSelect.value.length) {
    ElMessage({
      message: '请选择节点',
      type: 'warning',
    })
    return;
  }
  btnLoading.value = true;
  try {
    const node_names = nodeList.value.filter(e => nodeSelect.value.includes(e.nodeIp)).map(e => e.nodeName)
    const res = await nodeApi.joinNodes({
      node_names
    })
    if (res?.code === 200) {
      table.value.fetchData();
      componentKey.value += 1;
      dialogVisible.value = false;
    }
  } finally {
    btnLoading.value = false;
  }
}

// 复选框变化
const handleCheckboxChange = (ip) => {
  const index = nodeSelect.value.indexOf(ip);
  if (index > -1) {
    nodeSelect.value.splice(index, 1);
  } else {
    nodeSelect.value.push(ip);
  }
}

// 添加节点
const handleAdd = async () => {
  dialogVisible.value = true
  loading.value = true
  const res = await nodeApi.discoveredNodes({})
  nodeList.value = res?.list || []
  loading.value = false
}

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
    // filters: [
    //   {
    //     text: '可调度',
    //     value: 'true',
    //   },
    //   {
    //     text: '禁止调度',
    //     value: 'false',
    //   },
    // ],
  },
  {
    title: '显卡型号',
    dataIndex: 'type',
    // filters: (data) => {
    //   const r = data.reduce((all, item) => {
    //     return uniq([...all, ...item.type]);
    //   }, []);
    //
    //   return r.map((item) => ({ text: item, value: item }));
    // },
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
    title: '所属资源池',
    dataIndex: 'resourcePools',
    render: ({ resourcePools }) => `${resourcePools.join('、')}`,
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
    title: '禁用',
    hidden: (row) => !row.isSchedulable,
    onClick: async (row) => {
      ElMessageBox.confirm(`确认对该节点进行禁用操作？`, '操作确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
        .then(async () => {
          try {
            await nodeApi.stop(
              {
                nodeName: row.name,
                status: 'DISABLED'
              }
            ).then(
              () => {
                setTimeout(() => {
                  ElMessage.success('节点禁用成功');
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
  {
    title: '开启',
    hidden: (row) => row.isSchedulable,
    disabled: (row) => row.isExternal,
    onClick: async (row) => {
      ElMessageBox.confirm(`确认对该节点进行开启调度操作？`, '操作确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
        .then(async () => {
          try {
            await nodeApi.stop(
              {
                nodeName: row.name,
                status: 'ENABLE'
              }
            ).then(
              () => {
                setTimeout(() => {
                  ElMessage.success('节点开启调度成功');
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
</script>

<style></style>

<template>
  <list-header description="资源池管理用于统一管理和调度计算资源，支持资源的动态分配、回收与负载均衡，提升资源利用率和系统灵活性。">
    <template #actions>
      <el-button @click="editData = null; dialogVisible = true" style="margin-right: 24px;" type="primary"
        round>创建资源池</el-button>
    </template>
  </list-header>

  <block-box
    v-for="{ poolId, poolName, nodeNum, cpuCores, gpuNum, availableMemory, totalMemory, diskSize }, index in list"
    :key="poolId">
    <el-row style="align-items: center;">
      <div class="left">
        <b class="title">{{ poolName }}</b>
        <div class="tags">
          <span>节点数量&nbsp;&nbsp;{{ nodeNum }}</span>
          <span>CPU数&nbsp;&nbsp;{{ cpuCores }}核</span>
          <span>显卡数量&nbsp;&nbsp;{{ gpuNum }}张</span>
          <span>可用/总内存&nbsp;&nbsp;{{ bytesToGB(availableMemory) }}GB / {{ bytesToGB(totalMemory) }}GB</span>
          <span>磁盘大小&nbsp;&nbsp;{{ bytesToGB(diskSize) }}GB</span>
        </div>
      </div>
      <div class="right">
        <el-button type="text">查看详情</el-button>
        <template v-if="index === 0">
          <el-button type="text">配置</el-button>
        </template>
        <template v-else>
          <el-button type="text">编辑</el-button>
          <el-button type="text">删除</el-button>
        </template>
      </div>
    </el-row>
  </block-box>

  <el-dialog v-model="dialogVisible" :title="editData ? '编辑资源池' : '创建资源池'" width="1180" :before-close="handleClose">
    <el-row :wrap="false" style="align-items: center;">
      <span style="flex-shrink: 0; margin-right: 14px;">资源池名称</span>
      <el-input style="flex: 1;" v-model="input" size="large" />
    </el-row>
    <div style="margin-top: 20px; margin-bottom: 10px;">
      <span>选择节点</span>
      <span style="float: right;">已选<span style="color: #3061D0; margin: 0 5px;">{{ nodeSelect.length }}</span>个节点</span>
    </div>
    <div class="wrap">
      <div class="wrap-left">
        <div style="margin-top: 12px;"
          v-for="{ nodeIp, cpuCores, gpuNum, gpuMemory, totalMemory, diskSize }, index in nodeList" :key="nodeIp">
          <div style="display: flex; align-items: center;">
            <el-checkbox :model-value="nodeSelect.includes(nodeIp)" @change="handleCheckboxChange(nodeIp)" />
            <span style="color: #0B1524; margin-left: 8px;">{{ nodeIp }}</span>
          </div>
          <div class="wrap-row">
            <div>显卡数量<span>{{ gpuNum }}张</span></div>
            <div>显卡大小<span>{{ bytesToGB(gpuMemory) }}GB</span></div>
            <div>内存大小<span>{{ bytesToGB(totalMemory) }}GB</span></div>
            <div>磁盘大小<span>{{ bytesToGB(diskSize) }}GB</span></div>
            <div>CPU<span>{{ cpuCores }}核</span></div>
          </div>
        </div>

      </div>
      <div class="wrap-center">

      </div>
      <div class="wrap-right">
        <div style="margin-top: 12px;"
          v-for="{ nodeIp, cpuCores, gpuNum, gpuMemory, totalMemory, diskSize }, index in nodeList.filter(e => nodeSelect.includes(e.nodeIp))"
          :key="nodeIp">
          <div style="display: flex; align-items: center;">
            <el-icon :size="16" color="red" style="cursor: pointer;" @click="handleCheckboxChange(nodeIp)">
              <Remove />
            </el-icon>
            <span style="color: #0B1524; margin-left: 6px;">{{ nodeIp }}</span>
          </div>
          <div class="wrap-row" style="margin-top: 6px;">
            <div>显卡数量<span>{{ gpuNum }}张</span></div>
            <div>显卡大小<span>{{ bytesToGB(gpuMemory) }}GB</span></div>
            <div>内存大小<span>{{ bytesToGB(totalMemory) }}GB</span></div>
            <div>磁盘大小<span>{{ bytesToGB(diskSize) }}GB</span></div>
            <div>CPU<span>{{ cpuCores }}核</span></div>
          </div>
        </div>
      </div>
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="dialogVisible = false">确认</el-button>
      </span>
    </template>
  </el-dialog>

</template>

<script setup lang="jsx">
import pollApi from '~/vgpu/api/poll';
import { ref, onMounted } from 'vue';
import BlockBox from '@/components/BlockBox.vue';
import { Remove } from '@element-plus/icons-vue'

const list = ref([])
const dialogVisible = ref(false)
const editData = ref(null)
const nodeList = ref([])
const nodeSelect = ref([])
const input = ref('')

const bytesToGB = (bytes) => {
  return Math.round(bytes / (1024 * 1024 * 1024));
}

const handleCheckboxChange = (ip) => {
  const index = nodeSelect.value.indexOf(ip);
  if (index > -1) {
    nodeSelect.value.splice(index, 1);
  } else {
    nodeSelect.value.push(ip);
  }
}

onMounted(async () => {
  console.log(111)
  // list.value = await pollApi.getPollList({});
  list.value = [
    {
      "poolId": "4",
      "poolName": "master资源池",
      "cpuCores": "96",
      "nodeNum": "1",
      "gpuNum": "0",
      "availableMemory": "112892563456",
      "totalMemory": "134472257536",
      "diskSize": "7676310884352",
    },
    {
      "poolId": "5",
      "poolName": "worker资源池",
      "cpuCores": "80",
      "nodeNum": "1",
      "gpuNum": "1",
      "availableMemory": "134327394304",
      "totalMemory": "134432251904",
      "diskSize": "7675682045952",
    }
  ]
  nodeList.value = [
    {
      "nodeName": "k8s1",
      "cpuCores": "96",
      "gpuNum": "0",
      "gpuMemory": "0",
      "totalMemory": "134472257536",
      "diskSize": "7676310884352",
      "nodeIp": "172.16.100.14"
    },
    {
      "nodeName": "k8s2",
      "cpuCores": "80",
      "gpuNum": "1",
      "gpuMemory": "25769803776",
      "totalMemory": "134432251904",
      "diskSize": "7675682045952",
      "nodeIp": "172.16.100.15"
    }
  ]
});

</script>

<style scoped lang="scss">
.left {
  flex: 1;
}

.right {
  width: 170px;
  display: flex;
  justify-content: end;
}

.title {
  color: #3D4F62;
  font-size: 18px;
}

.tags {
  display: flex;
  gap: 60px;
  margin-top: 10px;

  span {
    font-size: 14px;
    color: #3D4F62;
  }
}

.wrap {
  display: flex;
  align-items: center;
}

.wrap-left,
.wrap-right {
  flex: 1;
  height: 320px;
  background: #F6F7F9;
  border-radius: 4px;
  overflow: auto;
  padding: 0 14px 14px 14px;

  .wrap-row {
    display: flex;
    gap: 18px;
    font-size: 12px;
    margin-left: 23px;

    div {
      color: #5A6B7D;

      span {
        color: #0B1524;
        margin-left: 4px;
      }
    }
  }
}

.wrap-center {
  width: 60px;
  flex-shrink: 0;
}
</style>

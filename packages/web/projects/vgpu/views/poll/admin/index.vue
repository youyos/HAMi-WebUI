<template>
  <list-header description="资源池管理用于统一管理和调度计算资源，支持资源的动态分配、回收与负载均衡，提升资源利用率和系统灵活性。">
    <template #actions>
      <el-button @click="editId = null; dialogVisible = true" style="margin-right: 24px;" type="primary"
        round>创建资源池</el-button>
    </template>
  </list-header>
  <div v-loading="loading" style="min-height: 90px;">
    <block-box
      v-for="{ poolId, poolName, nodeNum, cpuCores, gpuNum, availableMemory, totalMemory, diskSize, nodeList, linkUrl }, index in paginatedList"
      :key="poolId" style="margin: 15px 0 0 0;">
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
          <el-button @click="sendRouteChange(`/admin/vgpu/poll/admin/${poolId}?name=${poolName}`)"
            type="text">查看详情</el-button>
          <template v-if="index === 0 && currentPage === 1">
            <el-button @click="sendRouteChange(linkUrl, 'open')" type="text">配置</el-button>
          </template>
          <template v-else>
            <el-button
              @click="dialogVisible = true; editId = poolId; nodeSelect = nodeList.map(e => e.nodeIp); input = poolName"
              type="text">编辑</el-button>
            <el-button @click="() => handleDelete(poolId)" type="text">删除</el-button>
          </template>
        </div>
      </el-row>
    </block-box>
  </div>

  <!-- 分页组件 -->
  <el-pagination style="margin-top: 15px;" background v-model:current-page="currentPage" v-model:page-size="pageSize"
    layout="total, ->, sizes, jumper, prev, next" :page-sizes="[10, 20, 50, 100]" :total="list.length"
    @size-change="handleSizeChange" @current-change="handleCurrentChange" />

  <el-dialog @close="editId = null; input = ''; nodeSelect = []" v-model="dialogVisible"
    :title="editId ? '编辑资源池' : '创建资源池'" width="1180" :before-close="handleClose">
    <el-row :wrap="false" style="align-items: center;">
      <span style="flex-shrink: 0; margin-right: 14px;">资源池名称</span>
      <el-input style="flex: 1;" v-model="input" size="large" />
    </el-row>
    <div style="margin-top: 20px; margin-bottom: 10px;">
      <span>选择节点</span>
      <span style="float: right;">已选<span style="color: #3061D0; margin: 0 5px;">{{ nodeSelect.length
      }}</span>个节点</span>
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
      <div class="wrap-center"></div>
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
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button :loading="btnLoading" type="primary" @click="handleOk">确认</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="jsx">
import pollApi from '~/vgpu/api/poll';
import { ref, onMounted, computed, watchEffect } from 'vue';
import BlockBox from '@/components/BlockBox.vue';
import { Remove } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { bytesToGB } from '@/utils';
import useParentAction from '~/vgpu/hooks/useParentAction';

const { sendRouteChange } = useParentAction();

// 数据列表相关
const list = ref([])
const loading = ref(true)

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 对话框相关
const dialogVisible = ref(false)
const editId = ref(null)
const input = ref('')
const btnLoading = ref(false)

// 节点选择相关
const nodeList = ref([])
const nodeSelect = ref([])

// 计算分页后的数据
const paginatedList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return list.value.slice(start, end)
})

// watchEffect(()=>{
//   console.log(currentPage.value, pageSize.value, 88)
// })

// 分页大小变化
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1 // 切换每页条数时重置到第一页
}

// 当前页变化
const handleCurrentChange = (val) => {
  currentPage.value = val
}

// 确认操作
const handleOk = async () => {
  if (!input.value) {
    ElMessage({
      message: '请输入资源池名称',
      type: 'warning',
    })
    return;
  }
  if (!nodeSelect.value.length) {
    ElMessage({
      message: '请选择节点',
      type: 'warning',
    })
    return;
  }
  btnLoading.value = true;
  try {
    const nodes = nodeList.value.filter(e => nodeSelect.value.includes(e.nodeIp)).map(e => ({ nodeIp: e.nodeIp, nodeName: e.nodeName }))
    let res;
    if (editId.value) {
      res = await pollApi.update({
        pool_id: editId.value,
        pool_name: input.value,
        nodes
      })
    } else {
      res = await pollApi.create({
        pool_name: input.value,
        nodes
      })
    }
    if (res?.code === 200) {
      getList();
      dialogVisible.value = false;
    }
  } finally {
    btnLoading.value = false;
  }

  btnLoading.value = false;
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

// 删除操作
const handleDelete = async (id) => {
  ElMessageBox.confirm(`确定要删除当前资源池吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
    beforeClose: async (action, instance, done) => {
      if (action === 'confirm') {
        instance.confirmButtonLoading = true;
        const res = await pollApi.delete({ pool_id: id });
        if (res.code === 200) {
          ElMessage.success('删除成功');
          getList();
          done();
        }
        instance.confirmButtonLoading = false;
      } else {
        done();
      }
    }
  })
}

// 获取资源池列表
const getList = async () => {
  loading.value = true
  const res = await pollApi.getPollList()
  loading.value = false
  list.value = res.data
  total.value = res.data.length
}

// 初始化加载数据
onMounted(async () => {
  getList();
  const res = await pollApi.getNodeList()
  nodeList.value = res.data
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
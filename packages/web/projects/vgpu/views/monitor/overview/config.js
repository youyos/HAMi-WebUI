export const rangeConfigInit = [
  {
    title: '资源分配趋势',
    dataSource: [
      {
        name: 'vGPU',
        query: `sum(hami_container_vgpu_allocated) / sum(hami_vgpu_count) * 100`,
        data: [],
        type: 'line',
        areaStyle: {
          normal: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(250, 200, 88, 0.16)',
                },
                {
                  offset: 1,
                  color: 'rgba(250, 200, 88, 0.00)',
                },
              ],
              global: false,
            },
          },
        },
        itemStyle: {
          color: 'rgb(250, 200, 88)',
        },
        lineStyle: {
          color: 'rgb(250, 200, 88)',
        },
      },
      {
        name: '算力',
        query: `sum(hami_container_vcore_allocated) / sum(hami_core_size) * 100`,
        data: [],
        type: 'line',
        areaStyle: {
          normal: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(84, 112, 198, 0.16)',
                },
                {
                  offset: 1,
                  color: 'rgba(84, 112, 198, 0.00)',
                },
              ],
              global: false,
            },
          },
        },
        itemStyle: {
          color: 'rgb(84, 112, 198)',
        },
        lineStyle: {
          color: 'rgb(84, 112, 198)',
        },
      },
      {
        name: '显存',
        query: `sum(hami_container_vmemory_allocated) / sum(hami_memory_size) * 100`,
        data: [],
        type: 'line',
        areaStyle: {
          normal: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(34, 139, 34, 0.16)',
                },
                {
                  offset: 1,
                  color: 'rgba(34, 139, 34, 0.00)',
                },
              ],
              global: false,
            },
          },
        },
        itemStyle: {
          color: 'rgb(145, 204, 117)',
        },
        lineStyle: {
          color: 'rgb(145, 204, 117)',
        },
      },
      {
        name: 'CPU',
        query: `sum(hami_core_used) / sum(hami_core_size) * 100`,
        data: [],
        type: 'line',
        areaStyle: {
          normal: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(255, 99, 71, 0.16)',
                },
                {
                  offset: 1,
                  color: 'rgba(255, 99, 71, 0.00)',
                },
              ],
              global: false,
            },
          },
        },
        itemStyle: {
          color: 'rgb(255, 99, 71)',
        },
        lineStyle: {
          color: 'rgb(255, 99, 71)',
        },
      },
      {
        name: '内存',
        query: `sum(hami_memory_used) / sum(hami_memory_size) * 100`,
        data: [],
        type: 'line',
        areaStyle: {
          normal: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(138, 43, 226, 0.16)',
                },
                {
                  offset: 1,
                  color: 'rgba(138, 43, 226, 0.00)',
                },
              ],
              global: false,
            },
          },
        },
        itemStyle: {
          color: 'rgb(138, 43, 226)',
        },
        lineStyle: {
          color: 'rgb(138, 43, 226)',
        },
      }
    ]
  },
  {
    title: '资源使用趋势',
    dataSource: [
      {
        name: '算力',
        query: `avg(hami_core_util_avg)`,
        data: [],
        areaStyle: {
          normal: {
            color: {
              type: 'linear',
              x: 0, // 渐变起始点 0%
              y: 0, // 渐变起始点 0%
              x2: 0, // 渐变结束点 100%
              y2: 1, // 渐变结束点 100%
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(84, 112, 198, 0.16)', // 渐变起始颜色
                },
                {
                  offset: 1,
                  color: 'rgba(84, 112, 198, 0.00)', // 渐变结束颜色
                },
              ],
              global: false, // 缺省为 false
            },
          },
        },
        itemStyle: {
          color: 'rgb(84, 112, 198)', // 设置线条颜色为橙色
        },
        lineStyle: {
          color: 'rgb(84, 112, 198)', // 设置线条颜色为橙色
        },
      },
      {
        name: '显存',
        query: `sum(hami_memory_used) / sum(hami_memory_size) * 100`,
        data: [],
        areaStyle: {
          normal: {
            color: {
              type: 'linear',
              x: 0, // 渐变起始点 0%
              y: 0, // 渐变起始点 0%
              x2: 0, // 渐变结束点 100%
              y2: 1, // 渐变结束点 100%
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(34, 139, 34, 0.16)', // 渐变起始颜色
                },
                {
                  offset: 1,
                  color: 'rgba(34, 139, 34, 0.00)', // 渐变结束颜色
                },
              ],
              global: false, // 缺省为 false
            },
          },
        },
        itemStyle: {
          color: 'rgb(145, 204, 117)', // 设置线条颜色为橙色
        },
        lineStyle: {
          color: 'rgb(145, 204, 117)', // 设置线条颜色为橙色
        },
      },
      {
        name: 'CPU',
        query: `100 * (1 - sum(irate(node_cpu_seconds_total{mode="idle"}[1m])) / count(node_cpu_seconds_total{mode="idle"}))`,
        data: [],
        type: 'line',
        areaStyle: {
          normal: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(255, 99, 71, 0.16)',
                },
                {
                  offset: 1,
                  color: 'rgba(255, 99, 71, 0.00)',
                },
              ],
              global: false,
            },
          },
        },
        itemStyle: {
          color: 'rgb(255, 99, 71)',
        },
        lineStyle: {
          color: 'rgb(255, 99, 71)',
        },
      },
      {
        name: '内存',
        query: `(1 - sum(node_memory_MemAvailable_bytes) / sum(node_memory_MemTotal_bytes)) * 100`,
        data: [],
        type: 'line',
        areaStyle: {
          normal: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(138, 43, 226, 0.16)',
                },
                {
                  offset: 1,
                  color: 'rgba(138, 43, 226, 0.00)',
                },
              ],
              global: false,
            },
          },
        },
        itemStyle: {
          color: 'rgb(138, 43, 226)',
        },
        lineStyle: {
          color: 'rgb(138, 43, 226)',
        },
      }
    ],
  },
];

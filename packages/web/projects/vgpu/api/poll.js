import request from '@/utils/request';

const apiPrefix = '/api/vgpu';

class pollApi {
  getPollList(params) {
    return request({
      url: apiPrefix + '/v1/resource/pool/list',
      method: 'GET',
      params,
    });
  }

  getNodeList(params) {
    return request({
      url: apiPrefix + '/v1/available/nodes',
      method: 'GET',
      params,
    });
  }

  create(data) {
    return request({
      url: apiPrefix + '/v1/resource/pool/create',
      method: 'POST',
      data,
    });
  }

  update(data) {
    return request({
      url: apiPrefix + '/v1/resource/pool/update',
      method: 'POST',
      data,
    });
  }

  delete(data) {
    return request({
      url: apiPrefix + '/v1/resource/pool/delete',
      method: 'POST',
      data,
    });
  }

  remove(data) {
    return request({
      url: apiPrefix + '/v1/resource/pool/removeNode',
      method: 'POST',
      data,
    });
  }

  getDetailNodeList(data) {
    return request({
      url: apiPrefix + '/v1/resource/pool/detail',
      method: 'post',
      data,
    });
  }

  // getNodes(data) {
  //   return request({
  //     url: apiPrefix +  '/v1/nodes',
  //     method: 'POST',
  //     data,
  //   });
  // }

  // getNodeDetail(params) {
  //   return request({
  //     url: apiPrefix +  '/v1/node',
  //     method: 'GET',
  //     params,
  //   });
  // }

  // getNodeListReq(data) {
  //   return request(this.getNodeList(data));
  // }
}

export default new pollApi();

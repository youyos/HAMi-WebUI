import request from '@/utils/request';

const apiPrefix = '/api/vgpu';


class nodeApi {
  getNodeList(data) {
    return {
      url: apiPrefix + '/v1/nodes',
      method: 'POST',
      data,
    };
  }

  getNodes(data) {
    return request({
      url: apiPrefix + '/v1/nodes',
      method: 'POST',
      data,
    });
  }

  getNodeDetail(params) {
    return request({
      url: apiPrefix + '/v1/node',
      method: 'GET',
      params,
    });
  }

  getNodeListReq(data) {
    return request(this.getNodeList(data));
  }

  discoveredNodes(data) {
    return request({
      url: apiPrefix + '/v1/node/discovered',
      method: 'POST',
      data,
    });
  }
  joinNodes(data) {
    return request({
      url: apiPrefix + '/v1/node/join',
      method: 'POST',
      data,
    });
  }
}

export default new nodeApi();

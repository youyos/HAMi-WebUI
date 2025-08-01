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

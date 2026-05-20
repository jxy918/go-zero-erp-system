import request from './request'

// 活动日志API
export const activityApi = {
  // 获取活动日志列表
  getActivityList: (params) => {
    return request({
      url: '/activity/list',
      method: 'get',
      params
    })
  }
}

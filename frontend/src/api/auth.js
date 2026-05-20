import request from './request'

// 认证API
export const authApi = {
  // 登录
  login: (data) => {
    return request({
      url: '/auth/login',
      method: 'post',
      data
    })
  },

  // 登出
  logout: () => {
    return request({
      url: '/auth/logout',
      method: 'post'
    })
  },

  // 刷新Token
  refreshToken: (token) => {
    return request({
      url: '/auth/refresh',
      method: 'post',
      data: { token }
    })
  }
}
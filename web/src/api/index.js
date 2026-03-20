// Copyright © 2023 Ronnie Zhang (大脸怪). MIT License.

import { request } from '@/utils'

export default {
  // 获取用户信息
  getUser: () => request.post('/auth/detail'),
  // 刷新token
  refreshToken: (token) => request.post('/auth/refresh',{token}),
  // 登出
  logout: () => request.post('/auth/logout', {}, { needTip: false }),
  // 获取角色权限
  getRolePermissions: () => request.post('/roles/permissions'),



  
  // 切换当前角色
  switchCurrentRole: role => request.post(`/auth/current-role/switch/${role}`),
  // 验证菜单路径
  validateMenuPath: path => request.get(`/permission/menu/validate?path=${path}`),
}

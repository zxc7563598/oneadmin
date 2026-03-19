// Copyright © 2023 Ronnie Zhang (大脸怪). MIT License.

import { request } from '@/utils'

export default {
  changePassword: data => request.post('/auth/password', data),
  updateProfile: data => request.patch(`/user/profile/${data.id}`, data),
}

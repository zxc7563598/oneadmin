// Copyright © 2023 Ronnie Zhang (大脸怪). MIT License.

import { request } from '@/utils'

export default {
  // toggleRole: data => request.post('/auth/role/toggle', data),
  // getUser: () => request.get('/user/detail'),
  login: data => request.post('/auth/login', data, { needToken: false }),
}

import axios from 'axios'
import { ElMessage } from 'element-plus'

const pyService = axios.create({
  baseURL: import.meta.env.VITE_PY_API || '/pyapi',
  timeout: 120000,
})

pyService.interceptors.request.use(
  (config) => {
    config.headers = {
      'Content-Type': 'application/json',
      ...config.headers,
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

pyService.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (!error.response) {
      ElMessage({
        showClose: true,
        message: 'Python 分析服务暂不可用，请确认服务已启动（端口 8000）',
        type: 'warning',
        duration: 5000,
      })
      return Promise.reject(error)
    }

    const msg = error.response?.data?.detail || error.response?.statusText || '请求失败'
    ElMessage({
      showClose: true,
      message: msg,
      type: 'error',
    })
    return Promise.reject(error)
  }
)

export default pyService

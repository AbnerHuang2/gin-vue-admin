import pyService from '@/utils/pyRequest'

export const submitAnalysis = (data) => {
  return pyService({
    url: '/api/v1/selection/analyze',
    method: 'post',
    data,
  })
}

export const getTaskStatus = (taskId) => {
  return pyService({
    url: `/api/v1/selection/task/${taskId}/status`,
    method: 'get',
  })
}

export const getTaskReport = (taskId) => {
  return pyService({
    url: `/api/v1/selection/report/${taskId}`,
    method: 'get',
  })
}

export const getTaskList = (params) => {
  return pyService({
    url: '/api/v1/selection/tasks',
    method: 'get',
    params,
  })
}

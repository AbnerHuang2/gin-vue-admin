import pyService from '@/utils/pyRequest'

export const getCandidateList = (params) => {
  return pyService({
    url: '/api/v1/discovery/candidates',
    method: 'get',
    params,
  })
}

export const triggerDiscovery = () => {
  return pyService({
    url: '/api/v1/discovery/run',
    method: 'post',
  })
}

export const submitCandidate = (candidateId) => {
  return pyService({
    url: `/api/v1/discovery/candidates/${candidateId}/submit`,
    method: 'post',
  })
}

export const getDiscoveryStats = (params) => {
  return pyService({
    url: '/api/v1/discovery/stats',
    method: 'get',
    params,
  })
}

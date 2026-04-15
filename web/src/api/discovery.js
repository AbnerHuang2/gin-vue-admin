import pyService from '@/utils/pyRequest'

export const getCandidateList = (params) => {
  return pyService({
    url: '/api/v1/discovery/candidates',
    method: 'get',
    params,
  })
}

export const triggerDiscovery = (categories = []) => {
  return pyService({
    url: '/api/v1/discovery/run',
    method: 'post',
    data: { categories },
  })
}

export const submitCandidate = (candidateId) => {
  return pyService({
    url: `/api/v1/discovery/candidates/${candidateId}/submit`,
    method: 'post',
  })
}

export const analyzeCandidate = (candidateId) => {
  return pyService({
    url: `/api/v1/discovery/candidates/${candidateId}/analyze`,
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

// ── 品类关键词 CRUD ──

export const getKeywordList = (params) => {
  return pyService({
    url: '/api/v1/discovery/keywords',
    method: 'get',
    params,
  })
}

export const createKeyword = (data) => {
  return pyService({
    url: '/api/v1/discovery/keywords',
    method: 'post',
    data,
  })
}

export const updateKeyword = (id, data) => {
  return pyService({
    url: `/api/v1/discovery/keywords/${id}`,
    method: 'put',
    data,
  })
}

export const deleteKeyword = (id) => {
  return pyService({
    url: `/api/v1/discovery/keywords/${id}`,
    method: 'delete',
  })
}

import { Dispatch, Action } from 'redux'
import axios from 'axios'
import { ServerProject, Project, converter } from '@/entities/project'
import { Repository } from '@/entities/repository'

export const RequestGetProjects = 'RequestGetProjects' as const
export const ReceiveGetProjects = 'ReceiveGetProjects' as const
export const OpenNew = 'OpenNew' as const
export const CloseNew = 'CloseNew' as const
export const RequestGetRepositories = 'RequestGetRepositories' as const
export const ReceiveGetRepositories = 'ReceiveGetRepositories' as const

export const requestGetProjects = () => ({
  type: RequestGetProjects
})

export const receiveGetProjects = (projects: Array<Project>) => {
  return {
    type: ReceiveGetProjects,
    payload: projects
  }
}

export const getProjects = () => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestGetProjects())
    axios.get<Array<ServerProject>>('/api/projects').then(res => {
      const data: Array<Project> = res.data.map(p => converter(p))
      dispatch(receiveGetProjects(data))
    })
  }
}

export const openNew = () => ({
  type: OpenNew
})

export const closeNew = () => ({
  type: CloseNew
})

export const requestGetRepositories = () => ({
  type: RequestGetRepositories
})

export const receiveGetRepositories = (repositories: Array<Repository>) => {
  return {
    type: ReceiveGetRepositories,
    payload: repositories
  }
}

export const getRepositories = () => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestGetRepositories())
    axios.get<Array<Repository>>('/api/github/repositories').then(res => {
      dispatch(receiveGetRepositories(res.data))
    })
  }
}

type Actions = ReturnType<
  | typeof requestGetProjects
  | typeof receiveGetProjects
  | typeof openNew
  | typeof closeNew
  | typeof requestGetRepositories
  | typeof receiveGetRepositories
>

export default Actions

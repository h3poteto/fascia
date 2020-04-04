import { Dispatch, Action } from 'redux'
import axios from 'axios'

type ServerProject = {
  ID: number
  UserID: number
  Title: string
  Description: string
  RepositoryID: number | null
  ShowIssues: boolean
  ShowPullRequests: boolean
}

export type Project = {
  id: number
  userID: number
  title: string
  description: string
  repositoryID: number | null
  showIssues: boolean
  showPullRequests: boolean
}

export type Repository = {
  id: number
  full_name: string
}

export const RequestGetProjects = 'RequestGetProjects' as const
export const ReceiveGetProjects = 'ReceiveGetProjects' as const
export const OpenNew = 'OpenNew' as const
export const CloseNew = 'CloseNew' as const
export const RequestGetRepositories = 'RequestGetRepositories' as const
export const ReceiveGetRepositories = 'ReceiveGetRepositories' as const
export const RequestCreateProject = 'RequestCreateProject' as const
export const ReceiveCreateProject = 'ReceiveCreateProject' as const

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
      const data: Array<Project> = res.data.map(p => {
        return {
          id: p.ID,
          userID: p.UserID,
          title: p.Title,
          description: p.Description,
          repositoryID: p.RepositoryID,
          showIssues: p.ShowIssues,
          showPullRequests: p.ShowPullRequests
        }
      })
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

export const requestCreateProject = () => ({
  type: RequestCreateProject
})

export const receiveCreateProject = (project: Project) => ({
  type: ReceiveCreateProject,
  payload: project
})

export const createProject = (params: any) => {
  return async (dispatch: Function) => {
    dispatch(requestCreateProject())
    return axios.post<ServerProject>('/api/projects', params).then(res => {
      const data: Project = {
        id: res.data.ID,
        userID: res.data.UserID,
        title: res.data.Title,
        description: res.data.Description,
        repositoryID: res.data.RepositoryID,
        showIssues: res.data.ShowIssues,
        showPullRequests: res.data.ShowPullRequests
      }
      dispatch(receiveCreateProject(data))
      dispatch(getProjects())
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
  | typeof requestCreateProject
  | typeof receiveCreateProject
>

export default Actions

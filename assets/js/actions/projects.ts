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

export const RequestGetProjects = 'RequestGetProjects' as const
export const ReceiveGetProjects = 'ReceiveGetProjects' as const
export const OpenNew = 'OpenNew' as const
export const CloseNew = 'CloseNew' as const

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

type Actions = ReturnType<typeof requestGetProjects | typeof receiveGetProjects | typeof openNew | typeof closeNew>

export default Actions

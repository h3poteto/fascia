import axios from 'axios'

import { Project, ServerProject, getProjects } from '@/actions/projects'

export const RequestCreateProject = 'RequestCreateProject' as const
export const ReceiveCreateProject = 'ReceiveCreateProject' as const

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

type Actions = ReturnType<typeof requestCreateProject | typeof receiveCreateProject>

export default Actions

import { Dispatch, Action } from 'redux'
import axios from 'axios'

import { Project, ServerProject } from '@/entities/project'

export const ReceiveGetProject = 'ReceiveGetProject' as const
export const RequestGetProject = 'RequestGetProject' as const

export const requestGetProject = () => ({
  type: RequestGetProject
})

export const receiveGetProject = (project: Project) => ({
  type: ReceiveGetProject,
  payload: project
})

export const getProject = (id: number) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestGetProject())
    axios.get<ServerProject>(`/api/projects/${id}/show`).then(res => {
      const data: Project = {
        id: res.data.ID,
        userID: res.data.UserID,
        title: res.data.Title,
        description: res.data.Description,
        repositoryID: res.data.RepositoryID,
        showIssues: res.data.ShowIssues,
        showPullRequests: res.data.ShowPullRequests
      }
      dispatch(receiveGetProject(data))
    })
  }
}

type Actions = ReturnType<typeof requestGetProject | typeof receiveGetProject>

export default Actions

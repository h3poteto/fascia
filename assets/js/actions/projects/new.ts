import axios from 'axios'

import { getProjects } from '@/actions/projects'
import { ServerProject, Project, converter } from '@/entities/project'

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
      const data = converter(res.data)
      dispatch(receiveCreateProject(data))
      dispatch(getProjects())
    })
  }
}

type Actions = ReturnType<typeof requestCreateProject | typeof receiveCreateProject>

export default Actions

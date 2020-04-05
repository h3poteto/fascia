import axios from 'axios'
import { push } from 'connected-react-router'
import { Action, Dispatch } from 'redux'

export const RequestDeleteProject = 'RequestDeleteProject' as const
export const ReceiveDeleteProject = 'ReceiveDeleteProject' as const

export const requestDeleteProject = () => ({
  type: RequestDeleteProject
})

export const receiveDeleteProject = () => ({
  type: ReceiveDeleteProject
})

export const deleteProject = (id: number) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestDeleteProject())
    axios.delete<{}>(`/api/projects/${id}`).then(() => {
      dispatch(receiveDeleteProject())
      dispatch(push('/'))
    })
  }
}

type Actions = ReturnType<typeof requestDeleteProject | typeof receiveDeleteProject>

export default Actions

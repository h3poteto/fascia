import React from 'react'
import { RouteComponentProps } from 'react-router-dom'
import { ThunkDispatch } from 'redux-thunk'

import { RootStore } from '@/reducers/index'
import ListForm from './form.tsx'
import Actions, { getList, getListOptions } from '@/actions/projects/lists/show'
import ProjectActions, { getProject } from '@/actions/projects/show'
import UpdateActions, { updateList } from '@/actions/projects/lists/edit'


type Props = {
  dispatch: ThunkDispatch<any, any, Actions | ProjectActions | UpdateActions>
} & RootStore & RouteComponentProps<{ project_id: string, list_id: string}>

class Edit extends React.Component<Props> {
  componentDidMount() {
    const projectID = parseInt(this.props.match.params.project_id)
    const listID = parseInt(this.props.match.params.list_id)
    this.props.dispatch(getList(projectID, listID))
    this.props.dispatch(getProject(projectID))
    this.props.dispatch(getListOptions())
  }

  render() {
    const hide = () => {
      const projectID = this.props.match.params.project_id
      this.props.history.push(`/projects/${projectID}`)
    }

    const edit = (params: any) => {
      const projectID = parseInt(this.props.match.params.project_id)
      const listID = parseInt(this.props.match.params.list_id)
      this.props.dispatch(updateList(projectID, listID, params))
    }

    return (
      <div>
        <ListForm hide={hide} onSubmit={edit} list={this.props.list.list} project={this.props.list.project} list_options={this.props.list.list_options} />
      </div>
    )
  }
}

export default Edit

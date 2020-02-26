import React from 'react'
import { ThunkDispatch } from 'redux-thunk'

import Project from './projects/project.tsx'
import styles from './projects.scss'
import Actions, { getProjects } from '../actions/projects'
import { RootStore } from '../reducers/index'

type Props = {
  dispatch: ThunkDispatch<any, any, Actions>
} & RootStore

class ProjectsComponent extends React.Component<Props> {
  componentDidMount() {
    this.props.dispatch(getProjects())
    console.log(this.props)
  }

  render() {
    const projects = this.props.projects.projects
    return (
      <div className={styles.projects}>
        {projects.map(p => {
          return <Project title={p.title} />
        })}
        <div>{this.props.children}</div>
      </div>
    )
  }
}

export default ProjectsComponent

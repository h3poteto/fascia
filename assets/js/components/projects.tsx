import React from 'react'
import { ThunkDispatch } from 'redux-thunk'
import { Button } from 'react-bootstrap'

import Project from './projects/project.tsx'
import styles from './projects.scss'
import Actions, { getProjects, openNew, closeNew } from '../actions/projects'
import { RootStore } from '../reducers/index'
import New from './projects/new.tsx'

type Props = {
  dispatch: ThunkDispatch<any, any, Actions>
} & RootStore

class ProjectsComponent extends React.Component<Props> {
  componentDidMount() {
    this.props.dispatch(getProjects())
  }

  render() {
    const projects = this.props.projects.projects
    const   openNewModal = () => {
      return this.props.dispatch(openNew())
    }
    const closeNewModal = () =>  {
      return this.props.dispatch(closeNew())
    }

    return (
      <div className={styles.projects}>
        {projects.map(p => (
          <Project key={p.id} id={p.id} title={p.title} />
        ))}
        <Button className={styles.newButton} onClick={openNewModal}>New</Button>
        <New open={this.props.projects.newModal} close={closeNewModal}></New>
      </div>
    )
  }
}

export default ProjectsComponent

import React from 'react'

import Project from './projects/project.tsx'
import styles from './projects.scss'

type Props = {}

const ProjectsComponent: React.FC<Props> = ({ children }) => {
  const projects: Array<string> = ['hoge', 'fuga']
  return (
    <div className={styles.projects}>
      {projects.map(p => {
        return (
          <Project title={p} />
        )
      })}
      <div>{children}</div>
    </div>
  )
}

export default ProjectsComponent

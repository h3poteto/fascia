import React from 'react'
import { Link } from 'react-router-dom'

import styles from './project.scss'

type Props = {
  id: number
  title: string
}

const project: React.FC<Props> = props => {
  return (
    <Link className={styles.project} to={`/projects/${props.id}`}>
      {props.title}
    </Link>
  )
}

export default project

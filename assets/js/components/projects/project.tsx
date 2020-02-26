import React from 'react'
import styles from "./project.scss"

type Props = {
  title: string
}

const project: React.FC<Props> = (props) => {
  return (
    <div className={styles.project}>
      {props.title}
    </div>
  )
}

export default project

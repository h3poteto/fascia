import React from 'react'

type Props = {}

const projects: React.FC<Props> = ({ children }) => {
  return (
    <div>
      Projects
      <div>{children}</div>
    </div>
  )
}

export default projects

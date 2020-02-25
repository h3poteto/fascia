import React from 'react'

type Props = {}

const menu: React.FC<Props> = ({ children }) => {
  return (
    <div>
      <header className="top-header">
        header
      </header>
      <div className="contents">
        {children}
      </div>
    </div>
  )
}

export default menu

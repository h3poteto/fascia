import React from 'react'

class ListLoading extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    if (this.props.isLoading != undefined && this.props.isLoading) {
      return (
        <div className="list-loading">
          <div className="list-circle-wrapper">
            <div className="list-circle-body">
              <div className="list-spinner"></div>
            </div>
          </div>
        </div>
      )
    } else {
      return <span></span>
    }
  }
}

ListLoading.propTypes = {
  isLoading: React.PropTypes.bool,
}

export default ListLoading

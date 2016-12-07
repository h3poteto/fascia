import React from 'react'

class WholeLoading extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    if (this.props.isLoading) {
      return (
        <div className="whole-loading">
          <div className="whole-circle-wrapper">
            <div className="whole-circle-body">
              <div className="whole-spinner"></div>
            </div>
          </div>
        </div>
      )
    } else {
      return <div></div>
    }
  }
}

WholeLoading.propTypes = {
  isLoading: React.PropTypes.bool.isRequired,
}

export default WholeLoading

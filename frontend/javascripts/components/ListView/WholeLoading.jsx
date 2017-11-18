import React from 'react'
import PropTypes from 'prop-types'

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
  isLoading: PropTypes.bool.isRequired,
}

export default WholeLoading

import React from 'react'
import PropTypes from 'prop-types'
import Modal from 'react-modal'
import { Field, reduxForm } from 'redux-form'
import { GithubPicker } from 'react-color'

const customStyles = {
  overlay: {
    position          : 'fixed',
    top               : 0,
    left              : 0,
    right             : 0,
    bottom            : 0,
    backgroundColor   : 'rgba(255, 255, 255, 0.5)'
  },
  content: {
    position : 'fixed',
    top : '50%',
    left : '50%',
    right : 'auto',
    bottom : 'auto',
    marginRight : '-50%',
    transform : 'translate(-50%, -50%)'
  }
}

class NewListModal extends React.Component {
  componentWillReceiveProps(nextProps) {
    if (!nextProps.isListModalOpen) {
      this.handleInitialize()
    }
  }

  handleInitialize() {
    this.props.initialize({})
  }

  render() {
    const {
      handleSubmit,
      pristine,
      reset,
      submitting,
      onRequestClose,
      action,
      changeColor,
      color,
      projectID,
    } = this.props
    return (
      <Modal
        isOpen={this.props.isListModalOpen}
        onRequestClose={onRequestClose}
        style={customStyles}
        contentLabel="NewListModal"
      >
        <div className="list-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit((values) => { action(projectID, values) })}>
            <fieldset>
              <legend>Create List</legend>
              <label htmlFor="title">Title</label>
              <Field name="title" id="title" component="input" type="text" placeholder="List name" className="form-control" />
              <label htmlFor="color">Color</label>
              <div className="color-control-group">
                <div className="real-color" style={{backgroundColor: `#${color}`}}>　</div>
                <Field name="color" id="color" component="input" type="text" placeholder="008ed4" onChange={(e) => changeColor(e.target.value)} />
              </div>
              <GithubPicker
                onChangeComplete={(color) => {
                    this.props.array.removeAll('color')
                    this.props.array.push('color', color.hex.replace(/#/g, ''))
                    changeColor(color.hex.replace(/#/g, ''))
                }
                }
              />
              <div className="form-action">
                <button type="reset" className="pure-button pure-button-default" disabled={pristine || submitting} onClick={reset}>Reset</button>
                <button type="submit" className="pure-button pure-button-primary" disabled={pristine || submitting}>Create List</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
                    )
  }
}

NewListModal.propTypes = {
  initialize: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  pristine: PropTypes.bool,
  reset: PropTypes.func.isRequired,
  submitting: PropTypes.bool.isRequired,
  onRequestClose: PropTypes.func.isRequired,
  action: PropTypes.func.isRequired,
  projectID: PropTypes.string.isRequired,
  isListModalOpen: PropTypes.bool.isRequired,
  dirty: PropTypes.object,
  array: PropTypes.object,
  color: PropTypes.string,
  changeColor: PropTypes.func,
}

export default reduxForm({
  form: 'new-list-form',
})(NewListModal)

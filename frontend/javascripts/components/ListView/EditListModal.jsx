import React from 'react'
import PropTypes from 'prop-types'
import Modal from 'react-modal'
import { Field, reduxForm } from 'redux-form'
import { GithubPicker } from 'react-color'

import { RenderField, RenderColorField, validate } from './listForm'

const customStyles = {
  overlay : {
    position          : 'fixed',
    top               : 0,
    left              : 0,
    right             : 0,
    bottom            : 0,
    backgroundColor   : 'rgba(255, 255, 255, 0.5)'
  },
  content : {
    position : 'fixed',
    top : '50%',
    left : '50%',
    right : 'auto',
    bottom : 'auto',
    marginRight : '-50%',
    transform : 'translate(-50%, -50%)'
  }
}

class EditListModal extends React.Component {
  componentWillReceiveProps(nextProps) {
    // modalをcloseするタイミングでは初期化しておかないと，別のlistを選択したときに，現在の編集分が残っている可能性がある
    if (!nextProps.dirty || !nextProps.isListEditModalOpen) {
      this.handleInitialize(nextProps)
    }
  }

  handleInitialize(props) {
    const initData = {
      'title': props.list.Title,
      'color': props.list.Color,
      'option_id': props.list.ListOptionID,
    }

    this.props.initialize(initData)
  }

  listAction(project, listOptions) {
    if (project == null || project.RepositoryID == undefined || project.RepositoryID == null || project.RepositoryID == 0) {
      return null
    } else {
      return (
        <div>
          <label htmlFor="option_id">action</label>
          <Field name="option_id" id="option_id" component="select" className="form-control">
            <option value="0">nothing</option>
            {listOptions.map(function(option, index) {
               return <option key={index} value={option.ID}>{option.Action}</option>
             }, this)}
          </Field>
        </div>
      )
    }
  }

  render() {
    const {
      handleSubmit,
      pristine,
      reset,
      submitting,
      onRequestClose,
      action,
      project,
      list,
      listOptions,
      changeColor,
      color,
    } = this.props

    return (
      <Modal
        isOpen={this.props.isListEditModalOpen}
        onRequestClose={onRequestClose}
        style={customStyles}
        contentLabel="EditListModal"
      >
        <div className="list-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit((values) => { action(project.ID, list.ID, values) })}>
            <fieldset>
              <legend>Edit List</legend>
              <label htmlFor="title">Title</label>
              <Field component={RenderField} name="title" id="title" type="text" className="form-control" />
              <label htmlFor="color">Color</label>
              <Field component={RenderColorField} name="color" type="text" placeholder="008ed4" color={color} onChange={(e) => changeColor(e.target.value)} />
              <GithubPicker
                onChangeComplete={(color) => {
                    this.props.array.removeAll('color')
                    this.props.array.push('color', color.hex.replace(/#/g, ''))
                    changeColor(color.hex.replace(/#/g, ''))
                }
                }
              />
              {this.listAction(project, listOptions)}
              <div className="form-action">
                <button type="reset" className="pure-button pure-button-default" disabled={pristine || submitting} onClick={reset}>Reset</button>
                <button type="submit" className="pure-button pure-button-primary" disabled={pristine || submitting}>Update List</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}

EditListModal.propTypes = {
  initialize: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  pristine: PropTypes.bool,
  reset: PropTypes.func.isRequired,
  submitting: PropTypes.bool.isRequired,
  onRequestClose: PropTypes.func.isRequired,
  action: PropTypes.func.isRequired,
  project: PropTypes.object,
  list: PropTypes.object,
  listOptions: PropTypes.array,
  isListEditModalOpen: PropTypes.bool.isRequired,
  dirty: PropTypes.object,
  array: PropTypes.object,
  color: PropTypes.string,
  changeColor: PropTypes.func,
}

export default reduxForm({
  form: 'edit-list-form',
  validate,
})(EditListModal)

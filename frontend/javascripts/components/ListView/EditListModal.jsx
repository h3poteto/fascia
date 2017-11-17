import React from 'react'
import Modal from 'react-modal'
import { Field, reduxForm } from 'redux-form'
import { GithubPicker } from 'react-color'

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
              <Field name="title" id="title" component="input" type="text" className="form-control" />
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
  initialize: React.PropTypes.func.isRequired,
  handleSubmit: React.PropTypes.func.isRequired,
  pristine: React.PropTypes.bool,
  reset: React.PropTypes.func.isRequired,
  submitting: React.PropTypes.bool.isRequired,
  onRequestClose: React.PropTypes.func.isRequired,
  action: React.PropTypes.func.isRequired,
  project: React.PropTypes.object,
  list: React.PropTypes.object,
  listOptions: React.PropTypes.array,
  isListEditModalOpen: React.PropTypes.bool.isRequired,
  dirty: React.PropTypes.object,
  array: React.PropTypes.object,
  color: React.PropTypes.string,
  changeColor: React.PropTypes.func,
}

export default reduxForm({
  form: 'edit-list-form',
})(EditListModal)

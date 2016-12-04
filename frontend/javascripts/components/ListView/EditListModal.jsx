import React from 'react'
import Modal from 'react-modal'
import { Field, reduxForm } from 'redux-form'

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
  constructor(props) {
    super(props)
  }

  listAction(project, listOptions) {
    if (project == null || project.RepositoryID == undefined || project.RepositoryID == null || project.RepositoryID == 0) {
      return null
    } else {
      return (
        <div>
          <label htmlFor="option_id">action</label>
          <Field name="optiond_id" id="option_id" component="select" className="form-control">
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
    } = this.props
    return (
      <Modal
          isOpen={this.props.isListEditModalOpen}
          onRequestClose={onRequestClose}
          style={customStyles}
      >
        <div className="list-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit((values) => { action(project.ID, list.ID, values) })}>
            <fieldset>
              <legend>Edit List</legend>
              <label htmlFor="title">Title</label>
              <Field name="title" id="title" component="input" type="text" className="form-control" />
              <label htmlFor="color">Color</label>
              <Field name="color" id="color" component="input" type="text" className="form-control" />
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

export default reduxForm({
  form: 'edit-list-form',
})(EditListModal)

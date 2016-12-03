import React from 'react'
import Modal from 'react-modal'
import { Field, reduxForm } from 'redux-form'

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

class NewProjectModal extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    const {
      handleSubmit,
      pristine,
      reset,
      submitting,
      onRequestClose,
      action,
      repositories,
    } = this.props
    return (
      <Modal
          isOpen={this.props.isModalOpen}
          onRequestClose={onRequestClose}
          style={customStyles}
      >
        <div className="project-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit((values) => { action(values) })}>
            <fieldset>
              <legend>Create Project</legend>
              <label htmlFor="title">Title</label>
              <Field name="title" id="title" component="input" type="text" placeholder="Project name" className="form-control" />
              <label htmlFor="description">Description</label>
              <Field name="description" id="description" component="textarea" placeholder="Description" className="form-control" />
              <label htmlFor="repository_id">GitHub</label>
              <Field name="repository_id" id="repository_id" component="select" className="form-control">
                <option value="0">--</option>
                {repositories.map(function(repo, index) {
                  return <option key={index} value={repo.id}>{repo.full_name}</option>
                 }, this)}
              </Field>
              <div className="form-action">
                <button type="reset" className="pure-button pure-button-default" disabled={pristine || submitting} onClick={reset}>Reset</button>
                <button type="submit" className="pure-button pure-button-primary" disabled={pristine || submitting}>Create Project</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}

export default reduxForm({
  form: 'new-project-form',
})(NewProjectModal)

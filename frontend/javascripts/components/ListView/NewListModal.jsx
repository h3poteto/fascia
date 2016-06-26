import React from 'react'
import Modal from 'react-modal'

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

export default class NewListModal extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    return (
      <Modal
          isOpen={this.props.isListModalOpen}
          onRequestClose={this.props.closeNewListModal}
          style={customStyles}
      >
        <div className="list-form">
          <form className="pure-form pure-form-stacked">
            <fieldset>
              <legend>Create List</legend>
              <label htmlFor="title">Title</label>
              <input id="title" name="title" type="text" value={this.props.newList.title} onChange={this.props.updateNewListTitle} placeholder="List Name" className="form-control" />
              <label htmlFor="color">Color</label>
              <input id="color" name="color" type="text" value={this.props.newList.color} onChange={this.props.updateNewListColor} className="form-control" />
              <div className="form-action">
                <button onClick={e => this.props.fetchCreateList(this.props.projectID, this.props.newList.title, this.props.newList.color)} className="pure-button pure-button-primary" type="button">Create List</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}
